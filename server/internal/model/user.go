package model

import (
	"context"
	"time"
)

// maybe change password omitempty to -
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required,min=2,max=50"`
	Password  string    `json:"password,omitempty" db:"password" validate:"required,min=6"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type IUserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByName(ctx context.Context, name string) (*User, error)
	GetByID(ctx context.Context, id int) (*User, error)
}
