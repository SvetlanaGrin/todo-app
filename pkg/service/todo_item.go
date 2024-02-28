package service

import (
	"github.com/SvetlanaGrin/todo-app"
	"github.com/SvetlanaGrin/todo-app/pkg/repository"
)

type TodoItemService struct {
	listrepo repository.TodoList
	repo     repository.TodoItem
}

func NewTodoItemService(repo repository.TodoItem, listrepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listrepo: listrepo}
}

func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listrepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item)

}
func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}
func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}
func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	return s.repo.Update(userId, itemId, input)
}
