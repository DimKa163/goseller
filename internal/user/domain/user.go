package domain

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type UserID int64

func NewUserID(s string) (UserID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("invalid user ID: %w", err)
	}
	return UserID(id), nil
}

func (id UserID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

type User struct {
	ID        UserID    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     Email     `json:"email"`
	Phone     Phone     `json:"phone"`
	IsActive  bool      `json:"is_active"`
}

func (u *User) Update(name string, email Email, phone Phone) {
	u.Name = name
	u.Email = email
	u.Phone = phone
	u.UpdatedAt = time.Now()
}

func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

func CreateNewUser(name string, email Email, phone Phone) *User {
	return &User{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		IsActive:  true,
	}
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (UserID, error)
	GetByID(ctx context.Context, id UserID) (*User, error)
	GetByEmail(ctx context.Context, email Email) (*User, error)
	GetCountByEmail(ctx context.Context, email Email) (int64, error)
	GetByPhone(ctx context.Context, phone Phone) (*User, error)
	GetCountByPhone(ctx context.Context, phone Phone) (int64, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id UserID) error
}
