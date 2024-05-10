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
	"strconv"
	"strings"
)

func AddActivityField(a *app.App, args ...any) (err error) {
	ctx := context.Background()
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

	err = a.ActFieldSvc.Create(ctx, &field)
	if err != nil {
		return fmt.Errorf("ошибка добавления сферы деятельности: %w", err)
	}

	return nil
}

func DeleteActivityField(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	err = utils.PrintPaginatedCollection("Сферы деятельности", a.ActFieldSvc.GetAll, ctx)
	if err != nil {
		return fmt.Errorf("вывод сфер деятельности с пагинацией: %w", err)
	}

	fmt.Printf("Введите id сферы деятельности: ")
	activityFieldId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода сферы деятельности: %w", err)
	}
	activityFieldId = strings.TrimSpace(activityFieldId)

	activityFieldUuid, err := uuid.Parse(activityFieldId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	err = a.ActFieldSvc.DeleteById(ctx, activityFieldUuid)
	if err != nil {
		return fmt.Errorf("удаление сферы деятельности: %w", err)
	}

	return nil
}

func UpdateActivityField(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	err = utils.PrintPaginatedCollection("Сферы деятельности", a.ActFieldSvc.GetAll, ctx)
	if err != nil {
		return fmt.Errorf("вывод сфер деятельности с пагинацией: %w", err)
	}

	fmt.Printf("Введите id сферы деятельности: ")
	activityFieldId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода сферы деятельности: %w", err)
	}
	activityFieldId = strings.TrimSpace(activityFieldId)

	activityFieldUuid, err := uuid.Parse(activityFieldId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	field, err := a.ActFieldSvc.GetById(ctx, activityFieldUuid)
	if err != nil {
		return fmt.Errorf("получение сферы деятельности по id: %w", err)
	}

	var name, description, costStr string
	var cost float64
	fmt.Printf("Введите название cферы деятельности (%s): ", field.Name)
	name, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода названия сферы деятельности: %w", err)
	}
	name = strings.TrimSpace(name)
	if name != "" {
		field.Name = name
	}

	fmt.Printf("Введите описание (%s): ", field.Description)
	description, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода описания: %w", err)
	}
	description = strings.TrimSpace(description)
	if description != "" {
		field.Description = description
	}

	fmt.Printf("Введите вес сферы деятельности (%f): ", field.Cost)
	costStr, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода веса сферы деятельности: %w", err)
	}
	costStr = strings.TrimSpace(costStr)
	if costStr != "" {
		cost, err = strconv.ParseFloat(costStr, 32)
		if err != nil {
			return fmt.Errorf("ошибка конвертации во float: %w", err)
		}

		field.Cost = float32(cost)
	}

	err = a.ActFieldSvc.Update(ctx, field)
	if err != nil {
		return fmt.Errorf("ошибка обновления сферы деятельности: %w", err)
	}

	return nil
}

func GetActivityFields(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	err = utils.PrintPaginatedCollection("Сферы деятельности", a.ActFieldSvc.GetAll, ctx)
	if err != nil {
		return fmt.Errorf("вывод сфер деятельности с пагинацией: %w", err)
	}

	return nil
}
