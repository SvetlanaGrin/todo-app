package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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

type DataOut struct {
	Data []todo.TodoList `json:"data"`
}

func TestHandler_GetAllList(t *testing.T) {
	TestAuthHandler_signIn(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)

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
	//repoAuth1 := repos.Authorization
	repoList.EXPECT().GetAll(userId).Return(out_todoList, nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	var tok []string
	tok = append(tok, "Basic "+GlobalToken)
	r := gin.Default()
	r.GET("/api/lists", hand.userIdentity, hand.getAllList)
	reqGet := httptest.NewRequest(http.MethodGet,
		"/api/lists", nil)
	reqGet.Header["Authorization"] = tok

	w := httptest.NewRecorder()

	r.ServeHTTP(w, reqGet)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	var output DataOut
	if err := json.Unmarshal(data, &output); err != nil {
		log.Fatal(err.Error())
	}
	require.NoError(t, err)
	require.Equal(t, out_todoList, output.Data)
}
func TestHandler_GetAllListError(t *testing.T) {
	TestAuthHandler_signIn(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	var userId int = 1
	err := errors.New("Error")

	repoList.EXPECT().GetAll(userId).Return(nil, err)

	var tok []string
	tok = append(tok, "Basic "+GlobalToken)

	r := gin.Default()
	r.GET("/api/lists", hand.userIdentity, hand.getAllList)
	reqGet := httptest.NewRequest(http.MethodGet,
		"/api/lists", nil)
	reqGet.Header["Authorization"] = tok

	w := httptest.NewRecorder()

	r.ServeHTTP(w, reqGet)

	require.Equal(t, w.Code, http.StatusInternalServerError)
	//if w.Code != http.StatusInternalServerError{
	//	t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	//}
	res := w.Result()

	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	//err:= errors.New("Error")
	require.Equal(t, data, []byte("{\"message\":\"Error\"}"))
	//require.Equal(t, out_todoList, output.Data)
}

func TestHandler_GetListById(t *testing.T) {
	TestAuthHandler_signIn(t)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)
	var tok []string
	tok = append(tok, "Basic "+GlobalToken)

	r := gin.Default()

	out_todoList := todo.TodoList{
		Id:          1,
		Title:       "Купить одежду",
		Description: "Блузку",
	}

	var userId int = 1
	var listId int = 1
	//repoAuth1 := repos.Authorization
	repoList.EXPECT().GetById(userId, listId).Return(out_todoList, nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r.GET("/api/lists/:id", hand.userIdentity, hand.getListById)
	reqGet := httptest.NewRequest(http.MethodGet,
		"/api/lists/1",
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

	var output todo.TodoList
	if err := json.Unmarshal(data, &output); err != nil {
		log.Fatal(err.Error())
	}
	require.NoError(t, err)
	require.Equal(t, out_todoList, output)
}

func TestHandler_CreateList(t *testing.T) {
	TestAuthHandler_signIn(t)

	var tok []string
	tok = append(tok, "Basic "+GlobalToken)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)

	out_todoList := todo.TodoList{
		Id:          0,
		Title:       "Купить одежду",
		Description: "Блузку",
	}

	var userId int = 1
	//repoAuth1 := repos.Authorization
	repoList.EXPECT().Create(userId, out_todoList).Return(out_todoList.Id, nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r := gin.Default()

	r.POST("/api/lists", hand.userIdentity, hand.createList)
	reqGet := httptest.NewRequest(http.MethodPost,
		"/api/lists", bytes.NewBuffer(
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
	require.NoError(t, err)
	require.Equal(t, data, []byte("{\"id\":0}"))
}

func TestHandler_DeleteList(t *testing.T) {
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
	//repoAuth1 := repos.Authorization
	repoList.EXPECT().Delete(userId, listId).Return(nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r.DELETE("/api/lists/:id", hand.userIdentity, hand.deleteList)
	reqGet := httptest.NewRequest(http.MethodDelete,
		"/api/lists/1", nil)
	reqGet.Header["Authorization"] = tok

	w := httptest.NewRecorder()

	r.ServeHTTP(w, reqGet)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	var output todo.TodoList
	if err := json.Unmarshal(data, &output); err != nil {
		log.Fatal(err.Error())
	}
	require.NoError(t, err)
	require.Equal(t, data, []byte("{\"Status\":\"OK\"}"))
}

func TestHandler_UpdateList(t *testing.T) {
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

	out_todoList := todo.UpdateListInput{
		Title:       title,
		Description: description,
	}

	var userId int = 1
	var listId int = 1
	//repoAuth1 := repos.Authorization
	repoList.EXPECT().Update(userId, listId, out_todoList).Return(nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hand := NewHandler(services)

	r.PUT("/api/lists/:id", hand.userIdentity, hand.updateList)
	reqGet := httptest.NewRequest(http.MethodPut,
		"/api/lists/1",
		bytes.NewBuffer(
			[]byte(`{"title": "Купить одежду","description": "Сегодня"}`)))
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
	require.Equal(t, []byte("{\"Status\":\"ok\"}"), data)
}
