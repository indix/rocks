package testutils

import (
	"fmt"
	"io"
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
	stringop := StringConcatMergeOp{}
	opts.SetMergeOperator(stringop)
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
	err = db.Merge(wo, []byte("foo1"), []byte("bar"))
	assert.NoError(t, err)
	err = db.Merge(wo, []byte("foo2"), []byte("bar"))
	assert.NoError(t, err)
	db.Close()
}

// WriteTestDB1 writes dummy data into a rocksdb store
func WriteTestDB1(t *testing.T, dir string) {
	// create directory even if the file is not present
	err := os.MkdirAll(dir, os.ModePerm)
	assert.NoError(t, err)

	db := OpenDB(t, dir)
	wo := gorocksdb.NewDefaultWriteOptions()
	err = db.Merge(wo, []byte("foo1"), []byte("bar1"))
	assert.NoError(t, err)
	err = db.Merge(wo, []byte("foo3"), []byte("bar"))
	assert.NoError(t, err)
	err = db.Merge(wo, []byte("foo4"), []byte("bar"))
	assert.NoError(t, err)
	db.Close()
}

// CopyFile copies data from src to dst
func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

// StringConcatMergeOp implements rocksdb's MergeOperator for merging keys
type StringConcatMergeOp struct {
	gorocksdb.MergeOperator
}

func (op StringConcatMergeOp) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	fmt.Printf("FullMerge called on %v", string(key))
	return existingValue, true
}

func (op StringConcatMergeOp) PartialMerge(key, leftOperand, rightOperand []byte) ([]byte, bool) {
	fmt.Printf("PartialMerge called")
	return rightOperand, true
}
func (op StringConcatMergeOp) Name() string {
	return "StringConcatMergeOp"
}
