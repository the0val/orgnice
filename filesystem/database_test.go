package filesystem

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func initializeDatabase() {

}

func TestInitDb(t *testing.T) {
	testFile := "./test.db"

	// Standard path, file does not exist
	os.Remove(testFile)
	if _, err := InitDb(testFile); err != nil {
		t.Fail()
	}

	// Standard path, file does exist
	if _, err := InitDb(testFile); err != nil {
		t.Fail()
	}

	os.Remove(testFile)
	f, err := os.Create(testFile)
	if err != nil {
		t.Error(err)
	}
	f.Write([]byte("Hello"))
	f.Close()
	// File exists but not valid database, should fail
	if _, err := InitDb(testFile); err == nil {
		t.Fail()
	}

	os.Remove(testFile)
	testDir := "testDir"
	os.Mkdir(testDir, os.ModeDir+0500)
	if _, err := InitDb(path.Join(testDir, testFile)); err == nil {
		t.Fail()
	}
	os.RemoveAll(testDir)
}

func TestUseDatabase(t *testing.T) {
	// TODO custom init function for tests
	user, _ := InitDb("test.db")

	if _, err := user.SearchProjects(""); err != nil {
		fmt.Println("Should return all entries when searching with empty string")
		t.Fail()
	}

	if res, err := user.SearchProjects("NotInDB"); err != nil && len(res) == 0 {
		fmt.Println("Should return nil error code when nothing found")
		t.Fail()
	}

	if res, _ := user.SearchProjects("b"); len(res) != 1 {
		fmt.Println("Should fuzzy match names")
		t.Fail()
	}

	if res, _ := user.ProjectFromID(0); res.Name != "Inbox" {
		fmt.Println("Should find Inbox at ID 0")
		t.Fail()
	}
}
