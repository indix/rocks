package ops

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsitency(t *testing.T) {
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
	var check bool

	err = DoBackup(dataDir, backupDir)
	assert.NoError(t, err)
	assert.True(t, Exists(filepath.Join(backupDir, LatestBackup)))

	err = DoRestore(backupDir, restoreDir, restoreDir, false)
	assert.NoError(t, err)
	assert.True(t, Exists(filepath.Join(restoreDir, Current)))

	check, err = DoConsistency(dataDir, restoreDir)
	assert.NoError(t, err)
	assert.True(t, check)
}
