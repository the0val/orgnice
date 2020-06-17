package filesystem

import (
	"log"
	"os"
	"path"
)

var orgniceRoot string

// InitDir creates the neccesary folder structure
// for the apps datastorage.
func InitDir() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find users home dir\n%v", err)
	}
	orgniceRoot = path.Join(home, ".orgnice")
	err = createFolder(orgniceRoot)
	if err != nil {
		log.Fatalf("Unable to initialize folder %v\t%v", orgniceRoot, err)
	}
}

func createFolder(path string) error {
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
