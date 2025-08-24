package ad

import "errors"

type Ad struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	UserID      int    `json:"user_id"`
}

var ErrForbidden = errors.New("forbidden error")
