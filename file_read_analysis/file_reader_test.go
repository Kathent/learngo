package file_read_analysis

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewFilesReaderWithDir(t *testing.T) {
	reader, err := NewFilesReaderWithDir("C:\\GoPath\\src\\learngo")
	assert.Nil(t, err)

	var bt = make([]byte, 1)
	count, err := reader.Read(bt)
	assert.Nil(t, err)
	assert.Equal(t, count, 1)

	reader.Close()

	count, err = reader.Read(bt)
	assert.NotNil(t, err)
}

func TestNewFilesReaderWithDir2(t *testing.T) {
	reader, err := NewFilesReaderWithDir("C:\\GoPath\\src\\learngo")
	assert.Nil(t, err)

	var bt = make([]byte, 10240)
	count, err := reader.Read(bt)
	assert.Nil(t, err)
	t.Logf("count:%d, bt:%s", count, string(bt))
}
