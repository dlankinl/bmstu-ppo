package tui

import (
	"bufio"
	"fmt"
	"os"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/tui/handlers"
	"ppo/pkg/base"
)

const (
	admin = "admin"
	user  = "user"
)

type Action struct {
	Role string
	Name string
	Func func(...any) any
}

type TUI struct {
	userInfo *base.JwtPayload
	app      *app.App
}

func NewTUI(app *app.App) *TUI {
	return &TUI{
		app: app,
	}
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
	reader := bufio.NewReader(os.Stdin)
	_ = reader
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

			payload, err := base.VerifyAuthToken(token, a.Config.JwtKey)
			if err != nil {
				return fmt.Errorf("ошибка верификации JWT токена: %w", err)
			}

			if payload.Role == admin {
				err = adminMenu(a, payload)
			} else if payload.Role == user {
				err = userMenu(a, payload)
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			//err = userMenu(a, role)
			//if err != nil {
			//	fmt.Println(err)
			//	continue
			//}
		}
	}
}

func userMenu(a *app.App, payload *base.JwtPayload) (err error) {
	//for {
	//	actionsPrompt := generateActionsPrompt(role)
	//	fmt.Print(actionsPrompt)
	//
	//	var choice int
	//	_, err = fmt.Scanf("%d", &choice)
	//
	//	switch choice {
	//	case 0:
	//		return nil
	//	}
	//}
	var choice int
	for {
		fmt.Println(userPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Printf("ошибка ввода: %v", err)
		}

		switch choice {
		case 0:
			return nil
		case 1:
			err = handlers.GetAllUsers(a)
			if err != nil {
				return fmt.Errorf("ошибка просмотра списка предпринимателей: %w", err)
			}
		case 2:
			err = handlers.CalculateRating(a)
			if err != nil {
				return fmt.Errorf("ошибка вычисления рейтинга предпринимателя: %w", err)
			}
		}
	}
}

func guestMenu(a *app.App) (err error) {
	var choice int
	for {
		fmt.Println(guestPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Printf("ошибка ввода: %v", err)
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
			return nil
		case 2:
			err = handlers.GetAllUsers(a)
			if err != nil {
				return fmt.Errorf("ошибка просмотра списка предпринимателей: %w", err)
			}
		case 0:
			return nil
		}
	}
}

func adminMenu(a *app.App, payload *base.JwtPayload) (err error) {
	var choice int
	for {
		fmt.Println(adminPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("ошибка ввода: %w", err)
		}

		switch choice {
		case 1:
			err = handlers.CreateUser(a)
			if err != nil {
				return fmt.Errorf("ошибка заполнения карточки предпринимателя: %w", err)
			} else {
				fmt.Println("Карточка предпринимателя заполнена успешно")
			}
		case 2:
			err = handlers.UpdateUser(a)
			if err != nil {
				return err
			}
		case 3:

		case 0:
			return nil
		}
	}
}

func companiesMenu(a *app.App, payload *base.JwtPayload) (err error) {
	var choice int
	for {
		fmt.Println(companiesPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("ошибка ввода: %w", err)
		}

		switch choice {
		case 1:
			err = handlers.CreateUser(a)
			if err != nil {
				return fmt.Errorf("ошибка заполнения карточки предпринимателя: %w", err)
			} else {
				fmt.Println("Карточка предпринимателя заполнена успешно")
			}
		case 2:
			err = handlers.UpdateUser(a)
			if err != nil {
				return err
			}
		case 3:

		case 0:
			return nil
		}
	}
}

func generateActionsPrompt(role string) (actionsPrompt string) {
	allowedActions := make([]Action, 0)

	j := 1
	for _, action := range actions {
		if role == action.Role {
			allowedActions = append(allowedActions, action)
			actionsPrompt += fmt.Sprintf("\n%d. %s", j, action.Name)
			j++
		}
	}
	actionsPrompt += "\n\nВыберите действие: "

	return actionsPrompt
}
