package handlers

import (
	"bufio"
	"fmt"
	"os"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/config"
	"ppo/internal/tui/utils"
	"strings"
)

func getAllActivityFields(a *app.App) (err error) {
	page := 1
	for {
		fields, err := a.ActFieldSvc.GetAll(page)
		if err != nil {
			return fmt.Errorf("получение сфер деятельности: %w, err")
		}

		utils.PrintActivityFields(fields)

		fmt.Printf("1. Следующая страница.\n2. Предыдущая страница.\n0. Назад.\n\nВыберите действие: ")
		var option int
		_, err = fmt.Scanf("%d", &option)
		if err != nil {
			return fmt.Errorf("ошибка ввода следующего действия: %w", err)
		}

		switch option {
		case 1:
			if len(fields) == config.PageSize {
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

func AddActivityField(a *app.App, args ...any) (err error) {
	reader := bufio.NewReader(os.Stdin)

	var name, description string
	var cost float32
	fmt.Printf("Введите название cферы деятельности: ")
	_, err = fmt.Scanf("%s", &name)
	if err != nil {
		return fmt.Errorf("ошибка ввода названия сферы деятельности: %w", err)
	}

	fmt.Printf("Введите описание: ")
	description, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода описания: %w", err)
	}
	description = strings.TrimSpace(description)

	fmt.Printf("Введите вес сферы деятельности: ")
	_, err = fmt.Scanf("%f", &cost)
	if err != nil {
		return fmt.Errorf("ошибка ввода веса сферы деятельности: %w", err)
	}

	var field domain.ActivityField
	field.Name = name
	field.Description = description
	field.Cost = cost

	err = a.ActFieldSvc.Create(&field)
	if err != nil {
		return fmt.Errorf("ошибка добавления сферы деятельности: %w", err)
	}

	return nil
}
