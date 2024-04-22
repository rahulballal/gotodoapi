package internal

import (
	"database/sql"
	"github.com/rs/zerolog"
)

type HandlerMap struct {
	Logger      *zerolog.Logger
	Persistence TodoPersistence
}

type Config struct {
	Port     uint64
	LogLevel zerolog.Level
}

type TodosDb struct {
	OpenConnection  func() (*sql.DB, error)
	CloseConnection func(db *sql.DB)
	logger          *zerolog.Logger
}

type TodoPersistence interface {
	Create(todo *NewTodoRequest) (*Todo, error)
	GetAll() ([]*Todo, error)
	GetByID(id string) (*Todo, error)
	Update(todo *Todo) (*Todo, error)
}
