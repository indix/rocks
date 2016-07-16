package backup

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ind9/rocks/cmd/ops"
	"github.com/ind9/rocks/cmd/test-utils"
	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	defer os.RemoveAll(dataDir)
	assert.NoError(t, err)
	backupDir, err := ioutil.TempDir("", "ind9-rocks-backup")
	defer os.RemoveAll(backupDir)
	assert.NoError(t, err)

	testutils.WriteTestDB(t, dataDir)
	err = DoBackup(dataDir, backupDir)
	assert.NoError(t, err)

	assert.True(t, testutils.Exists(filepath.Join(backupDir, ops.LatestBackup)))
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
		testutils.WriteTestDB(t, filepath.Join(baseDataDir, relLocation))
	}

	err = DoRecursiveBackup(baseDataDir, baseBackupDir, 1)
	assert.NoError(t, err)

	for _, relLocation := range paths {
		assert.True(t, testutils.Exists(filepath.Join(baseBackupDir, relLocation, ops.LatestBackup)))
	}
}
