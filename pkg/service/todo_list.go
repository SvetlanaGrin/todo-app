package service

import (
	"github.com/SvetlanaGrin/todo-app"
	"github.com/SvetlanaGrin/todo-app/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}
func (s *TodoListService) GetAll(userid int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userid)
}

func (s *TodoListService) GetById(userid, id int) (todo.TodoList, error) {
	return s.repo.GetById(userid, id)
}
func (s *TodoListService) Delete(userid, id int) error {
	return s.repo.Delete(userid, id)
}
func (s *TodoListService) Update(userid, id int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userid, id, input)
}
