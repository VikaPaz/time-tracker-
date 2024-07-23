package user

import (
	"context"
	"fmt"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/VikaPaz/time_tracker/internal/server/user"
)

type UserService struct {
	repo     Repository
	userData Client
}

type Client interface {
	GetInf(ctx context.Context, p user.Person) (models.User, error)
}

type Repository interface {
	Create(user models.User) (models.User, error)
	Get(models.FilterRequest) (models.FilterResponse, error)
	Delete(request models.DeleteUserRequest) error
	Set(request models.User) error
}

type PeopleInfo interface {
	GetInfo(string2 string)
}

func NewService(repo Repository, userData Client) *UserService {
	return &UserService{repo: repo, userData: userData}
}

func (u *UserService) CreateUser(person user.Person, ctx context.Context) (models.User, error) { // TODO return full user model
	filter := models.FilterRequest{Fields: models.User{Passport: person.PassportNumber}, Limit: 1}
	result, err := u.repo.Get(filter)
	if err != nil {
		return models.User{}, err
	}
	if result.Users != nil {
		return result.Users[0], fmt.Errorf("user already exists")
	}

	info, err := u.userData.GetInf(ctx, person)
	if err != nil {
		return models.User{}, err
	}

	userInf, err := u.repo.Create(info)
	if err != nil {
		return models.User{}, err
	}

	return userInf, nil
}

func (u *UserService) DeleteUser(request models.DeleteUserRequest) error {
	err := u.repo.Delete(request)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) ChangeUser(request models.User) error {
	err := u.repo.Set(request)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetUsers(filter models.FilterRequest) (models.FilterResponse, error) {
	result, err := u.repo.Get(filter)
	if err != nil {
		return models.FilterResponse{}, err
	}
	return result, nil
}
