package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/yesseneon/todo"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}

type TodoList interface {
	Create(userID int, list todo.TodoList) (int, error)
	GetAll(userID int) ([]todo.TodoList, error)
	GetByID(userID, listID int) (todo.TodoList, error)
	Update(userID, listID int, list todo.TodoListUpdate) error
	Delete(userID, listID int) error
}

type TodoItem interface {
	Create(listID int, item todo.TodoItem) (int, error)
	GetAll(userID, listID int) ([]todo.TodoItem, error)
	GetByID(userID, itemID int) (todo.TodoItem, error)
	Update(userID, itemID int, item todo.TodoItemUpdate) error
	Delete(userID, itemID int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthRepository(db),
		TodoList:      newTodoListRepository(db),
		TodoItem:      newTodoItemRepository(db),
	}
}
