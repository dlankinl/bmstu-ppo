package handlers

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
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

func AddCompany(a *app.App, username string) (err error) {
	reader := bufio.NewReader(os.Stdin)

	user, err := a.UserSvc.GetByUsername(username)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	var name, activityFieldId, city string
	fmt.Printf("Введите название компании: ")
	_, err = fmt.Scanf("%s", &name)
	if err != nil {
		return fmt.Errorf("ошибка ввода названия компании: %w", err)
	}

	fmt.Printf("Введите город: ")
	city, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода города: %w", err)
	}
	city = strings.TrimSpace(city)

	fmt.Printf("Введите id сферы деятельности (-1 - вывести список):")
	activityFieldId, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода сферы деятельности: %w", err)
	}
	if activityFieldId == "-1" {
		err = getAllActivityFields(a)
		if err != nil {
			return fmt.Errorf("вывод сфер деятельности с пагинацией: %w", err)
		}
	}

	activityFieldUuid, err := uuid.Parse(activityFieldId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	var company domain.Company
	company.Name = name
	company.OwnerID = user.ID
	company.City = city
	company.ActivityFieldId = activityFieldUuid

	err = a.CompSvc.Create(&company)
	if err != nil {
		return fmt.Errorf("ошибка добавления информации в карточку предпринимателя: %w", err)
	}

	return nil
}
