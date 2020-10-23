package manifest

import (
	"errors"
	"os"
	"path/filepath"
)

type file struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
}

type manifest struct {
	Files            []file `json:"files"`
	HashingAlgorithm string `json:"hashing_algorithm"`
	FolderCount      int    `json:"folder_count"`
	FileCount        int    `json:"file_count"`
	TotalSize        int64  `json:"total_size"`
}

// GetFilesFromDir test
func GetFilesFromDir(directoryPath string) ([]string, error) {
	var files []string

	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		return files, errors.New("Path does not exist")
	}

	root := directoryPath
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files, nil
}
