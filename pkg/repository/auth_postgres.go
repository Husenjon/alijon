package repository

import (
	"errors"
	"fmt"

	inkassback "github.com/Husenjon/InkassBack"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NweAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user inkassback.User) (inkassback.User, error) {
	users := []inkassback.User{}
	query := fmt.Sprintf(
		`insert into
			%s (
				ism,
				familya,
				otasini_ismi,
				phone,
				username,
				password,
				branch_id,
				image,
				is_active
			)
		values
			(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9
			) returning *;`,
		usersTable,
	)
	err := r.db.Select(
		&users,
		query,
		user.Ism,
		user.Familya,
		user.OtasiniIsmi,
		user.Phone,
		user.Username,
		user.Password,
		user.BranchId,
		user.Image,
		user.IsActive,
	)
	if err != nil {
		return users[0], err
	}
	return users[0], nil
}
func (r *AuthPostgres) GetUser(username, password string) (inkassback.User, error) {
	users := []inkassback.User{}
	query := fmt.Sprintf(
		`SELECT
			*
		FROM
			%s
		WHERE
			username = $1
			AND password = $2
			AND is_active = true;`,
		usersTable)
	err := r.db.Select(&users, query, username, password)
	if err != nil {
		return users[0], err
	}
	if len(users) == 0 {
		u := inkassback.User{}
		return u, errors.New("user not found")
	}
	return users[0], nil
}
func (r *AuthPostgres) UpdateToken(id int, token string) (inkassback.User, error) {
	users := []inkassback.User{}
	query := fmt.Sprintf(
		`update
			%s
		set
			token = $1
		where
			id = $2
		returning *;`,
		usersTable)
	err := r.db.Select(&users, query, token, id)
	if err != nil {
		return users[0], err
	}
	if len(users) == 0 {
		u := inkassback.User{}
		return u, errors.New("user not found")
	}
	return users[0], nil
}
