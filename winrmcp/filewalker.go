package winrmcp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/masterzen/winrm/winrm"
)

type FileWalker struct {
	client  *winrm.Client
	config  *Config
	toDir   string
	fromDir string
}

func (fw *FileWalker) copyFile(fromPath string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !shouldUploadFile(fi) {
		return nil
	}

	hostPath, _ := filepath.Abs(fromPath)
	fromDir, _ := filepath.Abs(fw.fromDir)
	relPath, _ := filepath.Rel(fromDir, hostPath)
	toPath := filepath.Join(fw.toDir, relPath)

	f, err := os.Open(hostPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't read file %s: %v", fromPath, err))
	}

	return doCopy(fw.client, fw.config, f, winPath(toPath))
}

func shouldUploadFile(fi os.FileInfo) bool {
	// Ignore dir entries and OS X special hidden file
	return !fi.IsDir() && ".DS_Store" != fi.Name()
}
