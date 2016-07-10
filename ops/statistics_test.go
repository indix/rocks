package ops

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatistics(t *testing.T) {
	dataDir, err := ioutil.TempDir("", "ind9-rocks")
	defer os.RemoveAll(dataDir)
	assert.NoError(t, err)
	WriteTestDB(t, dataDir)
	err = DoStats(dataDir)
	assert.NoError(t, err)
}
