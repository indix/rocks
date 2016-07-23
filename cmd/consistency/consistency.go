package consistency

import (
	"fmt"
	log "log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ashwanthkumar/golang-utils/worker"
	"github.com/hashicorp/go-multierror"
	"github.com/ind9/rocks/cmd"
	"github.com/ind9/rocks/cmd/statistics"
	"github.com/spf13/cobra"
	"github.com/tecbot/gorocksdb"
)

var source string
var restoredLoc string
var threads int
var recursive bool
var paranoid bool

// Work struct contains Source and Restore locations for checking consistency
type Work struct {
	Source   string
	Restore  string
	Paranoid bool
}

// Result of DoConsistency
type Result struct {
	SourceDir     string
	RestoredDir   string
	Err           error
	SourceCount   int64
	RestoredCount int64
}

// ToString representation of the Result
func (result *Result) String() string {
	consistent := "FAIL"
	if result.IsConsistent() {
		consistent = "PASS"
	}
	return fmt.Sprintf(`Consistency Result = %s
SourceDir: %s RestoreDir: %s
Store Count: %d Restore Count: %d`, consistent, result.SourceDir, result.RestoredDir, result.SourceCount, result.RestoredCount)
}

// IsConsistent checks if the result is consistent
func (result *Result) IsConsistent() bool {
	return result.SourceCount == result.RestoredCount
}

var consistency = &cobra.Command{
	Use:   "consistency",
	Short: "Checks for the consistency between rocks store and it's corresponding restore using row counts",
	Long:  "Checks for the consistency between rocks store and it's corresponding restore using row counts",
	Run:   cmd.AttachHandler(checkConsistency),
}

func checkConsistency(args []string) (err error) {
	if source == "" {
		return fmt.Errorf("--src-dir was not set")
	}

	if restoredLoc == "" {
		return fmt.Errorf("--restore-dir was not set")
	}

	if recursive {
		return DoRecursiveConsistency(source, restoredLoc, threads)
	}
	result := DoConsistency(source, restoredLoc, paranoid)
	log.Println(result.String())
	return result.Err
}

// DoRecursiveConsistency checks for consistency recursively
func DoRecursiveConsistency(source, restore string, threads int) error {
	log.Printf("Initializing consistency check between %s as data directory and %s as it's restore directory\n", source, restore)

	workerPool := worker.Pool{
		MaxWorkers: threads,
		Op: func(request worker.Request) error {
			work := request.(Work)
			result := DoConsistency(work.Source, work.Restore, work.Paranoid)
			log.Println(result.String())
			return result.Err
		},
	}
	workerPool.Initialize()

	err := filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if info.Name() == cmd.Current {
			sourceDbLoc := filepath.Dir(path)
			sourceDbRelative, err := filepath.Rel(source, sourceDbLoc)
			if err != nil {
				return err
			}
			restoreDbLoc := filepath.Join(restore, sourceDbRelative)

			work := Work{
				Source:   sourceDbLoc,
				Restore:  restoreDbLoc,
				Paranoid: paranoid,
			}
			workerPool.AddWork(work)
			return filepath.SkipDir
		}
		return walkErr
	})

	var result error
	if errFromWorkers := workerPool.Join(); errFromWorkers != nil {
		result = multierror.Append(result, errFromWorkers)
	}

	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// DoConsistency checks for consistency between rocks source store and its restore
func DoConsistency(source, restore string, paranoid bool) Result {
	sourceOpts := gorocksdb.NewDefaultOptions()
	sourceOpts.SetParanoidChecks(paranoid)
	sourceDb, err := gorocksdb.OpenDb(sourceOpts, source)
	if err != nil {
		return Result{source, restore, err, 0, 0}
	}
	defer sourceDb.Close()

	restoreOpts := gorocksdb.NewDefaultOptions()
	restoreOpts.SetParanoidChecks(paranoid)
	restoredDb, err := gorocksdb.OpenDb(restoreOpts, restore)
	if err != nil {
		return Result{source, restore, err, 0, 0}
	}
	defer restoredDb.Close()

	sourceRowCount, err := statistics.DoStatsWithDB(sourceDb)
	if err != nil {
		return Result{source, restore, err, sourceRowCount, 0}
	}
	restoredRowCount, err := statistics.DoStatsWithDB(restoredDb)
	if err != nil {
		return Result{source, restore, err, sourceRowCount, restoredRowCount}
	}

	return Result{source, restore, err, sourceRowCount, restoredRowCount}
}

func init() {
	cmd.Rocks.AddCommand(consistency)

	consistency.PersistentFlags().StringVar(&source, "src-dir", "", "Original data location")
	consistency.PersistentFlags().StringVar(&restoredLoc, "restore-dir", "", "Restored location")
	consistency.PersistentFlags().BoolVar(&recursive, "recursive", false, "Recursively check for row counts across dbs")
	consistency.PersistentFlags().BoolVar(&paranoid, "paranoid", false, "Do paranoid checks on the DB")
	consistency.PersistentFlags().IntVar(&threads, "threads", 2*runtime.NumCPU(), "Number of threads to do backup")
}
