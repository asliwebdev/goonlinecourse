package repository

import (
	"database/sql"
	"lesson29/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) GetUserById(userId string) (*models.User, error) {
	user := &models.User{}
	row := u.db.QueryRow("SELECT id, name, age FROM users WHERE id = $1", userId)

	if err := row.Scan(&user.Id, &user.Name, &user.Age); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
