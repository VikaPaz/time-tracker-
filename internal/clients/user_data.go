package client

import (
	"context"
	user_data "github.com/VikaPaz/time_tracker/internal/clients/gen"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

//go:generate oapi-codegen -package=user_data -generate=types,client,spec  -o ./gen/user_data_client.go ./user_data.yaml

type UserInfo struct {
	client *user_data.ClientWithResponses
	server string
	log    *logrus.Logger
}

func NewClient(host string, logger *logrus.Logger) (*UserInfo, error) {
	client, err := user_data.NewClientWithResponses(host)
	if err != nil {
		logger.Error(err)
		return &UserInfo{}, models.ErrClientFailed
	}
	return &UserInfo{
			client: client,
			server: host,
			log:    logger,
		},
		nil
}

func (u *UserInfo) GetInf(ctx context.Context, p models.CreateUserRequest) (models.User, error) {
	u.log.Debugf("Validating passport %v", p.PassportNumber)
	passport := *p.PassportNumber
	if passport == "" {
		return models.User{}, models.ErrInvalidPassword
	}
	data := strings.Split(passport, " ")
	if len(data) != 2 || len(data[0]) != 4 || len(data[1]) != 6 {
		return models.User{}, models.ErrInvalidPassword
	}

	series, err := strconv.Atoi(data[0])
	if err != nil {
		return models.User{}, models.ErrInvalidPassword
	}
	number, err := strconv.Atoi(data[1])
	if err != nil {
		return models.User{}, models.ErrInvalidPassword
	}
	params := user_data.GetInfoParams{
		PassportSerie:  series,
		PassportNumber: number,
	}

	u.log.Debugf("Getting user information  with %v", params)
	resp, err := u.client.GetInfoWithResponse(ctx, &params)
	if err != nil {
		return models.User{}, err
	}
	userInfo := models.User{
		Name:       &resp.JSON200.Name,
		Surname:    &resp.JSON200.Surname,
		Patronymic: resp.JSON200.Patronymic,
		Address:    &resp.JSON200.Address,
		Passport:   &passport,
	}

	return userInfo, nil
}
