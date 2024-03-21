package services

import (
	"github.com/google/uuid"
	"ppo/domain"
)

type SkillService struct {
	skillRepo domain.ISkillRepository
}

func (s SkillService) Create(skill *domain.Skill) (err error) {
	return nil
}

func (s SkillService) GetById(id uuid.UUID) (skill *domain.Skill, err error) {
	return skill, nil
}

func (s SkillService) GetAll() (skills []*domain.Skill, err error) {
	return skills, nil
}

func (s SkillService) Update(skill *domain.Skill) (err error) {
	return nil
}

func (s SkillService) DeleteById(id uuid.UUID) (err error) {
	return nil
}
