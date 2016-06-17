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
	assert.NoError(t, err)
	defer os.RemoveAll(backupDir)
	restoreDir, err := ioutil.TempDir("", "ind9-rocks-restore")
	assert.NoError(t, err)
	defer os.RemoveAll(restoreDir)
	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	assert.NoError(t, err)
	defer os.RemoveAll(dataDir)

	db := openDB(t, dataDir)
	wo := gorocksdb.NewDefaultWriteOptions()
	db.Put(wo, []byte("foo1"), []byte("bar"))
	db.Put(wo, []byte("foo2"), []byte("bar"))
	db.Close()
	err = DoBackup(dataDir, backupDir)
	assert.NoError(t, err)

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

func openDB(t *testing.T, dir string) *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, dir)
	assert.NoError(t, err)
	return db
}
