package restore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ind9/rocks/backup"
	"github.com/ind9/rocks/ops"
	"github.com/stretchr/testify/assert"
)

func TestRestore(t *testing.T) {
	backupDir, err := ioutil.TempDir("", "ind9-rocks-backup")
	assert.NoError(t, err)
	defer os.RemoveAll(backupDir)
	restoreDir, err := ioutil.TempDir("", "ind9-rocks-restore")
	assert.NoError(t, err)
	defer os.RemoveAll(restoreDir)
	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	assert.NoError(t, err)
	defer os.RemoveAll(dataDir)

	ops.WriteTestDB(t, dataDir)

	err = backup.DoBackup(dataDir, backupDir)
	assert.NoError(t, err)
	assert.True(t, ops.Exists(filepath.Join(backupDir, LatestBackup)))

	err = DoRestore(backupDir, restoreDir, restoreDir, false)
	assert.NoError(t, err)
	assert.True(t, ops.Exists(filepath.Join(restoreDir, ops.Current)))
}

func TestRecursiveRestore(t *testing.T) {
	const ShardCount = 3
	const DBsInEachShard = 3
	var paths []string
	for shard := 0; shard < ShardCount; shard++ {
		for db := 0; db < DBsInEachShard; db++ {
			path := fmt.Sprintf("%d/store_%d/", shard, db)
			paths = append(paths, path)
		}
	}

	backupThreads := 1
	restoreThreads := 1
	RecursivelyTestRestore(t, paths, backupThreads, restoreThreads)
}

func TestRecursiveRestoreParallely(t *testing.T) {
	const ShardCount = 10
	const DBsInEachShard = 3
	var paths []string
	for shard := 0; shard < ShardCount; shard++ {
		for db := 0; db < DBsInEachShard; db++ {
			path := fmt.Sprintf("%d/store_%d/", shard, db)
			paths = append(paths, path)
		}
	}

	backupThreads := 1
	restoreThreads := 5
	RecursivelyTestRestore(t, paths, backupThreads, restoreThreads)
}

func RecursivelyTestRestore(t *testing.T, paths []string, backupThreads, restoreThreads int) {

	baseDataDir, err := ioutil.TempDir("", "baseDataDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseDataDir)

	baseBackupDir, err := ioutil.TempDir("", "baseBackupDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseBackupDir)

	baseRestoreDir, err := ioutil.TempDir("", "baseRestoreDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseRestoreDir)

	// recursively write data
	for _, relLocation := range paths {
		err = os.MkdirAll(filepath.Join(baseDataDir, relLocation), os.ModePerm)
		assert.NoError(t, err)
		ops.WriteTestDB(t, filepath.Join(baseDataDir, relLocation))
	}

	// recursive backup + assert it
	err = backup.DoRecursiveBackup(baseDataDir, baseBackupDir, 1)
	assert.NoError(t, err)
	for _, relLocation := range paths {
		assert.True(t, ops.Exists(filepath.Join(baseBackupDir, relLocation, LatestBackup)))
	}

	err = DoRecursiveRestore(baseBackupDir, baseRestoreDir, baseRestoreDir, 5, true)
	assert.NoError(t, err)
	for _, relLocation := range paths {
		assert.True(t, ops.Exists(filepath.Join(baseRestoreDir, relLocation, ops.Current)))
	}
}
