package ops

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tecbot/gorocksdb"
)

func TestRestore(t *testing.T) {
	backupDir, err := ioutil.TempDir("", "ind9-rocks-backup")
	defer os.RemoveAll(backupDir)
	assert.NoError(t, err)
	restoreDir, err := ioutil.TempDir("", "ind9-rocks-restore")
	defer os.RemoveAll(restoreDir)
	assert.NoError(t, err)
	dir, err := ioutil.TempDir("", "ind9-rocks")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	db := createDummyDB(t, dir)
	doBackup(t, backupDir, db)
	db.Close()

	err = DoRestore(backupDir, restoreDir, restoreDir)
	assert.NoError(t, err)

	db = openDB(t, restoreDir)
	ro := gorocksdb.NewDefaultReadOptions()
	value, err := db.GetBytes(ro, []byte("foo1"))
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(value))

	value, err = db.GetBytes(ro, []byte("foo2"))
	assert.NoError(t, err)
	assert.Equal(t, "bar", string(value))
}

func createDummyDB(t *testing.T, dir string) *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, dir)
	assert.NoError(t, err)
	wo := gorocksdb.NewDefaultWriteOptions()
	db.Put(wo, []byte("foo1"), []byte("bar"))
	db.Put(wo, []byte("foo2"), []byte("bar"))
	return db
}

func openDB(t *testing.T, dir string) *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	db, err := gorocksdb.OpenDb(opts, dir)
	assert.NoError(t, err)
	return db
}

func doBackup(t *testing.T, backupDir string, db *gorocksdb.DB) {
	opts := gorocksdb.NewDefaultOptions()
	backup, err := gorocksdb.OpenBackupEngine(opts, backupDir)
	assert.NoError(t, err)
	err = backup.CreateNewBackup(db)
	assert.NoError(t, err)
}
