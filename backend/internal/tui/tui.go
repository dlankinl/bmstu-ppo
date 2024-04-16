package tui

import (
	"bufio"
	"fmt"
	"os"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/config"
	"ppo/internal/tui/utils"
	"ppo/pkg/base"
	"strings"
	"time"
)

const (
	admin = "admin"
	user  = "user"
)

//func Run(a *app.App) (err error) {
//	var choice int
//	for {
//		fmt.Println(authPrompt)
//		_, err = fmt.Scanf("%d", &choice)
//		if err != nil {
//			fmt.Println("ошибка ввода: %w", err)
//			continue
//		}
//
//		switch choice {
//		case 0:
//			return nil
//		case 1:
//			err = guestMenu(a)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//		case 2:
//			var login, password string
//			fmt.Printf("Введите логин: ")
//			_, err = fmt.Scanf("%s", &login)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода логина")
//			}
//
//			fmt.Printf("Введите пароль: ")
//			_, err = fmt.Scanf("%s", &password)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода пароля")
//			}
//
//			ua := &domain.UserAuth{Username: login, Password: password}
//			token, err := a.AuthSvc.Login(ua)
//			if err != nil {
//				return fmt.Errorf("ошибка авторизации: %w", err)
//			}
//
//			fmt.Println(token)
//		case 3:
//			var login, password string
//			fmt.Printf("Введите логин: ")
//			_, err = fmt.Scanf("%s", &login)
//			if err != nil {
//				fmt.Println("Ошибка ввода логина")
//				continue
//			}
//
//			fmt.Printf("Введите пароль: ")
//			_, err = fmt.Scanf("%s", &password)
//			if err != nil {
//				fmt.Println("Ошибка ввода пароля")
//				continue
//			}
//
//			ua := &domain.UserAuth{Username: login, Password: password}
//			token, err := a.AuthSvc.Login(ua)
//			if err != nil {
//				fmt.Println("Ошибка авторизации: %v", err)
//				continue
//			}
//
//			fmt.Println("TOKEN: ", token)
//			err = adminMenu(a)
//			if err != nil {
//				fmt.Println(err)
//				continue
//			}
//		}
//	}
//}
//
//func guestMenu(a *app.App) (err error) {
//	var choice int
//	for {
//		fmt.Println(guestPrompt)
//		_, err = fmt.Scanf("%d", &choice)
//		if err != nil {
//			fmt.Println("ошибка ввода: %w", err)
//		}
//
//		switch choice {
//		case 1:
//			var login, password string
//
//			fmt.Printf("Введите логин: ")
//			_, err = fmt.Scanf("%s", &login)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода логина")
//			}
//
//			fmt.Printf("Введите пароль: ")
//			_, err = fmt.Scanf("%s", &password)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода пароля")
//			}
//
//			ua := &domain.UserAuth{Username: login, Password: password}
//			err = a.AuthSvc.Register(ua)
//			if err != nil {
//				return fmt.Errorf("ошибка регистрации: %w", err)
//			}
//		case 2:
//			page := 1
//			for {
//				users, err := a.UserSvc.GetAll(page)
//				if err != nil {
//					return fmt.Errorf("получение пользователей: %w, err")
//				}
//
//				if len(users) < config.PageSize {
//					break
//				} else {
//					fmt.Println("Выберите действие:\n1. Следующая страница.\n2. Предыдущая страница.\n0. Назад.")
//					var option int
//					_, err = fmt.Scanf("%d", &option)
//					if err != nil {
//						return fmt.Errorf("ошибка ввода следующего действия: %w", err)
//					}
//
//					switch option {
//					case 1:
//						page++
//					case 2:
//						if page > 1 {
//							page--
//						}
//					case 0:
//						break
//					}
//				}
//
//				utils.PrintUsers(users)
//			}
//		case 0:
//			return nil
//		}
//	}
//}
//
//func adminMenu(a *app.App) (err error) {
//	reader := bufio.NewReader(os.Stdin)
//	var choice int
//	for {
//		fmt.Println(adminPrompt)
//		_, err = fmt.Scanf("%d", &choice)
//		if err != nil {
//			fmt.Println("ошибка ввода: %w", err)
//		}
//
//		switch choice {
//		case 1:
//			var fullName, birthdayStr, city, gender, username, lastName, midName, firstName string
//
//			fmt.Printf("Введите имя пользователя:")
//			_, err = fmt.Scanf("%s", &username)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода имени пользователя")
//			}
//
//			fmt.Printf("Введите фамилию:")
//			_, err = fmt.Scanf("%s", &lastName)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода фамилии")
//			}
//
//			fmt.Printf("Введите имя:")
//			_, err = fmt.Scanf("%s", &firstName)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода имени")
//			}
//
//			fmt.Printf("Введите отчество:")
//			_, err = fmt.Scanf("%s", &midName)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода отчества")
//			}
//
//			fullName = strings.Join([]string{lastName, firstName, midName}, " ")
//
//			fmt.Printf("Введите дату рождения (ГГГГ-ММ-ДД): ")
//			_, err = fmt.Scanf("%s", &birthdayStr)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода даты рождения")
//			}
//
//			bday, err := time.Parse("2006-01-02", birthdayStr)
//			if err != nil {
//				return fmt.Errorf("ошибка перевода даты рождения в time.Time: %w", err)
//			}
//
//			fmt.Printf("Введите пол (m - мужской, f - женский): ")
//			_, err = fmt.Scanf("%s", &gender)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода пола")
//			}
//
//			fmt.Printf("Введите город: ")
//			_, err = fmt.Scanf("%s", &city)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода пола")
//			}
//
//			user := &domain.User{
//				FullName: fullName,
//				Username: username,
//				Birthday: bday,
//				Gender:   gender,
//				City:     city,
//			}
//			err = a.UserSvc.Update(user)
//			if err != nil {
//				return fmt.Errorf("ошибка заполнения карточки предпринимателя: %w", err)
//			} else {
//				fmt.Println("Карточка предпринимателя заполнена успешно")
//			}
//		case 2:
//			var username, fullName, gender, birthdayStr, city string
//			fmt.Printf("Введите имя пользователя: ")
//			_, err = fmt.Scanf("%s", &username)
//			if err != nil {
//				return fmt.Errorf("ошибка ввода имени пользователя: %w", err)
//			}
//
//			user, err := a.UserSvc.GetByUsername(username)
//			if err != nil {
//				return fmt.Errorf("пользователь не найден")
//			}
//
//			fmt.Printf("Введите полное имя (%s): ", user.FullName)
//			fullName, _ = reader.ReadString('\n')
//			if err != nil {
//				return fmt.Errorf("ошибка ввода полного имени: %w", err)
//			}
//			fullName = strings.TrimSpace(fullName)
//			if fullName != "" {
//				user.FullName = fullName
//			}
//
//			fmt.Printf("Введите дату рождения в формате ГГГГ-ММ-ДД (%s): ", user.Birthday.Format("2006-01-02"))
//			birthdayStr, _ = reader.ReadString('\n')
//			if err != nil {
//				return fmt.Errorf("ошибка ввода даты рождения")
//			}
//			birthdayStr = strings.TrimSpace(birthdayStr)
//			if birthdayStr != "" {
//				bday, err := time.Parse("2006-01-02", birthdayStr)
//				if err != nil {
//					return fmt.Errorf("ошибка перевода даты рождения в time.Time: %w", err)
//				}
//				user.Birthday = bday
//			}
//
//			fmt.Printf("Введите пол (%s): ", user.Gender)
//			gender, _ = reader.ReadString('\n')
//			if err != nil {
//				return fmt.Errorf("ошибка ввода пола: %w", err)
//			}
//			gender = strings.TrimSpace(gender)
//			if gender != "" {
//				user.Gender = gender
//			}
//
//			fmt.Printf("Введите город (%s): ", user.City)
//			city, _ = reader.ReadString('\n')
//			if err != nil {
//				return fmt.Errorf("ошибка ввода города: %w", err)
//			}
//			city = strings.TrimSpace(city)
//			if city != "" {
//				user.City = city
//			}
//
//			err = a.UserSvc.Update(user)
//			if err != nil {
//				return fmt.Errorf("ошибка обновления карточки предпринимателя: %w", err)
//			}
//		case 0:
//			return nil
//		}
//	}
//}

