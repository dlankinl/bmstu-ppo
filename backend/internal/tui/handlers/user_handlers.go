package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"ppo/internal/app"
	"ppo/internal/config"
	"ppo/internal/tui/utils"
	"strings"

	"github.com/google/uuid"
)

func GetAllUsers(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	page := 1
	for {
		users, err := a.UserSvc.GetAll(ctx, page)
		if err != nil {
			return fmt.Errorf("получение пользователей: %w", err)
		}

		utils.PrintCollection("Предприниматели", users)

		fmt.Printf("1. Предыдущая страница.\n2. Следующая страница.\n0. Назад.\n\nВыберите действие: ")
		var option int
		_, err = fmt.Scanf("%d", &option)
		if err != nil {
			return fmt.Errorf("ошибка ввода следующего действия: %w", err)
		}

		switch option {
		case 1:
			if page > 1 {
				page--
			}
		case 2:
			if len(users) == config.PageSize {
				page++
			}
		case 0:
			return nil
		}
	}
}

func CalculateRating(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)
	var idStr string

	err = GetAllUsers(a, args...)
	if err != nil {
		return fmt.Errorf("вывод пользователей: %w", err)
	}

	fmt.Printf("Введите id: ")
	idStr, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id: %w", err)
	}
	idStr = strings.TrimSpace(idStr)

	id, err := uuid.Parse(idStr)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	rating, err := a.Interactor.CalculateUserRating(ctx, id)
	if err != nil {
		return fmt.Errorf("расчёт рейтинга: %w", err)
	}

	fmt.Printf("Рейтинг пользователя с id=%s равен %f\n", id, rating)

	return nil
}
