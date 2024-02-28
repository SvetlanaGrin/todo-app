package service

import (
	"errors"
	"github.com/SvetlanaGrin/todo-app"
	repoMock "github.com/SvetlanaGrin/todo-app/pkg/mocks"
	"github.com/SvetlanaGrin/todo-app/pkg/repository"
	"github.com/golang/mock/gomock"
	require "github.com/stretchr/testify/require"
	"testing"
)

func Test_GeneratePasswordHash(t *testing.T) {
	password := "12345"
	output := generatePasswordHash(password)
	require.Equal(t, output, "33347266766767636867766365323133646664668cb2237d0679ca88db6464eac60da96345513964")
}

func TestAuthService_GenerateToken(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos := repoMock.NewMockAuthorization(ctl)
	services := NewAuthService(repos)
	var username string = "Alex"
	var password string = generatePasswordHash("12345")
	User1 := todo.User{
		Id:       1,
		Name:     "Sweta",
		Username: "Alex",
		Password: generatePasswordHash("12345"),
	}

	repos.EXPECT().GetUser(username, password).Return(User1, nil).Times(1)
	output, err := services.GenerateToken(username, "12345")
	require.NoError(t, err)
	require.NotEqual(t, "", output)

}

func TestAuthService_GenerateTokenError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos := repoMock.NewMockAuthorization(ctl)
	services := NewAuthService(repos)
	var username string = "Alex"
	var password string = generatePasswordHash("")

	err := errors.New("Error")
	repos.EXPECT().GetUser(username, password).Return(todo.User{}, err).Times(1)
	output, err := services.GenerateToken(username, "")
	require.EqualError(t, err, "Error")
	require.Equal(t, "", output)
}

/*func TestAuthService_ParseToken(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos := repoMock.NewMockAuthorization(ctl)
	services := NewAuthService(repos)
	var token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDY2ODgwNDAsImlhdCI6MTcwNjY0NDg0MCwiVXNlcklkIjoxfQ.A_WpkXOm4do3fptK9EgPEkZFLgNTY8Bk4e2U1GMHJFA"

	output, err := services.ParseToken(token)
	require.NoError(t, err)
	require.NotEqual(t, output, 0)

}*/

func TestAuthService_ParseTokenInvalidNumberError(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos := repoMock.NewMockAuthorization(ctl)
	services := NewAuthService(repos)
	var token string = "eyJhbGciOiJI223456"

	output, err := services.ParseToken(token)
	require.EqualError(t, err, "token contains an invalid number of segments")
	require.Equal(t, 0, output)
}

func TestAuthService_CreateUsers(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repos := repoMock.NewMockAuthorization(ctl)
	services := NewAuthService(repos)
	UserServices := todo.User{
		Id:       1,
		Name:     "Sweta",
		Username: "Alexndrovna",
		Password: "12345",
	}

	UserRepos := todo.User{
		Id:       1,
		Name:     "Sweta",
		Username: "Alexndrovna",
		Password: generatePasswordHash("12345"),
	}
	repos.EXPECT().CreateUsers(UserRepos).Return(UserRepos.Id, nil).Times(1)
	output, err := services.CreateUsers(UserServices)
	require.NoError(t, err)
	require.Equal(t, output, UserRepos.Id)
}

type Mocker struct {
	Authorization
	TodoList
	TodoItem
}

func TestNewService(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	repo := &repository.Repository{}
	services := NewService(repo)
	require.Equal(t, &Service{
		Authorization: &AuthService{},
		TodoList:      &TodoListService{},
		TodoItem:      &TodoItemService{},
	}, services)
}
