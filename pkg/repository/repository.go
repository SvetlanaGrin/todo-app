package repository

import (
	"github.com/SvetlanaGrin/todo-app"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUsers(user todo.User) (int, error)
	GetUser(username, password string) (todo.User, error)
}
type TodoList interface {
	Create(userId int, list todo.TodoList) (int, error)
	GetAll(userid int) ([]todo.TodoList, error)
	GetById(userid, id int) (todo.TodoList, error)
	Delete(userid, id int) error
	Update(userid, id int, input todo.UpdateListInput) error
}
type TodoItem interface {
	Create(listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(userId, itemId int) (todo.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input todo.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
