package ops

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tecbot/gorocksdb"
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

	WriteTestDB(t, dataDir)

	err = DoBackup(dataDir, backupDir)
	assert.NoError(t, err)
	assert.True(t, Exists(filepath.Join(backupDir, LatestBackup)))

	err = DoRestore(backupDir, restoreDir, restoreDir, false)
	assert.NoError(t, err)
	assert.True(t, Exists(filepath.Join(restoreDir, Current)))
}

func TestRecursiveRestore(t *testing.T) {
	baseDataDir, err := ioutil.TempDir("", "baseDataDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseDataDir)

	baseBackupDir, err := ioutil.TempDir("", "baseBackupDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseBackupDir)

	baseRestoreDir, err := ioutil.TempDir("", "baseRestoreDir")
	assert.NoError(t, err)
	defer os.RemoveAll(baseRestoreDir)

	paths := []string{
		"1/store_1/",
		"1/store_2/",
		"2/store_1/",
		"2/store_2/",
	}

	// recursively write data
	for _, relLocation := range paths {
		err = os.MkdirAll(filepath.Join(baseDataDir, relLocation), os.ModePerm)
		assert.NoError(t, err)
		WriteTestDB(t, filepath.Join(baseDataDir, relLocation))
	}

	// recursive backup + assert it
	err = DoRecursiveBackup(baseDataDir, baseBackupDir, 1)
	assert.NoError(t, err)
	for _, relLocation := range paths {
		assert.True(t, Exists(filepath.Join(baseBackupDir, relLocation, LatestBackup)))
	}

	err = DoRecursiveRestore(baseBackupDir, baseRestoreDir, baseRestoreDir, 1, true)
	assert.NoError(t, err)
	for _, relLocation := range paths {
		assert.True(t, Exists(filepath.Join(baseRestoreDir, relLocation, Current)))
	}
}

func openDB(t *testing.T, dir string) *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, dir)
	assert.NoError(t, err)
	return db
}

func WriteTestDB(t *testing.T, dir string) {
	// create directory even if the file is not present
	err := os.MkdirAll(dir, os.ModePerm)
	assert.NoError(t, err)

	db := openDB(t, dir)
	wo := gorocksdb.NewDefaultWriteOptions()
	err = db.Put(wo, []byte("foo1"), []byte("bar"))
	assert.NoError(t, err)
	err = db.Put(wo, []byte("foo2"), []byte("bar"))
	assert.NoError(t, err)
	db.Close()
}
