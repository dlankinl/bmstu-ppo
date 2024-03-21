package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       uuid.UUID
	Username string
	FullName string
	Gender   int
	Birthday time.Time
	City     string
}
