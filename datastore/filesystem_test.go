package datastore

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestCreateDir(t *testing.T) {
	const tPath = "./test"
	os.RemoveAll(tPath)

	err := createFolder(tPath)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	os.Remove(tPath)

	os.Remove(tPath)
	os.Mkdir(tPath, os.ModeDir)
	err = createFolder(tPath)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	tPath2 := path.Join(tPath, ".subdir")
	err = createFolder(tPath2)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	os.RemoveAll(tPath)

	os.Mkdir(tPath, os.ModeDir+os.ModePerm)
	err = createFolder(tPath)
	if err != nil {
		t.Fail()
	}
	os.Remove(tPath)

	f, err := os.Create(tPath)
	if err != nil {
		t.Error(err)
	}
	f.Close()
	err = createFolder(tPath)
	if err == nil {
		t.Fail()
	}
	os.Remove(tPath)
}
