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

func AddCompany(a *app.App, args ...any) (err error) {
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

	err = GetActivityFields(a, args...)
	if err != nil {
		return fmt.Errorf("добавление компании: %w", err)
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

	err = a.CompSvc.Create(ctx, &company)
	if err != nil {
		return fmt.Errorf("ошибка добавления компании: %w", err)
	}

	return nil
}

func DeleteCompany(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	err = GetMyCompanies(a, args...)
	if err != nil {
		return fmt.Errorf("удаление компании: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Введите id компании: ")
	companyId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id компании: %w", err)
	}
	companyId = strings.TrimSpace(companyId)

	companyUuid, err := uuid.Parse(companyId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	err = a.CompSvc.DeleteById(ctx, companyUuid)
	if err != nil {
		return fmt.Errorf("удаление компании: %w", err)
	}

	return nil
}

func UpdateCompany(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	err = GetMyCompanies(a, args...)
	if err != nil {
		return fmt.Errorf("обновление информации о компании: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Введите id компании: ")
	companyId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id компании: %w", err)
	}
	companyId = strings.TrimSpace(companyId)

	companyUuid, err := uuid.Parse(companyId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	comp, err := a.CompSvc.GetById(ctx, companyUuid)
	if err != nil {
		return fmt.Errorf("получение компании по id: %w", err)
	}

	var name, city, activityFieldId string
	fmt.Printf("Введите название компании (%s): ", comp.Name)
	name, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода названия компании: %w", err)
	}
	name = strings.TrimSpace(name)
	if name != "" {
		comp.Name = name
	}

	fmt.Printf("Введите название города (%s): ", comp.City)
	city, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода названия города: %w", err)
	}
	city = strings.TrimSpace(city)
	if city != "" {
		comp.City = city
	}

	err = GetActivityFields(a, args...)
	if err != nil {
		return fmt.Errorf("обновление компании: %w", err)
	}

	field, err := a.ActFieldSvc.GetById(ctx, comp.ActivityFieldId)
	if err != nil {
		return fmt.Errorf("обновление компании: %w", err)
	}

	fmt.Printf("Введите id сферы деятельности (%s): ", field.Name)
	activityFieldId, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода сферы деятельности: %w", err)
	}
	activityFieldId = strings.TrimSpace(activityFieldId)

	if activityFieldId != "" {
		activityFieldUuid, err := uuid.Parse(activityFieldId)
		if err != nil {
			return fmt.Errorf("парсинг uuid из строки: %w", err)
		}

		comp.ActivityFieldId = activityFieldUuid
	}

	err = a.CompSvc.Update(ctx, comp)
	if err != nil {
		return fmt.Errorf("ошибка обновления компании: %w", err)
	}

	return nil
}

func GetMyCompanies(a *app.App, args ...any) (err error) {
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

	err = utils.PrintPaginatedCollectionArgs("Компании", a.CompSvc.GetByOwnerId, ctx, user.ID, true)
	if err != nil {
		return fmt.Errorf("вывод компаний предпринимателя с пагинацией: %w", err)
	}

	return nil
}
