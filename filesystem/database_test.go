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

func TestUseProjects(t *testing.T) {
	user, _ := InitDb("test.db")
	defer os.Remove("test.db")

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

	user.NewProject("Project 2")
	if res, _ := user.SearchProjects(""); len(res) != 2 {
		fmt.Println("Should find two projects")
		t.Fail()
	}
}

func TestUseTasks(t *testing.T) {
	user, _ := InitDb("test.db")
	defer os.Remove("test.db")

	if res, err := user.AllTasks(); err != nil || len(res) != 0 {
		fmt.Println("Should find 0 tasks and give non-nil error code")
		fmt.Println(err.Error())
		t.Fail()
	}

	for i := 0; i < 10; i++ {
		_, err := user.NewTask("Task"+string(i), 0)
		if err != nil {
			fmt.Println("Should add task")
			t.FailNow()
		}
	}
	res, _ := user.AllTasks()
	if len(res) != 10 {
		fmt.Println("Should find 10 tastks")
		t.Fail()
	}
	for _, task := range res {
		if newTask, _ := user.TaskFromID(task.ID); newTask != task {
			fmt.Println("Should find tasks with TaskFromID")
			t.Fail()
		}
	}

	task, _ := user.NewTask("Name1", 0)
	task.Name = "New name"
	task.Done = true
	err := user.StoreTask(task)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	updatedTask, _ := user.TaskFromID(task.ID)
	if updatedTask != task {
		fmt.Println("Should update task")
		t.Fail()
	}

	err = user.StoreTask(Task{task.ID + 1, "Stored task", task.Project, false})
	if err != nil {
		fmt.Println("Should create new task with StoreTask")
		t.Fail()
	}
}