type Action struct {
	Role string
	Name string
}

var actions = []Action{
	{
		Role: admin,
		Name: "Редактировать карточку предпринимателя",
	},
	{
		Role: admin,
		Name: "Создать карточку предпринимателя",
	},
	{
		Role: admin,
		Name: "Удалить карточку предпринимателя",
	},
	{
		Role: admin,
		Name: "Сменить роль пользователя",
	},
	{
		Role: user,
		Name: "Просмотреть список предпринимателей",
	},
}

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

			role, err := base.VerifyAuthToken(token, a.Config.JwtKey)
			if err != nil {
				return fmt.Errorf("ошибка верификации JWT токена: %w", err)
			}

			allowedActions := make([]Action, 0)
			actionsPrompt := "Выберите действие:\n\n"
			j := 1
			for _, action := range actions {
				if role == action.Role {
					allowedActions = append(allowedActions, action)
					actionsPrompt += fmt.Sprintf("%d. %s\n", j, action.Name)
					j++
				}
			}

			fmt.Println(actionsPrompt)
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
	reader := bufio.NewReader(os.Stdin)
	var choice int
	for {
		fmt.Println(adminPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("ошибка ввода: %w", err)
		}

		switch choice {
		case 1:
			var fullName, birthdayStr, city, gender, username, lastName, midName, firstName string

			fmt.Printf("Введите имя пользователя:")
			_, err = fmt.Scanf("%s", &username)
			if err != nil {
				return fmt.Errorf("ошибка ввода имени пользователя")
			}

			fmt.Printf("Введите фамилию:")
			_, err = fmt.Scanf("%s", &lastName)
			if err != nil {
				return fmt.Errorf("ошибка ввода фамилии")
			}

			fmt.Printf("Введите имя:")
			_, err = fmt.Scanf("%s", &firstName)
			if err != nil {
				return fmt.Errorf("ошибка ввода имени")
			}

			fmt.Printf("Введите отчество:")
			_, err = fmt.Scanf("%s", &midName)
			if err != nil {
				return fmt.Errorf("ошибка ввода отчества")
			}

			fullName = strings.Join([]string{lastName, firstName, midName}, " ")

			fmt.Printf("Введите дату рождения (ГГГГ-ММ-ДД): ")
			_, err = fmt.Scanf("%s", &birthdayStr)
			if err != nil {
				return fmt.Errorf("ошибка ввода даты рождения")
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
		case 2:
			var username, fullName, gender, birthdayStr, city string
			fmt.Printf("Введите имя пользователя: ")
			_, err = fmt.Scanf("%s", &username)
			if err != nil {
				return fmt.Errorf("ошибка ввода имени пользователя: %w", err)
			}

			user, err := a.UserSvc.GetByUsername(username)
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

			err = a.UserSvc.Update(user)
			if err != nil {
				return fmt.Errorf("ошибка обновления карточки предпринимателя: %w", err)
			}
		case 0:
			return nil
		}
	}
}
