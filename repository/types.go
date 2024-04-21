// This file contains types that are used in the repository layer.
package repository

import (
	"database/sql"
	"time"
)

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type RegisterUser struct {
	ID       string `db:"id"`
	Phone    string `db:"phone"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

type User struct {
	UserID    string       `json:"userID" db:"id"`
	Phone     string       `json:"phone" db:"phone"`
	Name      string       `json:"name" db:"name"`
	Password  string       `json:"password" db:"password"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updatedAt" db:"updated_at"`
}

type UpdateUser struct {
	Phone string `db:"phone"`
	Name  string `db:"name"`
}
