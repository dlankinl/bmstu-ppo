package tui

import (
	"fmt"
	"ppo/domain"
	"ppo/internal/app"
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

		//if choice == 0 {
		//	return nil
		//}
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
			fmt.Println("Введите логин:")
			_, err = fmt.Scanf("%s", &login)
			if err != nil {
				return fmt.Errorf("ошибка ввода логина")
			}

			fmt.Println("Введите пароль:")
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
	fmt.Println(guestPrompt)

	var choice int
	switch choice {
	case 1:
		var login, password string

		fmt.Println("Введите логин:")
		_, err = fmt.Scanf("%s", &login)
		if err != nil {
			return fmt.Errorf("ошибка ввода логина")
		}

		fmt.Println("Введите пароль:")
		_, err = fmt.Scanf("%s", &password)
		if err != nil {
			return fmt.Errorf("ошибка ввода пароля")
		}

		ua := &domain.UserAuth{Username: login, Password: password}
		err = a.AuthSvc.Register(ua)
		if err != nil {
			return fmt.Errorf("ошибка регистрации: %w", err)
		}
	case 0:
		return nil
	}

	return nil
}
