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
	Func func(*app.App, ...any) error
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
		Func: handlers.UpdateUser,
	},
	{
		Role: admin,
		Name: "Создать карточку предпринимателя",
		Func: handlers.CreateUser,
	},
	{
		Role: user,
		Name: "Добавить компанию",
		Func: handlers.AddCompany,
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
		Func: handlers.GetAllUsers,
	},
	{
		Role: admin,
		Name: "Добавить сферу деятельности",
		Func: handlers.AddActivityField,
	},
}

func (t *TUI) Run() (err error) {
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
			err = t.guestMenu()
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
			token, err := t.app.AuthSvc.Login(ua)
			if err != nil {
				return fmt.Errorf("ошибка авторизации: %w", err)
			}

			payload, err := base.VerifyAuthToken(token, t.app.Config.JwtKey)
			if err != nil {
				return fmt.Errorf("ошибка верификации JWT токена: %w", err)
			}

			t.userInfo = payload
			err = t.userMenu()
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func (t *TUI) userMenu() (err error) {
	for {
		actions, prompt := generateActionsPrompt(t.userInfo.Role)
		fmt.Println(prompt)
		fmt.Println(actions)

		var choice int
		_, err = fmt.Scanf("%d", &choice)

		switch choice {
		case 0:
			return nil
		}
	}
	//var choice int
	//for {
	//	fmt.Println(userPrompt)
	//	_, err = fmt.Scanf("%d", &choice)
	//	if err != nil {
	//		fmt.Printf("ошибка ввода: %v", err)
	//	}
	//
	//	switch choice {
	//	case 0:
	//		return nil
	//	case 1:
	//		err = handlers.GetAllUsers(t.app)
	//		if err != nil {
	//			return fmt.Errorf("ошибка просмотра списка предпринимателей: %w", err)
	//		}
	//	case 2:
	//		err = handlers.CalculateRating(t.app)
	//		if err != nil {
	//			return fmt.Errorf("ошибка вычисления рейтинга предпринимателя: %w", err)
	//		}
	//	case 3:
	//		fmt.Println(t.userInfo)
	//		err = t.companiesMenu()
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}
}

func (t *TUI) guestMenu() (err error) {
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
			err = t.app.AuthSvc.Register(ua)
			if err != nil {
				return fmt.Errorf("ошибка регистрации: %w", err)
			}
			return nil
		case 2:
			err = handlers.GetAllUsers(t.app)
			if err != nil {
				return fmt.Errorf("ошибка просмотра списка предпринимателей: %w", err)
			}
		case 0:
			return nil
		}
	}
}

func (t *TUI) adminMenu() (err error) {
	var choice int
	for {
		fmt.Println(adminPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("ошибка ввода: %w", err)
		}

		switch choice {
		case 1:
			err = handlers.CreateUser(t.app)
			if err != nil {
				return fmt.Errorf("ошибка заполнения карточки предпринимателя: %w", err)
			} else {
				fmt.Println("Карточка предпринимателя заполнена успешно")
			}
		case 2:
			err = handlers.UpdateUser(t.app)
			if err != nil {
				return err
			}
		case 3:

		case 0:
			return nil
		}
	}
}

func (t *TUI) companiesMenu() (err error) {
	var choice int
	for {
		fmt.Println(companiesPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			fmt.Println("ошибка ввода: %w", err)
		}

		switch choice {
		case 1:
			err = handlers.AddCompany(t.app, t.userInfo.Username)
			if err != nil {
				return fmt.Errorf("ошибка добавления компании: %w", err)
			} else {
				fmt.Println("Компания успешно добавлена")
			}
		case 2:
			err = handlers.UpdateUser(t.app)
			if err != nil {
				return err
			}
		case 3:

		case 0:
			return nil
		}
	}
}

func generateActionsPrompt(role string) (actionsList []*Action, actionsPrompt string) {
	actionsList = make([]*Action, 0)

	j := 1
	for _, action := range actions {
		if role == action.Role {
			actionsList = append(actionsList, &action)
			actionsPrompt += fmt.Sprintf("\n%d. %s", j, action.Name)
			j++
		}
	}
	actionsPrompt += "\n\nВыберите действие: "

	return actionsList, actionsPrompt
}
