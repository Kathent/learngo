package file_read_analysis

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoadFilesWithReader(t *testing.T) {
	readErr := LoadFilesWithReader("C:\\GoPath\\src\\learngo\\load_file_dir", "2017-10-02 08:00:00", "2017-10-03 14:00:00")
	assert.Nil(t, readErr)
}
