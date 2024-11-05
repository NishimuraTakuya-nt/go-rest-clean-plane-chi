package models

import "time"

type Sample struct {
	ID        string    `json:"id"`
	StringVal string    `json:"string_val"`
	IntVal    int       `json:"int_val"`
	ArrayVal  []string  `json:"array_val"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
