package web

import (
	"github.com/google/uuid"
	"ppo/domain"
	"time"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`
	FullName string    `json:"full_name,omitempty"`
	Gender   string    `json:"gender,omitempty"`
	Birthday time.Time `json:"birthday,omitempty"`
	City     string    `json:"city,omitempty"`
	Role     string    `json:"role,omitempty"`
}

type Skill struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
}

func toUserTransport(user *domain.User) User {
	return User{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Gender:   user.Gender,
		Birthday: user.Birthday,
		City:     user.City,
		Role:     user.Role,
	}
}

func toUserModel(user *User) domain.User {
	return domain.User{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Gender:   user.Gender,
		Birthday: user.Birthday,
		City:     user.City,
		Role:     user.Role,
	}
}

func toSkillTransport(skill *domain.Skill) Skill {
	return Skill{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
	}
}

func toSkillModel(skill *Skill) domain.Skill {
	return domain.Skill{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
	}
}
