package internal

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type TodoPersistence interface {
	Create(todo *Todo) (*Todo, error)
	GetAll() ([]*Todo, error)
	GetByID(id string) (*Todo, error)
	Update(todo *Todo) (*Todo, error)
}

func InitializePersistence() {
	db, err := OpenConnection()
	if err != nil {
		panic(err)
	}
	defer CloseConnection(db)
	db.Exec("CREATE TABLE IF NOT EXISTS todos (id TEXT PRIMARY KEY, title TEXT, description TEXT, status TEXT)")
}

func OpenConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "todos.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseConnection(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}

type TodosDb struct {
	OpenConnection  func() (*sql.DB, error)
	CloseConnection func(db *sql.DB) error
}

func (t *TodosDb) Create(todo *NewTodoRequest) (*Todo, error) {
	db, err := t.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer t.CloseConnection(db)

	stmt, err := db.Prepare("INSERT INTO todos (title, description, status) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(todo.Title, todo.Description, todo.Status)
	if err != nil {
		return nil, err
	}

	return &Todo{}, nil
}

func (t *TodosDb) GetAll() ([]*Todo, error) {
	db, err := t.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer t.CloseConnection(db)

	rows, err := db.Query("SELECT * FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []*Todo{}
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func (t *TodosDb) GetByID(id string) (*Todo, error) {
	db, err := t.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer t.CloseConnection(db)

	stmt, err := db.Prepare("SELECT * FROM todos WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	todo := &Todo{}
	err = row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (t *TodosDb) Update(todo *Todo) (*Todo, error) {
	db, err := t.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer t.CloseConnection(db)

	stmt, err := db.Prepare("UPDATE todos SET title = ?, description = ?, status = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(todo.Title, todo.Description, todo.Status, todo.ID)
	if err != nil {
		return nil, err
	}

	return todo, nil

}
