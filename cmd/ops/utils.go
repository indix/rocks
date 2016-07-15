package ops

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tecbot/gorocksdb"
)

// Exists function checks for existence of a file/directory
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// OpenDB creates a dummy rocksdb store
func OpenDB(t *testing.T, dir string) *gorocksdb.DB {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, dir)
	assert.NoError(t, err)
	return db
}

// WriteTestDB writes dummy data into a rocksdb store
func WriteTestDB(t *testing.T, dir string) {
	// create directory even if the file is not present
	err := os.MkdirAll(dir, os.ModePerm)
	assert.NoError(t, err)

	db := OpenDB(t, dir)
	wo := gorocksdb.NewDefaultWriteOptions()
	err = db.Put(wo, []byte("foo1"), []byte("bar"))
	assert.NoError(t, err)
	err = db.Put(wo, []byte("foo2"), []byte("bar"))
	assert.NoError(t, err)
	db.Close()
}
