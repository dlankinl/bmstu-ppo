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

func AddContact(a *app.App, args ...any) (err error) {
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

	var name, value string
	fmt.Printf("Введите название средства связи: ")
	name, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода названия средства связи: %w", err)
	}

	fmt.Printf("Введите значение: ")
	value, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода значения: %w", err)
	}
	value = strings.TrimSpace(value)

	var contact domain.Contact
	contact.OwnerID = user.ID
	contact.Name = name
	contact.Value = value

	err = a.ConSvc.Create(ctx, &contact)
	if err != nil {
		return fmt.Errorf("ошибка добавления средства связи: %w", err)
	}

	return nil
}

func DeleteContact(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	err = GetMyContacts(a, args...)
	if err != nil {
		return fmt.Errorf("удаление средства связи: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Введите id средства связи: ")
	conId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id навыка: %w", err)
	}
	conId = strings.TrimSpace(conId)

	conUuid, err := uuid.Parse(conId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	err = a.ConSvc.DeleteById(ctx, conUuid)
	if err != nil {
		return fmt.Errorf("удаление средства связи: %w", err)
	}

	return nil
}

func UpdateContact(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	err = GetMyContacts(a, args...)
	if err != nil {
		return fmt.Errorf("обновление информации о средстве связи: %w", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Введите id средства связи: ")
	conId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id навыка: %w", err)
	}
	conId = strings.TrimSpace(conId)

	conUuid, err := uuid.Parse(conId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	contact, err := a.ConSvc.GetById(ctx, conUuid)
	if err != nil {
		return fmt.Errorf("получение средства связи по id: %w", err)
	}

	var name, value string
	fmt.Printf("Введите название средства связи (%s): ", contact.Name)
	name, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода названия навыка: %w", err)
	}
	name = strings.TrimSpace(name)
	if name != "" {
		contact.Name = name
	}

	fmt.Printf("Введите значение (%s): ", contact.Value)
	value, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода описания: %w", err)
	}
	value = strings.TrimSpace(value)
	if value != "" {
		contact.Value = value
	}

	err = a.ConSvc.Update(ctx, contact)
	if err != nil {
		return fmt.Errorf("ошибка обновления средства связи: %w", err)
	}

	return nil
}

func GetMyContacts(a *app.App, args ...any) (err error) {
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

	err = utils.PrintPaginatedCollectionArgs("Средства связи", a.ConSvc.GetByOwnerId, ctx, user.ID, true)
	if err != nil {
		return fmt.Errorf("вывод средств связи с пагинацией: %w", err)
	}

	return nil
}
