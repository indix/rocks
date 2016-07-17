package statistics

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ind9/rocks/cmd/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStatistics(t *testing.T) {
	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	defer os.RemoveAll(dataDir)
	assert.NoError(t, err)
	testutils.WriteTestDB(t, dataDir)
	count, err := DoStats(dataDir)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func TestRecursiveStatistics(t *testing.T) {
	baseDataDir, err := ioutil.TempDir("", "baseDataDir")
	err = os.MkdirAll(baseDataDir, os.ModePerm)
	defer os.RemoveAll(baseDataDir)
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

	count, err := DoRecursiveStats(baseDataDir, 3)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), count)
}
