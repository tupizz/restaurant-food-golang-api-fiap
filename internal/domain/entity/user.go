package entity

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}
