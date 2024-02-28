package service

import (
	"github.com/SvetlanaGrin/todo-app"
	repoMock "github.com/SvetlanaGrin/todo-app/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTodoListService_Create(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	services := NewTodoListService(repos_List)

	out_todoList := todo.TodoList{
		Id:          1,
		Title:       "Купить одежду",
		Description: "Блузку",
	}

	var userId int = 1

	repos_List.EXPECT().Create(userId, out_todoList).Return(out_todoList.Id, nil)
	output, err := services.Create(userId, out_todoList)
	require.NoError(t, err)
	require.Equal(t, out_todoList.Id, output)
}

func TestTodoListService_GetAll(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	services := NewTodoListService(repos_List)

	out_todoList := []todo.TodoList{}
	out_todoList1 := todo.TodoList{
		Id:          1,
		Title:       "Купить одежду",
		Description: "Блузку",
	}
	out_todoList2 := todo.TodoList{
		Id:          2,
		Title:       "Купить одежду",
		Description: "Блузку",
	}
	out_todoList = append(out_todoList, out_todoList1)
	out_todoList = append(out_todoList, out_todoList2)

	var userId int = 1

	repos_List.EXPECT().GetAll(userId).Return(out_todoList, nil)
	output, err := services.GetAll(userId)
	require.NoError(t, err)
	require.Equal(t, out_todoList, output)
}

func TestTodoListService_GetById(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	services := NewTodoListService(repos_List)

	out_todoList := todo.TodoList{
		Id:          1,
		Title:       "Купить одежду",
		Description: "Блузку",
	}

	var userId int = 1
	var listId int = 1
	repos_List.EXPECT().GetById(userId, listId).Return(out_todoList, nil)
	output, err := services.GetById(userId, listId)
	require.NoError(t, err)
	require.Equal(t, out_todoList, output)
}

func TestTodoListService_Delete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	services := NewTodoListService(repos_List)

	var userId int = 1
	var listId int = 1
	repos_List.EXPECT().Delete(userId, listId).Return(nil)
	err := services.Delete(userId, listId)
	require.NoError(t, err)
}

func TestTodoListService_Update(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	services := NewTodoListService(repos_List)

	var title *string
	title1 := "Купить одежду"
	title = &title1
	var description *string
	description1 := "Сегодня"
	description = &description1

	out_todoList := todo.UpdateListInput{
		Title:       title,
		Description: description,
	}

	var userId int = 1
	var listId int = 1
	repos_List.EXPECT().Update(userId, listId, out_todoList).Return(nil)
	err := services.Update(userId, listId, out_todoList)
	require.NoError(t, err)
}

func TestTodoListService_UpdateError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	services := NewTodoListService(repos_List)

	var title *string

	title = nil
	var description *string

	description = nil

	out_todoList := todo.UpdateListInput{
		Title:       title,
		Description: description,
	}

	var userId int = 1
	var listId int = 1
	err := services.Update(userId, listId, out_todoList)
	require.EqualError(t, err, "update structure has no values")
}
