package tui

import (
	"bufio"
	"business-logic/domain"
	"fmt"
	"os"
	"ppo/app/internal/app"
	handlers2 "ppo/app/internal/tui/handlers"
	"ppo/app/internal/tui/utils"
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
		Name: "[ Предприниматели ] Редактировать карточку предпринимателя",
		Func: handlers2.UpdateUser,
	},
	{
		Role: admin,
		Name: "[ Предприниматели ] Создать карточку предпринимателя",
		Func: handlers2.CreateUser,
	},
	{
		Role: user,
		Name: "[ Компании ] Добавить компанию",
		Func: handlers2.AddCompany,
	},
	{
		Role: user,
		Name: "[ Компании ] Удалить компанию",
		Func: handlers2.DeleteCompany,
	},
	{
		Role: user,
		Name: "[ Компании ] Обновить информацию о компании",
		Func: handlers2.UpdateCompany,
	},
	{
		Role: user,
		Name: "[ Компании ] Посмотреть список своих компаний",
		Func: handlers2.GetMyCompanies,
	},
	{
		Role: user,
		Name: "[ Предприниматели ] Просчитать рейтинг",
		Func: handlers2.CalculateRating,
	},
	{
		Role: admin,
		Name: "[ Предприниматели ] Сменить роль пользователя",
		Func: handlers2.ChangeUserRole,
	},
	{
		Role: user,
		Name: "[ Предприниматели ] Просмотреть список предпринимателей",
		Func: handlers2.GetAllUsers,
	},
	{
		Role: admin,
		Name: "[ Сферы деятельности ] Добавить сферу деятельности",
		Func: handlers2.AddActivityField,
	},
	{
		Role: admin,
		Name: "[ Сферы деятельности ] Удалить сферу деятельности",
		Func: handlers2.DeleteActivityField,
	},
	{
		Role: admin,
		Name: "[ Сферы деятельности ] Редактировать сферу деятельности",
		Func: handlers2.UpdateActivityField,
	},
	{
		Role: admin,
		Name: "[ Навыки ] Добавить навык",
		Func: handlers2.AddSkill,
	},
	{
		Role: admin,
		Name: "[ Навыки ] Удалить навык",
		Func: handlers2.DeleteSkill,
	},
	{
		Role: admin,
		Name: "[ Навыки ] Редактировать навык",
		Func: handlers2.UpdateSkill,
	},
	{
		Role: user,
		Name: "[ Навыки ] Добавить навык",
		Func: handlers2.AddUserSkill,
	},
	{
		Role: user,
		Name: "[ Навыки ] Удалить навык",
		Func: handlers2.DeleteUserSkill,
	},
	{
		Role: user,
		Name: "[ Средства связи ] Добавить средство связи",
		Func: handlers2.AddContact,
	},
	{
		Role: user,
		Name: "[ Средства связи ] Редактировать средство связи",
		Func: handlers2.UpdateContact,
	},
	{
		Role: user,
		Name: "[ Средства связи ] Удалить средство связи",
		Func: handlers2.DeleteContact,
	},
	{
		Role: user,
		Name: "[ Финансовые показатели ] Добавить отчёт",
		Func: handlers2.AddReport,
	},
	{
		Role: user,
		Name: "[ Финансовые показатели ] Редактировать отчёт",
		Func: handlers2.UpdateFinReport,
	},
	{
		Role: user,
		Name: "[ Финансовые показатели ] Удалить отчёт",
		Func: handlers2.DeleteFinReport,
	},
	{
		Role: user,
		Name: "[ Финансовые показатели ] Сформировать полный отчёт",
		Func: handlers2.GetUserFinReport,
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
				t.app.Logger.Error("ошибка ввода логина")
				return fmt.Errorf("ошибка ввода логина")
			}

			fmt.Printf("Введите пароль: ")
			_, err = fmt.Scanf("%s", &password)
			if err != nil {
				t.app.Logger.Error("ошибка ввода пароля")
				return fmt.Errorf("ошибка ввода пароля")
			}

			ua := &domain.UserAuth{Username: login, Password: password}
			token, err := t.app.AuthSvc.Login(ua)
			if err != nil {
				t.app.Logger.Errorf("ошибка авторизации: %v", err)
				fmt.Printf("ошибка авторизации: %v\n", err)
				continue
			}

			payload, err := base.VerifyAuthToken(token, t.app.Config.Server.JwtKey)
			if err != nil {
				t.app.Logger.Errorf("ошибка верификации JWT токена: %v", err)
				return fmt.Errorf("ошибка верификации JWT токена: %w", err)
			}

			t.userInfo = payload
			err = t.userMenu()
			if err != nil {
				t.app.Logger.Error(err)
				fmt.Println(err)
				continue
			}
		}
	}
}

func (t *TUI) userMenu() (err error) {
	for {
		allowedActions, prompt := generateActionsPrompt(t.userInfo.Role)
		fmt.Println(prompt)

		var choice int
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			t.app.Logger.Error("ошибка ввода: %v", err)
			fmt.Printf("Ошибка ввода: %v", err)
			continue
		}

		if choice == 0 {
			return nil
		} else if choice > len(allowedActions) || choice < 0 {
			fmt.Println("Действия с таким номером нет.")
		} else {
			err = allowedActions[choice-1].Func(t.app, t.userInfo.Username)
			if err != nil {
				t.app.Logger.Error(err)
				fmt.Println(err)
			}
		}
	}
}

func (t *TUI) guestMenu() (err error) {
	var choice int
	for {
		fmt.Println(guestPrompt)
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			t.app.Logger.Error("ошибка ввода: %v", err)
			fmt.Printf("ошибка ввода: %v", err)
		}

		switch choice {
		case 1:
			var login, password string

			fmt.Printf("Введите логин: ")
			_, err = fmt.Scanf("%s", &login)
			if err != nil {
				t.app.Logger.Error("ошибка ввода логина")
				return fmt.Errorf("ошибка ввода логина")
			}

			fmt.Printf("Введите пароль: ")
			_, err = fmt.Scanf("%s", &password)
			if err != nil {
				t.app.Logger.Error("ошибка ввода пароля")
				return fmt.Errorf("ошибка ввода пароля")
			}

			ua := &domain.UserAuth{Username: login, Password: password}
			err = t.app.AuthSvc.Register(ua)
			if err != nil {
				t.app.Logger.Errorf("ошибка регистрации: %v", err)
				return fmt.Errorf("ошибка регистрации: %w", err)
			}
			return nil
		case 2:
			err = utils.PrintPaginatedCollection("Предприниматели", t.app.UserSvc.GetAll)
			if err != nil {
				t.app.Logger.Errorf("ошибка просмотра списка предпринимателей: %v", err)
				return fmt.Errorf("ошибка просмотра списка предпринимателей: %w", err)
			}
		case 0:
			return nil
		}
	}
}

func generateActionsPrompt(role string) (actionsList []Action, actionsPrompt string) {
	actionsList = make([]Action, 0)

	j := 1
	for _, action := range actions {
		if role == action.Role {
			actionsList = append(actionsList, action)
			actionsPrompt += fmt.Sprintf("\n%d. %s", j, action.Name)
			j++
		}
	}
	actionsPrompt += "\n\nВыберите действие: "

	return actionsList, actionsPrompt
}