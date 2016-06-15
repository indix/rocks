package ops

import (
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
