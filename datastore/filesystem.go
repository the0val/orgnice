package datastore

import (
	"os"
)

var orgniceRoot string

// CreateFolder creates a directory at specified path with proper permissions.
func CreateFolder(path string) error {
	err := os.Mkdir(path, os.ModeDir+os.ModePerm)
	if err != nil {
		fileInfo, e := os.Stat(path)
		fileMode := fileInfo.Mode()
		if e == nil && fileMode.IsDir() {
			if fileMode.Perm()&0700 != 0700 {
				return os.Chmod(path, fileMode|0700)
			}
			return nil
		}
		return err
	}
	return nil
}
