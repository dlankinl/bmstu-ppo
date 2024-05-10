package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"ppo/internal/app"
	"strings"
	"time"
)

func UpdateUser(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	var username, fullName, gender, birthdayStr, city, role string
	fmt.Printf("Введите имя пользователя: ")
	_, err = fmt.Scanf("%s", &username)
	if err != nil {
		return fmt.Errorf("ошибка ввода имени пользователя: %w", err)
	}

	user, err := a.UserSvc.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("пользователь не найден")
	}

	fmt.Printf("Введите полное имя (%s): ", user.FullName)
	fullName, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода полного имени: %w", err)
	}
	fullName = strings.TrimSpace(fullName)
	if fullName != "" {
		user.FullName = fullName
	}

	fmt.Printf("Введите дату рождения в формате ГГГГ-ММ-ДД (%s): ", user.Birthday.Format("2006-01-02"))
	birthdayStr, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода даты рождения")
	}
	birthdayStr = strings.TrimSpace(birthdayStr)
	if birthdayStr != "" {
		bday, err := time.Parse("2006-01-02", birthdayStr)
		if err != nil {
			return fmt.Errorf("ошибка перевода даты рождения в time.Time: %w", err)
		}
		user.Birthday = bday
	}

	fmt.Printf("Введите пол (%s): ", user.Gender)
	gender, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода пола: %w", err)
	}
	gender = strings.TrimSpace(gender)
	if gender != "" {
		user.Gender = gender
	}

	fmt.Printf("Введите город (%s): ", user.City)
	city, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода города: %w", err)
	}
	city = strings.TrimSpace(city)
	if city != "" {
		user.City = city
	}

	fmt.Printf("Введите роль (%s): ", user.Role)
	role, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода роли: %w", err)
	}
	role = strings.TrimSpace(role)
	if role != "" {
		user.Role = role
	}

	err = a.UserSvc.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("ошибка обновления карточки предпринимателя: %w", err)
	}

	return nil
}

func CreateUser(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	var username, fullName, gender, birthdayStr, city, role string
	fmt.Printf("Введите имя пользователя: ")
	_, err = fmt.Scanf("%s", &username)
	if err != nil {
		return fmt.Errorf("ошибка ввода имени пользователя: %w", err)
	}

	user, err := a.UserSvc.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	fmt.Printf("Введите полное имя: ")
	fullName, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода полного имени: %w", err)
	}
	fullName = strings.TrimSpace(fullName)

	fmt.Printf("Введите дату рождения в формате ГГГГ-ММ-ДД: ")
	birthdayStr, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода даты рождения")
	}
	birthdayStr = strings.TrimSpace(birthdayStr)
	var bday time.Time
	if birthdayStr != "" {
		bday, err = time.Parse("2006-01-02", birthdayStr)
		if err != nil {
			return fmt.Errorf("ошибка перевода даты рождения в time.Time: %w", err)
		}
	}

	fmt.Printf("Введите пол: ")
	gender, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода пола: %w", err)
	}
	gender = strings.TrimSpace(gender)

	fmt.Printf("Введите город: ")
	city, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода города: %w", err)
	}
	city = strings.TrimSpace(city)

	fmt.Printf("Введите роль: ")
	role, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода роли: %w", err)
	}
	role = strings.TrimSpace(role)

	user.Gender = gender
	user.City = city
	user.Birthday = bday
	user.FullName = fullName
	user.Role = role

	err = a.UserSvc.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("ошибка добавления информации в карточку предпринимателя: %w", err)
	}

	return nil
}

func ChangeUserRole(a *app.App, args ...any) (err error) {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	var username, role string
	fmt.Printf("Введите имя пользователя: ")
	_, err = fmt.Scanf("%s", &username)
	if err != nil {
		return fmt.Errorf("ошибка ввода имени пользователя: %w", err)
	}

	user, err := a.UserSvc.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("пользователь не найден: %w", err)
	}

	fmt.Printf("Введите роль: ")
	role, _ = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода роли: %w", err)
	}
	role = strings.TrimSpace(role)

	user.Role = role

	err = a.UserSvc.Update(ctx, user)
	if err != nil {
		return fmt.Errorf("изменение роли пользователя: %w", err)
	}

	return nil
}
