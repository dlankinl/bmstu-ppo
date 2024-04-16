package handlers

import (
	"fmt"
	"ppo/internal/app"
	"ppo/internal/config"
	"ppo/internal/tui/utils"
)

func GetAllUsers(a *app.App) (err error) {
	page := 1
	flag := true
	for flag {
		users, err := a.UserSvc.GetAll(page)
		if err != nil {
			return fmt.Errorf("получение пользователей: %w, err")
		}

		if len(users) < config.PageSize {
			flag = false
		} else {
			fmt.Println("Выберите действие:\n1. Следующая страница.\n2. Предыдущая страница.\n0. Назад.")
			var option int
			_, err = fmt.Scanf("%d", &option)
			if err != nil {
				return fmt.Errorf("ошибка ввода следующего действия: %w", err)
			}

			switch option {
			case 1:
				page++
			case 2:
				if page > 1 {
					page--
				}
			case 0:
				break
			}
		}

		utils.PrintUsers(users)
	}

	return nil
}
