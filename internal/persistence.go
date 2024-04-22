package internal

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	_ "github.com/mattn/go-sqlite3"
)

func InitializePersistence(logger *zerolog.Logger) {
	db, err := OpenConnection()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to open connection to database")
		panic(err)
	}
	defer CloseConnection(db)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (id TEXT PRIMARY KEY, title TEXT, description TEXT, status TEXT, severity TEXT)`)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create table")
		panic(err)
	}
	logger.Info().Msg("Database initialized")
}

func OpenConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "todos.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CloseConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
}

func NewTodosDb(logger *zerolog.Logger) TodosDb {
	return TodosDb{
		OpenConnection:  OpenConnection,
		CloseConnection: CloseConnection,
		logger:          logger,
	}
}

func (t *TodosDb) Create(todo *NewTodoRequest) (*Todo, error) {
	db, err := t.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer t.CloseConnection(db)
	uuidStr := uuid.New().String()
	stmt, err := db.Prepare("INSERT INTO todos (id, title, description, status, severity) VALUES (?, ?, ?, ?,?)")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	_, err = stmt.Exec(uuidStr, todo.Title, todo.Description, todo.Status, todo.Severity)
	if err != nil {
		t.logger.Err(err).Msg("Failed to insert todo")
		return nil, err
	}

	return &Todo{
		ID:          uuidStr,
		Title:       todo.Title,
		Description: todo.Description,
		Status:      todo.Status,
		Severity:    todo.Severity,
	}, nil
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
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var todos []*Todo
	for rows.Next() {
		todo := &Todo{}
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Severity)
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
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	row := stmt.QueryRow(id)
	todo := &Todo{}
	err = row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Severity)
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
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	_, err = stmt.Exec(todo.Title, todo.Description, todo.Status, todo.ID)
	if err != nil {
		return nil, err
	}

	return todo, nil

}
