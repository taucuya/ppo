package structs

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id            uuid.UUID
	Name          string
	Date_of_birth time.Time
	Mail          string
	Password      string
	Phone         string
	Address       string
	Status        string
	Role          string
}

var (
	ErrUserNotFound = errors.New("user not found")
)
