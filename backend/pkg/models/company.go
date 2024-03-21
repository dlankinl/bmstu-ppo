package models

import "github.com/google/uuid"

type Company struct {
	ID              uuid.UUID
	Name            string
	FieldOfActivity string
	City            string
}
