package dto

import (
	"bulletin-board/internal/user"
	"time"
)

type ResponseUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Birthday string `json:"birthday"`
	Contact  string `json:"contact"`
}

type RequestUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Birthday string `json:"birthday"`
	Contact  string `json:"contact"`
}

func ToDto(user user.User) ResponseUser {
	return ResponseUser{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Birthday: user.Birthday.Format("2006-01-02"),
		Contact:  user.Contact,
	}
}

func ToUser(requestUser RequestUser) user.User {
	birthday, err := time.Parse("2006-01-02", requestUser.Birthday)
	if err != nil {
		return user.User{}
	}
	return user.User{
		Name:     requestUser.Name,
		Email:    requestUser.Email,
		Password: requestUser.Password,
		Birthday: birthday,
		Contact:  requestUser.Contact,
	}
}
