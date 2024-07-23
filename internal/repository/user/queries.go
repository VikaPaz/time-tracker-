package user

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/google/uuid"
)

type UserRepository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn: conn}
}

func (r *UserRepository) Create(user models.User) (models.User, error) {
	row := r.conn.QueryRow("INSERT INTO users (passport, name, surname, patronymic, address) values "+
		"($1, $2, $3, $4, $5) RETURNING id", user.Passport, user.Name, user.Surname, user.Patronymic, user.Address)
	if err := row.Err(); err != nil {
		return models.User{}, err
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return models.User{}, err
	}

	user.ID = &id
	return user, nil
}

func (r *UserRepository) Get(f models.FilterRequest) (models.FilterResponse, error) {
	var users []models.User

	builder := sq.Select("count(*) over ()", "id", "passport", "name", "surname", "patronymic", "address").From("users")
	builder = builder.PlaceholderFormat(sq.Dollar)
	if f.Fields.ID != nil {
		builder = builder.Where(sq.Eq{"ILike": f.Fields.ID})
	}
	if f.Fields.Name != nil {
		//builder = builder.Where(sq.Eq{"name": *f.Fields.Name})
		builder = builder.Where(sq.ILike{"name": fmt.Sprintf("%%%v%%", *f.Fields.Name)})
	}
	if f.Fields.Surname != nil {
		builder = builder.Where(sq.ILike{"surname": fmt.Sprintf("%%%v%%", *f.Fields.Surname)})
	}
	if f.Fields.Patronymic != nil {
		builder = builder.Where(sq.ILike{"patronymic": fmt.Sprintf("%%%v%%", *f.Fields.Patronymic)})
	}
	if f.Fields.Address != nil {
		builder = builder.Where(sq.ILike{"address": fmt.Sprintf("%%%v%%", *f.Fields.Address)})
	}
	if f.Fields.Passport != nil {
		builder = builder.Where(sq.ILike{"passport": fmt.Sprintf("%%%v%%", *f.Fields.Passport)})
	}

	if f.Limit != 0 {
		builder = builder.Limit(f.Limit)
	}
	if f.Offset != 0 {
		builder = builder.Offset(f.Offset)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return models.FilterResponse{}, err
	}

	rows, err := r.conn.Query(query, args...)
	if err != nil {
		return models.FilterResponse{}, err
	}

	result := models.FilterResponse{}
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&result.Total, &user.ID, &user.Passport, &user.Name, &user.Surname, &user.Patronymic, &user.Address)
		if err != nil {
			if err == sql.ErrNoRows {
				return models.FilterResponse{}, nil
			}
			return models.FilterResponse{}, err
		}
		users = append(users, user)
	}
	result.Users = users
	return result, nil
}

func (r *UserRepository) Delete(request models.DeleteUserRequest) error {
	builder := sq.Delete("users").Where(sq.Eq{"id": request.ID})
	builder = builder.PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Set(user models.User) error {
	builder := sq.Update("users").Where(sq.Eq{"id": user.ID})
	builder = builder.PlaceholderFormat(sq.Dollar)
	if user.Name != nil {
		builder = builder.Set("name", user.Name)
	}
	if user.Surname != nil {
		builder = builder.Set("surname", user.Surname)
	}
	if user.Patronymic != nil {
		builder = builder.Set("patronymic", user.Patronymic)
	}
	if user.Address != nil {
		builder = builder.Set("address", user.Address)
	}
	if user.Passport != nil {
		builder = builder.Set("passport", user.Passport)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.conn.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
