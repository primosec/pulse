package domain

import (
	"errors"
	"time"
)

var ErrNotFound = errors.New("not found")

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
