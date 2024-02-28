package handler

import (
	"bytes"
	"encoding/json"
	"github.com/SvetlanaGrin/todo-app"
	repoMock "github.com/SvetlanaGrin/todo-app/pkg/mocks"
	"github.com/SvetlanaGrin/todo-app/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetItemById(t *testing.T) {
	TestAuthHandler_signIn(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)
	var tok []string
	tok = append(tok, "Basic "+GlobalToken)

	r := gin.Default()

	out_todoItem := todo.TodoItem{
		Id:          0,
		Title:       "Купить одежду",
		Description: "Блузку",
		Done:        false,
	}

	var userId int = 1
	var listId int = 1
	repoItem.EXPECT().GetById(userId, listId).Return(out_todoItem, nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r.GET("/api/items/:id", hand.userIdentity, hand.getItemById)
	reqGet := httptest.NewRequest(http.MethodGet,
		"/api/items/1",
		bytes.NewBuffer(
			[]byte(`{"title": "Купить одежду","description": "Блузку"}`)))
	reqGet.Header["Authorization"] = tok

	w := httptest.NewRecorder()

	r.ServeHTTP(w, reqGet)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	var output todo.TodoItem
	if err := json.Unmarshal(data, &output); err != nil {
		log.Fatal(err.Error())
	}
	require.NoError(t, err)
	require.Equal(t, out_todoItem, output)
}

func TestHandler_CreateItem(t *testing.T) {
	TestAuthHandler_signIn(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)
	var tok []string
	tok = append(tok, "Basic "+GlobalToken)

	r := gin.Default()

	out_todoItem := todo.TodoItem{
		Id:          0,
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
	repoList.EXPECT().GetById(userId, listId).Return(out_todoList, nil)
	repoItem.EXPECT().Create(listId, out_todoItem).Return(out_todoItem.Id, nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r.POST("/api/lists/:id/items", hand.userIdentity, hand.createItem)
	reqGet := httptest.NewRequest(http.MethodPost,
		"/api/lists/1/items",
		bytes.NewBuffer(
			[]byte(`{"title": "Купить одежду","description": "В ТЦ"}`)))
	reqGet.Header["Authorization"] = tok

	w := httptest.NewRecorder()

	r.ServeHTTP(w, reqGet)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, []byte("{\"id\":0}"), data)
}

func TestHandler_GetAllItem(t *testing.T) {
	TestAuthHandler_signIn(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)
	var tok []string
	tok = append(tok, "Basic "+GlobalToken)

	r := gin.Default()

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

	repoItem.EXPECT().GetAll(userId, listId).Return(out_todoItem, nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r.GET("/api/lists/:id/items", hand.userIdentity, hand.getAllItem)
	reqGet := httptest.NewRequest(http.MethodGet,
		"/api/lists/1/items", nil)
	reqGet.Header["Authorization"] = tok

	w := httptest.NewRecorder()

	r.ServeHTTP(w, reqGet)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	var output []todo.TodoItem
	if err := json.Unmarshal(data, &output); err != nil {
		log.Fatal(err.Error())
	}

	require.NoError(t, err)
	require.Equal(t, out_todoItem, output)
}

func TestHandler_DeleteItem(t *testing.T) {
	TestAuthHandler_signIn(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)
	var tok []string
	tok = append(tok, "Basic "+GlobalToken)

	r := gin.Default()

	var userId int = 1
	var listId int = 1
	repoItem.EXPECT().Delete(userId, listId).Return(nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r.DELETE("/api/items/:id", hand.userIdentity, hand.deleteItem)
	reqGet := httptest.NewRequest(http.MethodDelete,
		"/api/items/1", nil)
	reqGet.Header["Authorization"] = tok

	w := httptest.NewRecorder()

	r.ServeHTTP(w, reqGet)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.Equal(t, []byte("{\"Status\":\"ok\"}"), data)
	require.NoError(t, err)
}

func TestHandler_UpdateItem(t *testing.T) {
	TestAuthHandler_signIn(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)
	var tok []string
	tok = append(tok, "Basic "+GlobalToken)

	r := gin.Default()

	var title *string
	title1 := "Купить одежду"
	title = &title1
	var description *string
	description1 := "Сегодня"
	description = &description1
	var done *bool
	done1 := false
	done = &done1

	out_todoItem := todo.UpdateItemInput{
		Title:       title,
		Description: description,
		Done:        done,
	}

	var userId int = 1
	var itemId int = 1

	repoItem.EXPECT().Update(userId, itemId, out_todoItem).Return(nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r.PUT("/api/items/:id", hand.userIdentity, hand.updateItem)
	reqGet := httptest.NewRequest(http.MethodPut,
		"/api/items/1",
		bytes.NewBuffer(
			[]byte(`{"title": "Купить одежду","description": "Сегодня","done": false }`)))
	reqGet.Header["Authorization"] = tok

	w := httptest.NewRecorder()

	r.ServeHTTP(w, reqGet)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.Equal(t, []byte("{\"Status\":\"ok\"}"), data)
	require.NoError(t, err)
}
