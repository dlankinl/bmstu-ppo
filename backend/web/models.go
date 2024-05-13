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

type Contact struct {
	ID      uuid.UUID `json:"id,omitempty"`
	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Value   string    `json:"value,omitempty"`
}

type ActivityField struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Cost        float32   `json:"cost,omitempty"`
}

type Company struct {
	ID              uuid.UUID `json:"id,omitempty"`
	OwnerID         uuid.UUID `json:"owner_id,omitempty"`
	ActivityFieldId uuid.UUID `json:"activity_field_id,omitempty"`
	Name            string    `json:"name,omitempty"`
	City            string    `json:"city,omitempty"`
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

func toContactTransport(contact *domain.Contact) Contact {
	return Contact{
		ID:      contact.ID,
		OwnerID: contact.OwnerID,
		Name:    contact.Name,
		Value:   contact.Value,
	}
}

func toContactModel(contact *Contact) domain.Contact {
	return domain.Contact{
		ID:      contact.ID,
		OwnerID: contact.OwnerID,
		Name:    contact.Name,
		Value:   contact.Value,
	}
}

func toActFieldTransport(field *domain.ActivityField) ActivityField {
	return ActivityField{
		ID:          field.ID,
		Name:        field.Name,
		Description: field.Description,
		Cost:        field.Cost,
	}
}

func toActFieldModel(field *ActivityField) domain.ActivityField {
	return domain.ActivityField{
		ID:          field.ID,
		Name:        field.Name,
		Description: field.Description,
		Cost:        field.Cost,
	}
}

func toCompanyTransport(company *domain.Company) Company {
	return Company{
		ID:              company.ID,
		OwnerID:         company.OwnerID,
		ActivityFieldId: company.ActivityFieldId,
		Name:            company.Name,
		City:            company.City,
	}
}

func toCompanyModel(company *Company) domain.Company {
	return domain.Company{
		ID:              company.ID,
		OwnerID:         company.OwnerID,
		ActivityFieldId: company.ActivityFieldId,
		Name:            company.Name,
		City:            company.City,
	}
}
