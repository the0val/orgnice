package filesystem

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	// Runs to give a database handler
	_ "github.com/mattn/go-sqlite3"
)

// User represents all the data about a user. It's stored in a database file.
type User struct {
	db *sql.DB
}

// Task is a task stored in the database
type Task struct {
	ID      int
	Name    string
	Project Project
	Done    bool
}

// Project is a project stored in the database
type Project struct {
	ID   int
	Name string
}

// InitDb will create a database at path if it doesn't exist.
func InitDb(path string) (User, error) {
	user := User{}
	// If the database doesn't exist, should create tables
	// after opening the database.
	createTables := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		createTables = true
	}

	var err error
	user.db, err = sql.Open("sqlite3", path)
	if err != nil {
		return user, err
	}
	if e := user.db.Ping(); e != nil {
		return user, e
	}

	if createTables {
		if err := user.createTables(); err != nil {
			return user, err
		}
	}

	return user, nil
}

// NewProject creates a new project with given name
// in the database db.
func (user *User) NewProject(name string) error {
	tx, err := user.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO projects (name) VALUES (?)", name)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SearchProjects searches the database for projects with name that
// contains the given string (case-insensitive).
func (user *User) SearchProjects(name string) ([]Project, error) {
	rows, err := user.db.Query("SELECT id, name FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Project, 0)
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			fmt.Println(err)
			return nil, err
		}
		if strings.Contains(p.Name, name) {
			out = append(out, p)
		}
	}

	return out, nil
}

// ProjectFromID returns a project from the database with the given ID
// If no match found error will be sql.ErrNoRows
func (user *User) ProjectFromID(ID int) (Project, error) {
	row := user.db.QueryRow("SELECT id, name FROM projects WHERE id=?", ID)
	out := Project{}
	if err := row.Scan(&out.ID, &out.Name); err != nil {
		return Project{}, err
	}
	return out, nil
}

// NewTask creates a new task with the given projectID.
// Use projcetID 0 to put it in the default location Inbox.
func (user *User) NewTask(name string, projectID int) (Task, error) {
	p, err := user.ProjectFromID(projectID)
	if err != nil {
		return Task{}, err
	}

	res, err := user.db.Exec("INSERT INTO tasks (name, project) VALUES (?, ?)", name, projectID)
	if err != nil {
		return Task{}, err
	}
	taskID, _ := res.LastInsertId()

	return Task{ID: int(taskID), Name: name, Project: p}, nil
}

// TaskFromID returns a task from the database with the given ID
// If no match found error will be sql.ErrNoRows
func (user *User) TaskFromID(ID int) (Task, error) {
	row := user.db.QueryRow("SELECT id, name, project, done FROM tasks WHERE id=?", ID)
	out := Task{}
	var doneInt, projectID int
	if err := row.Scan(&out.ID, &out.Name, &projectID, &doneInt); err != nil {
		return Task{}, err
	}
	out.interpretDatabase(user, doneInt, projectID)
	return out, nil
}

// TasksInProject returns a list of all tasks in project
func (user *User) TasksInProject(projectID int) ([]Task, error) {
	if _, err := user.ProjectFromID(projectID); err != nil {
		return []Task{}, err
	}
	rows, err := user.db.Query("SELECT id, name, done FROM tasks WHERE project=?", projectID)
	if err != nil {
		return []Task{}, err
	}
	defer rows.Close()

	out := make([]Task, 0)
	for rows.Next() {
		var task Task
		var doneInt int
		if err := rows.Scan(&task.ID, &task.Name, &doneInt); err != nil {
			return []Task{}, err
		}
		task.interpretDatabase(user, doneInt, projectID)
		out = append(out, task)
	}
	return out, nil
}

// AllTasks returns a list of all tasks belonging to user
func (user *User) AllTasks() ([]Task, error) {
	rows, err := user.db.Query("SELECT id, name, done, project FROM tasks ORDER BY id")
	if err != nil {
		return nil, err
	}
	out := make([]Task, 0)
	for rows.Next() {
		var task Task
		var doneInt int
		var projectID int
		if err := rows.Scan(&task.ID, &task.Name, &doneInt, &projectID); err != nil {
			return []Task{}, err
		}
		task.interpretDatabase(user, doneInt, projectID)
		out = append(out, task)
	}
	return out, nil
}

// StoreTask takes a task as parameter and stores it in the database.
// If a task with the same ID already exists it will update that instead.
// For creating a new task, NewTask is preferred.
func (user *User) StoreTask(task Task) error {
	if _, err := user.TaskFromID(task.ID); err != nil {
		_, err := user.db.Exec("INSERT INTO tasks (id, name, project) VALUES (?, ?, ?)", task.ID, task.Name, task.Project.ID)
		if err != nil {
			return err
		}
	} else {
		doneInt := boolToInt(task.Done)
		projectID := task.Project.ID
		user.db.Exec("UPDATE tasks SET name = ?, project = ?, done = ? WHERE id = ?", task.Name, projectID, doneInt, task.ID)
	}
	return nil
}

// createTables should be called whenever a new user database
// is created.
func (user *User) createTables() error {
	sqlStmt := `CREATE TABLE tasks (
		id INTEGER PRIMARY KEY,
		name STRING NOT NULL,
		project INTEGER DEFAULT 0,
		done INTEGER DEFAULT 0
	)`
	_, err := user.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `CREATE TABLE projects (
		id INTEGER PRIMARY KEY,
		name STRING NOT NULL
	)`
	_, err = user.db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	// Create project Inbox, used as default for tasks
	user.db.Exec("INSERT INTO projects (id, name) VALUES (0, 'Inbox')")

	return nil
}

// intepretDatabase takes fields from the tasks database table
// and sets the corresponding feilds in task
func (task *Task) interpretDatabase(db *User, doneInt, projectID int) {
	project, _ := db.ProjectFromID(projectID)
	task.Project = project
	task.Done = intToBool(doneInt)
}

// intToBool takes an int as input and return false
// if the int is 0. Otherwise true.
func intToBool(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

// boolToInt takes a bool and returns 1 if it's true,
// 0 if it's false.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
