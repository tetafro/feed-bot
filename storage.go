package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Storage describes persistent datastorage.
type Storage interface {
	Get() ([]int64, error)
	Save([]int64) error
}

// FileStorage is a storage that uses plain text file for storing data.
type FileStorage struct {
	file string
}

// NewFileStorage creates new file storage.
func NewFileStorage(file string) (*FileStorage, error) {
	fs := &FileStorage{file: file}

	f, err := os.Create(fs.file)
	if os.IsExist(err) {
		return fs, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "create file")
	}
	f.Close()

	return fs, nil
}

// Get gets list of numbers from file with comma-separated data.
func (s *FileStorage) Get() ([]int64, error) {
	b, err := ioutil.ReadFile(s.file)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}
	if len(b) == 0 {
		return nil, nil
	}

	// Parse comma-separated numbers
	var nums []int64
	for _, str := range strings.Split(string(b), ",") {
		n, err := strconv.Atoi(str)
		if err != nil {
			return nil, errors.Wrap(err, "convert to integer")
		}
		nums = append(nums, int64(n))
	}

	return nums, nil
}

// Save saves list of numbers in comma-separated format to file.
func (s *FileStorage) Save(nums []int64) error {
	var str string
	if len(nums) > 0 {
		for _, n := range nums {
			str += strconv.Itoa(int(n)) + ","
		}
		str = str[:len(str)-1]
	}
	return ioutil.WriteFile(s.file, []byte(str), 0600)
}
