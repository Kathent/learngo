package file_read_analysis

import (
	"os"
	"errors"
	"io/ioutil"
	"path/filepath"
	"io"
)

type FilesReader struct {
	files []*os.File
	index int
}

func NewFilesReader(files []*os.File) *FilesReader{
	return &FilesReader{files: files}
}

func NewFilesReaderWithPaths(files []string) (*FilesReader, error){
	tmp := make([]*os.File, 0)
	var suc bool = true
	var err error
	for _, v := range files{
		file, err := os.Open(v)
		if err != nil {
			suc = false
			break
		}

		tmp = append(tmp, file)
	}

	if !suc{
		for _, v := range tmp {
			v.Close()
		}
	}

	return &FilesReader{files:tmp}, err
}

func NewFilesReaderWithDir(path string) (*FilesReader, error){
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("path is not dir.")
	}

	infos, err := ioutil.ReadDir(path)
	tmp := make([]string, 0)
	for _, info := range infos {
		if info.IsDir(){
			continue
		}
		tmp = append(tmp, filepath.Join(path, info.Name()))
	}

	return NewFilesReaderWithPaths(tmp)
}

func (r *FilesReader) Read(b []byte) (int, error) {
	if len(b) <= 0 {
		return 0, nil
	}

	n, err := r.files[r.index].Read(b)
	if err != nil {
		return 0, err
	}

	if n < len(b) {
		if r.index == len(r.files) - 1 {
			return n, nil
		}

		r.index++
		newCount, err := r.Read(b[n:])
		return n + newCount, err
	}else {
		return n, nil
	}
}

func (r *FilesReader) Close() (err error) {
	for _, f := range r.files {
		closeErr := f.Close()
		if closeErr != nil {
			err = closeErr
		}
	}
	return
}

func (r *FilesReader) ReadByte() (byte, error) {
	var bt = make([]byte, 1)
	_, err := r.files[r.index].Read(bt)
	if err != nil && err == io.EOF{
		if r.index == len(r.files) - 1{
			return 0, err
		}

		r.index++
		return r.ReadByte()
	}

	return bt[0], err
}
