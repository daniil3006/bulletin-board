package user

import (
	"errors"
	"time"
)

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Birthday time.Time `json:"birthday"`
	Contact  string    `json:"contact"`
}

var ErrInvalidUserId = errors.New("invalid id")
