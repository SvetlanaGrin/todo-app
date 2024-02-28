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

type GetToken struct {
	Token string `json:"token"`
}

var GlobalToken string

func TestAuthHandler_signUp(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)
	UserRepos := todo.User{
		Id:       0,
		Name:     "Sveta",
		Username: "Alexndrovna",
		Password: "33347266766767636867766365323133646664668cb2237d0679ca88db6464eac60da96345513964",
	}
	var newId int = 1
	//repoAuth1 := repos.Authorization
	repoAuth.EXPECT().CreateUsers(UserRepos).Return(newId, nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hend := NewHandler(services)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/auth/sign-up", hend.signUp)
	req := httptest.NewRequest(http.MethodPost,
		"/auth/sign-up",
		bytes.NewBuffer(
			[]byte(`{"name": "Sveta","username": "Alexndrovna","password": "12345"}`)))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	expected := "{\"id\":1}"
	require.Equal(t, expected, string(data))

}

func TestAuthHandler_signUpJsonError(t *testing.T) {
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
	hend := NewHandler(services)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/auth/sign-up", hend.signUp)
	req := httptest.NewRequest(http.MethodPost,
		"/auth/sign-up",
		bytes.NewBuffer(
			[]byte(`{"name": "Sveta","username": "Alexndrovna","password "12345"}`)))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	require.Equal(t, w.Code, http.StatusBadRequest)
	require.Equal(t, data, []uint8([]byte{0x7b, 0x22, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
		0x22, 0x3a, 0x22, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x20,
		0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x20, 0x27,
		0x31, 0x27, 0x20, 0x61, 0x66, 0x74, 0x65, 0x72, 0x20, 0x6f, 0x62,
		0x6a, 0x65, 0x63, 0x74, 0x20, 0x6b, 0x65, 0x79, 0x22, 0x7d}))

}

func TestAuthHandler_signIn(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repoAuth := repoMock.NewMockAuthorization(ctl)
	repoList := repoMock.NewMockTodoList(ctl)
	repoItem := repoMock.NewMockTodoItem(ctl)

	UserRepos := todo.User{
		Id:       1,
		Name:     "Sveta",
		Username: "Alexndrovna",
		Password: "33347266766767636867766365323133646664668cb2237d0679ca88db6464eac60da96345513964",
	}
	//repoAuth1 := repos.Authorization
	repoAuth.EXPECT().GetUser(UserRepos.Username, UserRepos.Password).Return(UserRepos, nil)

	services := &service.Service{
		Authorization: service.NewAuthService(repoAuth),
		TodoList:      service.NewTodoListService(repoList),
		TodoItem:      service.NewTodoItemService(repoItem, repoList),
	}
	hend := NewHandler(services)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/auth/sign-in", hend.signIn)
	req := httptest.NewRequest(http.MethodPost,
		"/auth/sign-in",
		bytes.NewBuffer(
			[]byte(`{"username": "Alexndrovna","password": "12345"}`)))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	res := w.Result()

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	var getToken GetToken
	if err := json.Unmarshal(data, &getToken); err != nil {
		log.Fatal(err.Error())
	}
	GlobalToken = getToken.Token

}
func TestAuthHandler_signInJconError(t *testing.T) {
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
	hend := NewHandler(services)
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.POST("/auth/sign-in", hend.signIn)
	req := httptest.NewRequest(http.MethodPost,
		"/auth/sign-in",
		bytes.NewBuffer(
			[]byte(`{"username": "Alexndrovna","password0"12345"}`)))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	res := w.Result()

	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)

	require.Equal(t, w.Code, http.StatusBadRequest)
	require.Equal(t, data, []uint8([]byte{0x7b, 0x22, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
		0x22, 0x3a, 0x22, 0x69, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x20,
		0x63, 0x68, 0x61, 0x72, 0x61, 0x63, 0x74, 0x65, 0x72, 0x20, 0x27,
		0x31, 0x27, 0x20, 0x61, 0x66, 0x74, 0x65, 0x72, 0x20, 0x6f, 0x62,
		0x6a, 0x65, 0x63, 0x74, 0x20, 0x6b, 0x65, 0x79, 0x22, 0x7d}))

}
