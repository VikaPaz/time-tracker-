package user

import (
	"context"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	repo     Repository
	userData Client
	log      *logrus.Logger
}

type Client interface {
	GetInf(ctx context.Context, p models.CreateUserRequest) (models.User, error)
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

func NewService(repo Repository, userData Client, logger *logrus.Logger) *UserService {
	return &UserService{
		repo:     repo,
		userData: userData,
		log:      logger,
	}
}

func (u *UserService) CreateUser(person models.CreateUserRequest, ctx context.Context) (models.User, error) {
	u.log.Debugf("Checking user exists")
	filter := models.FilterRequest{Fields: models.User{Passport: person.PassportNumber}, Limit: 1}
	result, err := u.repo.Get(filter)
	if err != nil {
		return models.User{}, err
	}
	if result.Users != nil {
		return result.Users[0], models.ErrUserExists
	}

	u.log.Infof("Getting user information")
	info, err := u.userData.GetInf(ctx, person)
	if err != nil {
		return models.User{}, err
	}

	u.log.Debugf("Creating user: %v", info)
	userInf, err := u.repo.Create(info)
	if err != nil {
		return models.User{}, err
	}

	return userInf, nil
}

func (u *UserService) DeleteUser(request models.DeleteUserRequest) error {
	u.log.Debugf("Deleting user with ID: %v", request.ID)
	err := u.repo.Delete(request)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) ChangeUser(request models.User) error {
	u.log.Debugf("Changing user information: %v", request)
	err := u.repo.Set(request)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetUsers(filter models.FilterRequest) (models.FilterResponse, error) {
	u.log.Debugf("Getting users with filter: %v", filter)
	result, err := u.repo.Get(filter)
	if err != nil {
		return models.FilterResponse{}, err
	}
	return result, nil
}
