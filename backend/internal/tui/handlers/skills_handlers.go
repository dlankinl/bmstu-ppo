package handlers

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"os"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/tui/utils"
	"strings"
)

func AddSkill(a *app.App, args ...any) (err error) {
	ctx := context.Background()
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

	err = a.SkillSvc.Create(ctx, &skill)
	if err != nil {
		return fmt.Errorf("ошибка добавления навыка: %w", err)
	}

	return nil
}

func DeleteSkill(a *app.App, args ...any) (err error) {
	ctx := context.Background()
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

	err = a.SkillSvc.DeleteById(ctx, skillUuid)
	if err != nil {
		return fmt.Errorf("удаление навыка: %w", err)
	}

	return nil
}

func UpdateSkill(a *app.App, args ...any) (err error) {
	ctx := context.Background()
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

	skill, err := a.SkillSvc.GetById(ctx, skillUuid)
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

	err = a.SkillSvc.Update(ctx, skill)
	if err != nil {
		return fmt.Errorf("ошибка обновления навыка: %w", err)
	}

	return nil
}

func GetSkills(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	err = utils.PrintPaginatedCollection("Навыки", a.SkillSvc.GetAll, ctx)
	if err != nil {
		return fmt.Errorf("вывод навыков с пагинацией: %w", err)
	}

	return nil
}
