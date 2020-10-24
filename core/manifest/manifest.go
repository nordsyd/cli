package manifest

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nordsyd/cli/core/hash"
)

// File is a child of the Manifest struct
type File struct {
	Key  string `json:"key"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
}

// Manifest defines the manifest structure
type Manifest struct {
	Files            []File `json:"files"`
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

// GenerateManifest test
func GenerateManifest(files []string, rootPath string) Manifest {
	var (
		fileList    []File
		totalSize   int64
		fileCount   int
		folderCount int
	)

	for _, filePath := range files {
		fileStat, err := os.Stat(filePath)

		if err != nil {
			continue
		}

		if !fileStat.IsDir() {
			hash, err := hash.FileSHA256(filePath)

			if err != nil {
				fmt.Println("Error creating hash for: ", filePath)
				continue
			}

			fileCount++
			totalSize += fileStat.Size()

			fileStruct := File{
				Key:  strings.Replace(filePath, rootPath, "", 1),
				Size: fileStat.Size(),
				Hash: hash,
			}

			fileList = append(fileList, fileStruct)
		} else {
			folderCount++
		}
	}

	return Manifest{
		Files:            fileList,
		HashingAlgorithm: "sha256",
		FolderCount:      folderCount - 1,
		FileCount:        fileCount,
		TotalSize:        totalSize,
	}
}
