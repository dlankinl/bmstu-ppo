package handlers

import (
	"bufio"
	"business-logic/domain"
	"fmt"
	"github.com/google/uuid"
	"os"
	"ppo/app/internal/app"
	"ppo/app/internal/tui/utils"
	"strings"
)

func AddSkill(a *app.App, args ...any) (err error) {
	reader := bufio.NewReader(os.Stdin)

	var name, description string
	fmt.Printf("Введите название навыка: ")
	name, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода названия навыка: %w", err)
	}

	fmt.Printf("Введите описание: ")
	description, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода описания: %w", err)
	}
	description = strings.TrimSpace(description)

	var skill domain.Skill
	skill.Name = name
	skill.Description = description

	err = a.SkillSvc.Create(&skill)
	if err != nil {
		return fmt.Errorf("ошибка добавления навыка: %w", err)
	}

	return nil
}

func DeleteSkill(a *app.App, args ...any) (err error) {
	err = GetSkills(a, args...)
	if err != nil {
		return fmt.Errorf("удаление навыка: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Введите id навыка: ")
	skillId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id навыка: %w", err)
	}
	skillId = strings.TrimSpace(skillId)

	skillUuid, err := uuid.Parse(skillId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	err = a.SkillSvc.DeleteById(skillUuid)
	if err != nil {
		return fmt.Errorf("удаление навыка: %w", err)
	}

	return nil
}

func UpdateSkill(a *app.App, args ...any) (err error) {
	err = GetSkills(a, args...)
	if err != nil {
		return fmt.Errorf("обновление информации о навыке: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Введите id навыка: ")
	skillId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id навыка: %w", err)
	}
	skillId = strings.TrimSpace(skillId)

	skillUuid, err := uuid.Parse(skillId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	skill, err := a.SkillSvc.GetById(skillUuid)
	if err != nil {
		return fmt.Errorf("получение навыка по id: %w", err)
	}

	var name, description string
	fmt.Printf("Введите название навыка (%s): ", skill.Name)
	name, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода названия навыка: %w", err)
	}
	name = strings.TrimSpace(name)
	if name != "" {
		skill.Name = name
	}

	fmt.Printf("Введите описание (%s): ", skill.Description)
	description, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода описания: %w", err)
	}
	description = strings.TrimSpace(description)
	if description != "" {
		skill.Description = description
	}

	err = a.SkillSvc.Update(skill)
	if err != nil {
		return fmt.Errorf("ошибка обновления навыка: %w", err)
	}

	return nil
}

func GetSkills(a *app.App, args ...any) (err error) {
	err = utils.PrintPaginatedCollection("Навыки", a.SkillSvc.GetAll)
	if err != nil {
		return fmt.Errorf("вывод навыков с пагинацией: %w", err)
	}

	return nil
}
