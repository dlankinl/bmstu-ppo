package utils

import (
	"fmt"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/config"
)

func PrintUser(user *domain.User) {
	var gender string
	if user.Gender == "m" {
		gender = "Мужской"
	} else if user.Gender == "w" {
		gender = "Женский"
	}

	fmt.Printf("%s | %s | %s | %s | %s | %s\n", user.ID.String(), user.Username, user.FullName, gender, user.Birthday.Format("2006-01-02"), user.City)
}

func PrintUsers(users []*domain.User) {
	fmt.Println("Предприниматели:")
	for _, user := range users {
		PrintUser(user)
	}
	fmt.Println()
}

func PrintActivityField(field *domain.ActivityField) {
	fmt.Printf("%s | %s | %s | %f", field.ID, field.Name, field.Description, field.Cost)
}

func PrintActivityFields(fields []*domain.ActivityField) {
	fmt.Println("Сферы деятельности:")
	for _, field := range fields {
		PrintActivityField(field)
	}
}

func PrintPaginatedCollection(a *app.App) (err error) {
	page := 1
	for {
		users, err := a.UserSvc.GetAll(page)
		if err != nil {
			return fmt.Errorf("получение пользователей: %w", err)
		}

		PrintUsers(users)

		fmt.Printf("1. Следующая страница.\n2. Предыдущая страница.\n0. Назад.\n\nВыберите действие: ")
		var option int
		_, err = fmt.Scanf("%d", &option)
		if err != nil {
			return fmt.Errorf("ошибка ввода следующего действия: %w", err)
		}

		switch option {
		case 1:
			if len(users) == config.PageSize {
				page++
			}
		case 2:
			if page > 1 {
				page--
			}
		case 0:
			return nil
		}
	}
}
