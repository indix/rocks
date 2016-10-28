package repair

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ind9/rocks/cmd/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/tecbot/gorocksdb"
)

func asTestRepairOnNormalDB(t *testing.T) {
	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	defer os.RemoveAll(dataDir)
	assert.NoError(t, err)

	testutils.WriteTestDB(t, dataDir)
	err = DoRepair(dataDir)
	assert.NoError(t, err)
}

func TestRepairOnCustomSSTFiles(t *testing.T) {
	db1DataDir, err := ioutil.TempDir("", "ind9-rocks")
	assert.NoError(t, err)
	fmt.Printf("%s\n", db1DataDir)
	// defer os.RemoveAll(db1DataDir)
	testutils.WriteTestDB(t, db1DataDir)
	assert.NoError(t, DoRepair(db1DataDir))

	db2DataDir, err := ioutil.TempDir("", "ind9-rocks")
	assert.NoError(t, err)
	fmt.Printf("%s\n", db2DataDir)
	// defer os.RemoveAll(db2DataDir)
	testutils.WriteTestDB1(t, db2DataDir)
	assert.NoError(t, DoRepair(db2DataDir))
}

func asTestOpenDB(t *testing.T) {
	dir := "/tmp/ind9-rocks532690606"
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(false)
	stringop := &testutils.StringConcatMergeOp{}
	opts.SetMergeOperator(stringop)
	db, err := gorocksdb.OpenDb(opts, dir)
	assert.NoError(t, err)
	defer db.Close()

	keyRange := gorocksdb.Range{}
	db.CompactRange(keyRange)

	ro := gorocksdb.NewDefaultReadOptions()
	it := db.NewIterator(ro)
	defer it.Close()
	it.SeekToFirst()
	for ; it.Valid(); it.Next() {
		fmt.Printf("Key: %v Value: %v\n", string(it.Key().Data()), string(it.Value().Data()))
	}
	assert.NoError(t, it.Err())
}
