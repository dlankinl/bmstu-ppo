package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"ppo/domain"
)

func UserDbToUser(in *User) *domain.User {
	return &domain.User{
		ID:       in.ID,
		Username: in.Username.String,
		FullName: in.FullName.String,
		Birthday: in.Birthday.Time,
		Gender:   in.Gender.String,
		City:     in.City.String,
		Role:     in.Role.String,
	}
}

type User struct {
	ID       uuid.UUID
	Username sql.NullString
	FullName sql.NullString
	Gender   sql.NullString
	Birthday sql.NullTime
	City     sql.NullString
	Role     sql.NullString
}
