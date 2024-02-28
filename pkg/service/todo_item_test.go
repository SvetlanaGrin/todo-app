package service

import (
	"errors"
	"github.com/SvetlanaGrin/todo-app"
	repoMock "github.com/SvetlanaGrin/todo-app/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTodoItemService_GetAll(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	repos_Item := repoMock.NewMockTodoItem(ctl)
	services := NewTodoItemService(repos_Item, repos_List)

	out_todoItem := []todo.TodoItem{}
	out_todoItem1 := todo.TodoItem{
		Id:          1,
		Title:       "Купить продукты",
		Description: "Срочно",
		Done:        false,
	}

	out_todoItem2 := todo.TodoItem{
		Id:          2,
		Title:       "Купить одежду",
		Description: "В ТЦ",
		Done:        false,
	}

	out_todoItem = append(out_todoItem, out_todoItem1)
	out_todoItem = append(out_todoItem, out_todoItem2)
	var userId int = 1
	var listId int = 1
	repos_Item.EXPECT().GetAll(userId, listId).Return(out_todoItem, nil)
	out_todoItem_service, err := services.GetAll(userId, listId)
	require.NoError(t, err)
	require.Equal(t, out_todoItem, out_todoItem_service)

}

func TestTodoItemService_GetById(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	repos_Item := repoMock.NewMockTodoItem(ctl)
	services := NewTodoItemService(repos_Item, repos_List)

	out_todoItem1 := todo.TodoItem{
		Id:          2,
		Title:       "Купить одежду",
		Description: "В ТЦ",
		Done:        false,
	}

	var userId int = 1
	var listId int = 1
	repos_Item.EXPECT().GetById(userId, listId).Return(out_todoItem1, nil)
	out_todoItem_service, err := services.GetById(userId, listId)
	require.NoError(t, err)
	require.Equal(t, out_todoItem1, out_todoItem_service)
}

func TestTodoItemService_Delete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	repos_Item := repoMock.NewMockTodoItem(ctl)
	services := NewTodoItemService(repos_Item, repos_List)

	var userId int = 1
	var listId int = 1
	repos_Item.EXPECT().Delete(userId, listId).Return(nil)
	err := services.Delete(userId, listId)
	require.NoError(t, err)
}

func TestTodoItemService_Create(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	repos_Item := repoMock.NewMockTodoItem(ctl)
	services := NewTodoItemService(repos_Item, repos_List)

	out_todoItem1 := todo.TodoItem{
		Id:          1,
		Title:       "Купить одежду",
		Description: "В ТЦ",
		Done:        false,
	}
	out_todoList := todo.TodoList{
		Id:          1,
		Title:       "Купить одежду",
		Description: "Блузку",
	}

	var userId int = 1
	var listId int = 1
	repos_List.EXPECT().GetById(userId, listId).Return(out_todoList, nil)
	repos_Item.EXPECT().Create(listId, out_todoItem1).Return(out_todoItem1.Id, nil)
	output, err := services.Create(userId, listId, out_todoItem1)
	require.NoError(t, err)
	require.Equal(t, out_todoItem1.Id, output)
}

func TestTodoItemService_CreateError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	repos_Item := repoMock.NewMockTodoItem(ctl)
	services := NewTodoItemService(repos_Item, repos_List)
	out_todoItem1 := todo.TodoItem{}
	out_todoList := todo.TodoList{}

	var userId int = 1
	var listId int = 1

	err := errors.New("Error")
	repos_List.EXPECT().GetById(userId, listId).Return(out_todoList, err)
	output, err := services.Create(userId, listId, out_todoItem1)
	require.EqualError(t, err, "Error")
	require.Equal(t, 0, output)
}

func TestTodoItemService_Update(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos_List := repoMock.NewMockTodoList(ctl)
	repos_Item := repoMock.NewMockTodoItem(ctl)
	services := NewTodoItemService(repos_Item, repos_List)
	var title *string
	title1 := "Купить одежду"
	title = &title1
	var description *string
	description1 := "Сегодня"
	description = &description1
	var done *bool
	done1 := false
	done = &done1
	out_todoItem1 := todo.UpdateItemInput{
		Title:       title,
		Description: description,
		Done:        done,
	}

	var userId int = 1
	var itemId int = 1

	repos_Item.EXPECT().Update(userId, itemId, out_todoItem1).Return(nil)
	err := services.Update(userId, itemId, out_todoItem1)
	require.NoError(t, err)
}
