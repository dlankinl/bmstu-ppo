package models

import "github.com/google/uuid"

type Contact struct {
	ID    uuid.UUID
	Name  string
	Value string
}
