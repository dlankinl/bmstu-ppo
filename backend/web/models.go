package web

import (
	"ppo/domain"
	"time"

	"github.com/google/uuid"
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

type UserSkill struct {
	UserId  uuid.UUID `json:"user_id,omitempty"`
	SkillId uuid.UUID `json:"skill_id,omitempty"`
}

type FinancialReport struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CompanyID uuid.UUID `json:"company_id,omitempty"`
	Revenue   float32   `json:"revenue,omitempty"`
	Costs     float32   `json:"costs,omitempty"`
	Year      int       `json:"year,omitempty"`
	Quarter   int       `json:"quarter,omitempty"`
}

type Period struct {
	StartYear    int `json:"start_year"`
	StartQuarter int `json:"start_quarter"`
	EndYear      int `json:"end_year"`
	EndQuarter   int `json:"end_quarter"`
}

type Review struct {
	ID          uuid.UUID `json:"id"`
	Target      uuid.UUID `json:"target_id"`
	Reviewer    uuid.UUID `json:"reviewer_id"`
	Pros        string    `json:"pros"`
	Cons        string    `json:"cons"`
	Description string    `json:"description"`
	Rating      int       `json:"rating"`
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

func toUserSkillTransport(userSkill *domain.UserSkill) UserSkill {
	return UserSkill{
		UserId:  userSkill.UserId,
		SkillId: userSkill.SkillId,
	}
}

func toUserSkillModel(userSkill *UserSkill) domain.UserSkill {
	return domain.UserSkill{
		UserId:  userSkill.UserId,
		SkillId: userSkill.SkillId,
	}
}

func toFinReportTransport(finReport *domain.FinancialReport) FinancialReport {
	return FinancialReport{
		ID:        finReport.ID,
		CompanyID: finReport.CompanyID,
		Revenue:   finReport.Revenue,
		Costs:     finReport.Costs,
		Year:      finReport.Year,
		Quarter:   finReport.Quarter,
	}
}

func toFinReportModel(finReport *FinancialReport) domain.FinancialReport {
	return domain.FinancialReport{
		ID:        finReport.ID,
		CompanyID: finReport.CompanyID,
		Revenue:   finReport.Revenue,
		Costs:     finReport.Costs,
		Year:      finReport.Year,
		Quarter:   finReport.Quarter,
	}
}

func toPeriodTransport(per *domain.Period) Period {
	return Period{
		StartYear:    per.StartYear,
		StartQuarter: per.StartQuarter,
		EndYear:      per.EndYear,
		EndQuarter:   per.EndQuarter,
	}
}

func toReviewTransport(rev *domain.Review) Review {
	return Review{
		ID:          rev.ID,
		Target:      rev.Target,
		Reviewer:    rev.Reviewer,
		Pros:        rev.Pros,
		Cons:        rev.Cons,
		Description: rev.Description,
		Rating:      rev.Rating,
	}
}

func toReviewModel(rev *Review) domain.Review {
	return domain.Review{
		ID:          rev.ID,
		Target:      rev.Target,
		Reviewer:    rev.Reviewer,
		Pros:        rev.Pros,
		Cons:        rev.Cons,
		Description: rev.Description,
		Rating:      rev.Rating,
	}
}
