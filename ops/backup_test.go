package ops

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {

	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	defer os.RemoveAll(dataDir)
	assert.NoError(t, err)
	backupDir, err := ioutil.TempDir("", "ind9-rocks-backup")
	defer os.RemoveAll(backupDir)
	assert.NoError(t, err)

	db := createDummyDB(t, dataDir)
	db.Close()
	err = DoBackup(dataDir, backupDir)
	assert.NoError(t, err)

	assert.True(t, Exists(filepath.Join(backupDir, LatestBackup)))
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil
}

func TestRecursiveBackup(t *testing.T) {
	baseDataDir, err := ioutil.TempDir("", "baseDataDir")
	err = os.MkdirAll(baseDataDir, os.ModePerm)
	defer os.RemoveAll(baseDataDir)
	assert.NoError(t, err)

	baseBackupDir, err := ioutil.TempDir("", "baseBackupDir")
	err = os.MkdirAll(baseBackupDir, os.ModePerm)
	defer os.RemoveAll(baseBackupDir)
	assert.NoError(t, err)

	paths := []string{
		"1/store_1/",
		"1/store_2/",
		"2/store_1/",
		"2/store_2/",
	}

	for _, relLocation := range paths {
		err = os.MkdirAll(filepath.Join(baseDataDir, relLocation), os.ModePerm)
		assert.NoError(t, err)
		db := createDummyDB(t, filepath.Join(baseDataDir, relLocation))
		db.Close()
	}

	err = walkSourceDir(baseDataDir, baseBackupDir)
	assert.NoError(t, err)

	for _, relLocation := range paths {
		assert.True(t, Exists(filepath.Join(baseBackupDir, relLocation, LatestBackup)))
	}
}
