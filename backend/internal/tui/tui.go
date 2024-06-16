package tui

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/config"
	"ppo/internal/tui/handlers"
	"ppo/internal/tui/utils"
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
		Func: handlers.UpdateUser,
	},
	{
		Role: admin,
		Name: "[ Предприниматели ] Создать карточку предпринимателя",
		Func: handlers.CreateUser,
	},
	{
		Role: user,
		Name: "[ Компании ] Добавить компанию",
		Func: handlers.AddCompany,
	},
	{
		Role: user,
		Name: "[ Компании ] Удалить компанию",
		Func: handlers.DeleteCompany,
	},
	{
		Role: user,
		Name: "[ Компании ] Обновить информацию о компании",
		Func: handlers.UpdateCompany,
	},
	{
		Role: user,
		Name: "[ Компании ] Посмотреть список своих компаний",
		Func: handlers.GetMyCompanies,
	},
	{
		Role: user,
		Name: "[ Предприниматели ] Просчитать рейтинг",
		Func: handlers.CalculateRating,
	},
	{
		Role: admin,
		Name: "[ Предприниматели ] Сменить роль пользователя",
		Func: handlers.ChangeUserRole,
	},
	{
		Role: user,
		Name: "[ Предприниматели ] Просмотреть список предпринимателей",
		Func: handlers.GetAllUsers,
	},
	{
		Role: admin,
		Name: "[ Сферы деятельности ] Добавить сферу деятельности",
		Func: handlers.AddActivityField,
	},
	{
		Role: admin,
		Name: "[ Сферы деятельности ] Удалить сферу деятельности",
		Func: handlers.DeleteActivityField,
	},
	{
		Role: admin,
		Name: "[ Сферы деятельности ] Редактировать сферу деятельности",
		Func: handlers.UpdateActivityField,
	},
	{
		Role: admin,
		Name: "[ Навыки ] Добавить навык",
		Func: handlers.AddSkill,
	},
	{
		Role: admin,
		Name: "[ Навыки ] Удалить навык",
		Func: handlers.DeleteSkill,
	},
	{
		Role: admin,
		Name: "[ Навыки ] Редактировать навык",
		Func: handlers.UpdateSkill,
	},
	{
		Role: user,
		Name: "[ Навыки ] Добавить навык",
		Func: handlers.AddUserSkill,
	},
	{
		Role: user,
		Name: "[ Навыки ] Удалить навык",
		Func: handlers.DeleteUserSkill,
	},
	{
		Role: user,
		Name: "[ Средства связи ] Добавить средство связи",
		Func: handlers.AddContact,
	},
	{
		Role: user,
		Name: "[ Средства связи ] Редактировать средство связи",
		Func: handlers.UpdateContact,
	},
	{
		Role: user,
		Name: "[ Средства связи ] Удалить средство связи",
		Func: handlers.DeleteContact,
	},
	{
		Role: user,
		Name: "[ Финансовые показатели ] Добавить отчёт",
		Func: handlers.AddReport,
	},
	{
		Role: user,
		Name: "[ Финансовые показатели ] Редактировать отчёт",
		Func: handlers.UpdateFinReport,
	},
	{
		Role: user,
		Name: "[ Финансовые показатели ] Удалить отчёт",
		Func: handlers.DeleteFinReport,
	},
	{
		Role: user,
		Name: "[ Финансовые показатели ] Сформировать полный отчёт",
		Func: handlers.GetUserFinReport,
	},
}

func (t *TUI) Run() (err error) {
	ctx := context.Background()

	var choice int
	reader := bufio.NewReader(os.Stdin)
	_ = reader
	for {
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
			token, err := t.app.AuthSvc.Login(ctx, ua)
			if err != nil {
				fmt.Printf("ошибка авторизации: %v\n", err)
				continue
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
		allowedActions, prompt := generateActionsPrompt(t.userInfo.Role)
		fmt.Println(prompt)

		var choice int
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
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
				fmt.Println(err)
			}
		}
	}
}

func (t *TUI) guestMenu() (err error) {
	ctx := context.Background()

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
			err = t.app.AuthSvc.Register(ctx, ua)
			if err != nil {
				return fmt.Errorf("ошибка регистрации: %w", err)
			}
			return nil
		case 2:
			err = printPaginatedUsers(ctx, t.app)
			if err != nil {
				return fmt.Errorf("ошибка просмотра списка предпринимателей: %w", err)
			}
		case 0:
			return nil
		}
	}
}

func printPaginatedUsers(ctx context.Context, app *app.App) (err error) {
	page := 1

	for {
		tmp, err := app.UserSvc.GetAll(ctx, page)
		if err != nil {
			return fmt.Errorf("получение пагинированных данных: %w", err)
		}

		utils.PrintCollection("Предприниматели", tmp)

		fmt.Printf("1. Предыдущая страница.\n2. Следующая страница.\n0. Назад.\n\nВыберите действие: ")
		var option int
		_, err = fmt.Scanf("%d", &option)
		if err != nil {
			return fmt.Errorf("ошибка ввода следующего действия: %w", err)
		}

		switch option {
		case 1:
			if page > 1 {
				page--
			}
		case 2:
			if len(tmp) == config.PageSize {
				page++
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
