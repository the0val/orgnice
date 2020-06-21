package filesystem

import (
	"database/sql"
	"os"

	// Runs to give a database handler
	_ "github.com/mattn/go-sqlite3"
)

// User represents all the data about a user. It's stored in a database file.
type User struct {
	*sql.DB
}

// Task is a task stored in the database
type Task struct {
	ID      int
	Name    string
	Project string
	Done    bool
}

// Project is a project stored in the databas
type Project struct {
	ID   int
	Name string
}

// InitDb will create a database at path if it doesn't exist.
func InitDb(path string) (User, error) {
	db := User{}
	var err error
	db.DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return db, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		db.createTables()
	}

	return db, nil
}

// NewProject creates a new project with given name
// in the database db.
func (db *User) NewProject(name string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO projects (name) VALUES (?)", name)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// NewTask creates a new task with the given projectID.
// Use projcetID 0 to put it in the default location Inbox.
func (db *User) NewTask(name string, projectID int) (Task, error) {
	p, err := db.ProjectFromID(projectID)
	if err != nil {
		return Task{}, err
	}

	tx, err := db.Begin()
	if err != nil {
		return Task{}, err
	}

	res, err := tx.Exec("INSERT INTO tasks (name, project) VALUES (?, ?)", name, projectID)
	if err != nil {
		return Task{}, err
	}
	taskID, _ := res.LastInsertId()

	return Task{ID: int(taskID), Name: name, Project: p.Name}, tx.Commit()
}

// SearchProjects searches the database for projects with name that
// contains the given string (case-insensitive).
func (db *User) SearchProjects(name string) ([]Project, error) {
	// https://pkg.go.dev/database/sql?tab=doc#example-DB.Query-MultipleResultSets
	rows, err := db.Query("SELECT id, name FROM projects WHERE name LIKE '%?%", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Project, 0)
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, p.Name); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

// ProjectFromID returns a project from the database with the given ID
// If no match found error will be sql.ErrNoRows
func (db *User) ProjectFromID(ID int) (Project, error) {
	row := db.QueryRow("SELECT id, name FROM projects WHERE id=?", ID)
	out := Project{}
	if err := row.Scan(&out.ID, &out.Name); err != nil {
		return Project{}, err
	}
	return out, nil
}

// TaskFromID returns a task from the database with the given ID
// If no match found error will be sql.ErrNoRows
func (db *User) TaskFromID(ID int) (Task, error) {
	row := db.QueryRow("SELECT id, name, project, done FROM tasks WHERE id=?", ID)
	out := Task{}
	var doneInt, projectID int
	if err := row.Scan(&out.ID, &out.Name, &projectID, &doneInt); err != nil {
		return Task{}, err
	}
	project, _ := db.ProjectFromID(projectID)
	out.Project = project.Name
	if doneInt == 0 {
		out.Done = false
	} else {
		out.Done = true
	}
	return out, nil
}

// TasksInProject returns a list of all tasks in project
func (db *User) TasksInProject(projectID int) ([]Task, error) {
	if _, err := db.ProjectFromID(projectID); err != nil {
		return []Task{}, err
	}
	rows, err := db.Query("SELECT id, name, done FROM tasks WHERE project=?", projectID)
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
		if doneInt == 0 {
			task.Done = false
		} else {
			task.Done = true
		}
		out = append(out, task)
	}
	return out, nil
}

func (db *User) createTables() error {
	sqlStmt := `CREATE TABLE tasks (
		id INTEGER PRIMARY KEY,
		name STRING NOT NULL,
		project INTEGER DEFAULT 0,
		done INTEGER
	)`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `CREATE TABLE projects (
		id INTEGER PRIMARY KEY,
		name STRING NOT NULL
	)`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	// Create project Inbox, used as default for tasks
	db.Exec("INSERT INTO projects (id, name) VALUES (0, 'Inbox')")

	return nil
}
