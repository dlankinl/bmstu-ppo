package models

import "github.com/google/uuid"

type Skill struct {
	ID          uuid.UUID
	Name        string
	Description string
}
