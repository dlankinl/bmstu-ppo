package tui

import (
	"fmt"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/config"
	"ppo/internal/tui/utils"
	"time"
)

func Run(a *app.App) (err error) {
	var choice int
	for {
		fmt.Println(authPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("ошибка ввода: %w", err)
			continue
		}

		switch choice {
		case 0:
			return nil
		case 1:
			err = guestMenu(a)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case 2:
			var login, password string
			fmt.Printf("Введите логин: ")
			_, err = fmt.Scanf("%s", &login)
			if err != nil {
				return fmt.Errorf("ошибка ввода логина")
			}

			fmt.Printf("Введите пароль: ")
			_, err = fmt.Scanf("%s", &password)
			if err != nil {
				return fmt.Errorf("ошибка ввода пароля")
			}

			ua := &domain.UserAuth{Username: login, Password: password}
			token, err := a.AuthSvc.Login(ua)
			if err != nil {
				return fmt.Errorf("ошибка авторизации: %w", err)
			}

			fmt.Println(token)
		}
	}
}

func guestMenu(a *app.App) (err error) {
	var choice int
	for {
		fmt.Println(guestPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("ошибка ввода: %w", err)
		}

		switch choice {
		case 1:
			var login, password string

			fmt.Printf("Введите логин: ")
			_, err = fmt.Scanf("%s", &login)
			if err != nil {
				return fmt.Errorf("ошибка ввода логина")
			}

			fmt.Printf("Введите пароль: ")
			_, err = fmt.Scanf("%s", &password)
			if err != nil {
				return fmt.Errorf("ошибка ввода пароля")
			}

			ua := &domain.UserAuth{Username: login, Password: password}
			err = a.AuthSvc.Register(ua)
			if err != nil {
				return fmt.Errorf("ошибка регистрации: %w", err)
			}
		case 2:
			page := 1
			for {
				users, err := a.UserSvc.GetAll(page)
				if err != nil {
					return fmt.Errorf("получение пользователей: %w, err")
				}

				if len(users) < config.PageSize {
					break
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
		case 0:
			return nil
		}
	}
}

func adminMenu(a *app.App) (err error) {
	var choice int
	for {
		fmt.Println(adminPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("ошибка ввода: %w", err)
		}

		switch choice {
		case 1:
			var fullName, birthdayStr, city, gender, username string

			fmt.Printf("Введите имя пользователя:")
			_, err = fmt.Scanf("%s", &username)
			if err != nil {
				return fmt.Errorf("ошибка ввода имени пользователя")
			}

			fmt.Printf("Введите ФИО через пробел:")
			_, err = fmt.Scanf("%s", &username)
			if err != nil {
				return fmt.Errorf("ошибка ввода имени пользователя")
			}

			fmt.Printf("Введите дату рождения (ГГГГ-ММ-ДД): ")
			_, err = fmt.Scanf("%s", &birthdayStr)
			if err != nil {
				return fmt.Errorf("ошибка ввода пароля")
			}

			bday, err := time.Parse("2006-01-02", birthdayStr)
			if err != nil {
				return fmt.Errorf("ошибка перевода даты рождения в time.Time: %w", err)
			}

			fmt.Printf("Введите пол (m - мужской, f - женский): ")
			_, err = fmt.Scanf("%s", &gender)
			if err != nil {
				return fmt.Errorf("ошибка ввода пола")
			}

			fmt.Printf("Введите город: ")
			_, err = fmt.Scanf("%s", &city)
			if err != nil {
				return fmt.Errorf("ошибка ввода пола")
			}

			user := &domain.User{
				FullName: fullName,
				Username: username,
				Birthday: bday,
				Gender:   gender,
				City:     city,
			}
			err = a.UserSvc.Update(user)
			if err != nil {
				return fmt.Errorf("ошибка заполнения карточки предпринимателя: %w", err)
			} else {
				fmt.Println("Карточка предпринимателя заполнена успешно")
			}
		// TODO: сделать возможность оставить текущее значение (как в утилитах "y/N")
		case 2:
			page := 1
			for {
				users, err := a.UserSvc.GetAll(page)
				if err != nil {
					return fmt.Errorf("получение пользователей: %w, err")
				}

				if len(users) < config.PageSize {
					break
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
		case 0:
			return nil
		}
	}
}
