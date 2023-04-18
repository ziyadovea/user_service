package entity

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("empty name")
	}
	if u.Email == "" {
		return errors.New("empty email")
	}
	if u.Password == "" {
		return errors.New("empty password")
	}
	return nil
}

func (u *User) HashPassword() error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to generate hash from password: %w", err)
	}
	u.Password = string(hashedPwd)
	return nil
}

func (u *User) ComparePassword(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}
