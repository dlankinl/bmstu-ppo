package handlers

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"os"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/tui/utils"
	"strings"
)

//func getPaginatedCompanies()

func AddCompany(a *app.App, args ...any) (err error) {
	reader := bufio.NewReader(os.Stdin)

	var username string
	var ok bool
	if len(args) > 0 {
		username, ok = args[0].(string)
		if !ok {
			return fmt.Errorf("приведение аргумента к string")
		}
	}

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

	err = utils.PrintPaginatedCollection("Сферы деятельности", a.ActFieldSvc.GetAll)
	if err != nil {
		return fmt.Errorf("вывод сфер деятельности с пагинацией: %w", err)
	}

	fmt.Printf("Введите id сферы деятельности: ")
	activityFieldId, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода сферы деятельности: %w", err)
	}
	activityFieldId = strings.TrimSpace(activityFieldId)

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
		return fmt.Errorf("ошибка добавления компании: %w", err)
	}

	return nil
}

func GetMyCompanies(a *app.App, args ...any) (err error) {
	var username string
	var ok bool
	if len(args) > 0 {
		username, ok = args[0].(string)
		if !ok {
			return fmt.Errorf("приведение аргумента к string")
		}
	}

	user, err := a.UserSvc.GetByUsername(username)
	if err != nil {
		return fmt.Errorf("пользователь не найден")
	}

	err = utils.PrintPaginatedCollectionArgs("Компании", a.CompSvc.GetByOwnerId, user.ID)
	if err != nil {
		return fmt.Errorf("вывод компаний предпринимателя с пагинацией: %w", err)
	}

	return nil
}
