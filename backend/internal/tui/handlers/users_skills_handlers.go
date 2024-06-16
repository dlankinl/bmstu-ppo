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

func AddUserSkill(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	var username string
	var ok bool
	if len(args) > 0 {
		username, ok = args[0].(string)
		if !ok {
			return fmt.Errorf("приведение аргумента к string")
		}
	}

	user, err := a.UserSvc.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	err = GetSkills(a, args...)
	if err != nil {
		return fmt.Errorf("добавление навыка пользователю: %w", err)
	}

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

	var userSkill domain.UserSkill
	userSkill.UserId = user.ID
	userSkill.SkillId = skillUuid

	err = a.UserSkillSvc.Create(ctx, &userSkill)
	if err != nil {
		return fmt.Errorf("ошибка добавления навыка пользователю: %w", err)
	}

	return nil
}

func DeleteUserSkill(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	var username string
	var ok bool
	if len(args) > 0 {
		username, ok = args[0].(string)
		if !ok {
			return fmt.Errorf("приведение аргумента к string")
		}
	}

	user, err := a.UserSvc.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	err = GetMySkills(a, args...)
	if err != nil {
		return fmt.Errorf("удаление навыка пользователя: %w", err)
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

	pair := &domain.UserSkill{
		UserId:  user.ID,
		SkillId: skillUuid,
	}

	err = a.UserSkillSvc.Delete(ctx, pair)
	if err != nil {
		return fmt.Errorf("удаление навыка пользователя: %w", err)
	}

	return nil
}

func GetMySkills(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	var username string
	var ok bool
	if len(args) > 0 {
		username, ok = args[0].(string)
		if !ok {
			return fmt.Errorf("приведение аргумента к string")
		}
	}

	user, err := a.UserSvc.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("пользователь не найден")
	}

	err = utils.PrintPaginatedCollectionArgs("Навыки", a.UserSkillSvc.GetSkillsForUser, ctx, user.ID, true)
	if err != nil {
		return fmt.Errorf("вывод навыков с пагинацией: %w", err)
	}

	return nil
}
