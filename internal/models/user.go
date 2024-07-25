package models

import (
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	PassportNumber *string `json:"passportNumber,omitempty"`
}

type User struct {
	ID         *uuid.UUID `json:"id,omitempty"`
	Passport   *string    `json:"passport,omitempty"`
	Name       *string    `json:"name,omitempty"`
	Surname    *string    `json:"surname,omitempty"`
	Patronymic *string    `json:"patronymic,omitempty"`
	Address    *string    `json:"address,omitempty"`
}

type FilterRequest struct {
	Fields User   `json:"fields,omitempty"`
	Limit  uint64 `json:"limit,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
}

type FilterResponse struct {
	Users []User
	Total int64
}

type DeleteUserRequest struct {
	ID uuid.UUID `json:"id,omitempty"`
}
