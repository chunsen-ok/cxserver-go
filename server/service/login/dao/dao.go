package dao

import (
	"context"
	"cxfw/db"
	"cxfw/model"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(u *model.User) (int, *model.User, error) {
	m := model.User{}
	err := db.S().QueryRow(context.Background(), `select * from users where account = $1`, u.Account).Scan(&m.ID, &m.Account, &m.Name, &m.Password)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return http.StatusInternalServerError, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(u.Password)); err != nil {
		return http.StatusInternalServerError, nil, err
	}

	m.Password = ""

	return http.StatusOK, &m, nil
}

func Logout(u *model.User) bool {
	return true
}
