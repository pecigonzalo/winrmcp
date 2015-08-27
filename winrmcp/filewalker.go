package winrmcp

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileWalker struct {
	fromDir string
}

func (fw *FileWalker) getFiles(fromPath string, fi os.FileInfo, err error) (*File, error) {
	if !shouldUploadFile(fi) {
		return nil, nil
	}
	var file *File

	file.path, _ = filepath.Abs(fromPath)
	// fromDir, _ := filepath.Abs(fw.fromDir)
	// relPath, _ := filepath.Rel(fromDir, hostPath)

	file.reader, err = os.Open(file.path)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read file %s: %v", fromPath, err)
	}

	//return doCopy(fw.client, fw.config, f, winPath(toPath))
	return file, err
}

func shouldUploadFile(fi os.FileInfo) bool {
	// Ignore dir entries and OS X special hidden file
	return !fi.IsDir() && ".DS_Store" != fi.Name()
}
