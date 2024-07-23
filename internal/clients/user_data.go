package client

import (
	"context"
	"fmt"
	user_data "github.com/VikaPaz/time_tracker/internal/clients/gen"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/VikaPaz/time_tracker/internal/server/user"
	"strconv"
	"strings"
)

//go:generate oapi-codegen -package=user_data -generate=types,client,spec  -o ./gen/user_data_client.go ./user_data.yaml

type UserInfo struct {
	client *user_data.ClientWithResponses
	server string
}

func NewClient(host string) *UserInfo {
	client, _ := user_data.NewClientWithResponses(host) //"http://127.0.0.1:8080"
	return &UserInfo{
		client: client,
		server: host,
	}
}

func (u *UserInfo) GetInf(ctx context.Context, p user.Person) (models.User, error) {
	passport := *p.PassportNumber
	if passport == "" {
		return models.User{}, fmt.Errorf("passport is nil")
	}
	data := strings.Split(passport, " ")
	if len(data) != 2 || len(data[0]) != 4 || len(data[1]) != 6 {
		return models.User{}, fmt.Errorf("passport is invalid")
	}

	series, err := strconv.Atoi(data[0])
	if err != nil {
		return models.User{}, err
	}
	number, err := strconv.Atoi(data[1])
	if err != nil {
		return models.User{}, err
	}
	params := user_data.GetInfoParams{
		PassportSerie:  series,
		PassportNumber: number,
	}

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
