package handlers

import (
	"bufio"
	"business-logic/domain"
	"fmt"
	"github.com/google/uuid"
	"os"
	"ppo/app/internal/app"
	"ppo/app/internal/tui/utils"
	"strings"
)

func AddReport(a *app.App, args ...any) (err error) {
	reader := bufio.NewReader(os.Stdin)

	err = GetMyCompanies(a, args...)
	if err != nil {
		return fmt.Errorf("добавление финансового отчета: %w", err)
	}

	fmt.Printf("Введите id компании: ")
	compId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id компании: %w", err)
	}
	compId = strings.TrimSpace(compId)

	compUuid, err := uuid.Parse(compId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	var revenue, costs float32
	var year, quarter int

	fmt.Printf("Введите выручку: ")
	_, err = fmt.Scanf("%f", &revenue)
	if err != nil {
		return fmt.Errorf("ошибка ввода выручки: %w", err)
	}

	fmt.Printf("Введите расходы: ")
	_, err = fmt.Scanf("%f", &costs)
	if err != nil {
		return fmt.Errorf("ошибка ввода расходов: %w", err)
	}

	fmt.Printf("Введите год: ")
	_, err = fmt.Scanf("%d", &year)
	if err != nil {
		return fmt.Errorf("ошибка ввода года: %w", err)
	}

	fmt.Printf("Введите квартал (1-4): ")
	_, err = fmt.Scanf("%d", &quarter)
	if err != nil {
		return fmt.Errorf("ошибка ввода квартала: %w", err)
	}

	var report domain.FinancialReport
	report.CompanyID = compUuid
	report.Revenue = revenue
	report.Costs = costs
	report.Year = year
	report.Quarter = quarter

	err = a.FinSvc.Create(&report)
	if err != nil {
		return fmt.Errorf("ошибка добавления финансового отчета: %w", err)
	}

	return nil
}

func DeleteFinReport(a *app.App, args ...any) (err error) {
	reader := bufio.NewReader(os.Stdin)

	err = GetMyCompanies(a, args...)
	if err != nil {
		return fmt.Errorf("удаление финансового отчета: %w", err)
	}

	fmt.Printf("Введите id компании: ")
	compId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id компании: %w", err)
	}
	compId = strings.TrimSpace(compId)

	var year int
	fmt.Printf("Укажите год: ")
	_, err = fmt.Scanf("%d", &year)
	if err != nil {
		return fmt.Errorf("ошибка ввода года: %w", err)
	}

	err = GetCompanyReports(a, compId, year)
	if err != nil {
		return fmt.Errorf("обновление финансового отчета: %w", err)
	}

	fmt.Printf("Введите id отчета: ")
	repId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id отчета: %w", err)
	}
	repId = strings.TrimSpace(repId)

	repUuid, err := uuid.Parse(repId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	err = a.FinSvc.DeleteById(repUuid)
	if err != nil {
		return fmt.Errorf("ошибка удаления финансового отчета: %w", err)
	}

	return nil
}

func UpdateFinReport(a *app.App, args ...any) (err error) {
	reader := bufio.NewReader(os.Stdin)

	err = GetMyCompanies(a, args...)
	if err != nil {
		return fmt.Errorf("обновление финансового отчета: %w", err)
	}

	fmt.Printf("Введите id компании: ")
	compId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id компании: %w", err)
	}
	compId = strings.TrimSpace(compId)

	var year int
	fmt.Printf("Укажите год: ")
	_, err = fmt.Scanf("%d", &year)
	if err != nil {
		return fmt.Errorf("ошибка ввода года: %w", err)
	}

	err = GetCompanyReports(a, compId, year)
	if err != nil {
		return fmt.Errorf("обновление финансового отчета: %w", err)
	}

	fmt.Printf("Введите id отчета: ")
	repId, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода id отчета: %w", err)
	}
	repId = strings.TrimSpace(repId)

	repUuid, err := uuid.Parse(repId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	err = a.FinSvc.DeleteById(repUuid)
	if err != nil {
		return fmt.Errorf("ошибка удаления финансового отчета: %w", err)
	}

	return nil
}

func GetCompanyReports(a *app.App, args ...any) (err error) {
	var compId string
	var year int
	var ok bool
	if len(args) == 2 {
		compId, ok = args[0].(string)
		if !ok {
			return fmt.Errorf("приведение аргумента к string")
		}

		year, ok = args[1].(int)
		if !ok {
			return fmt.Errorf("приведение аргумента к int")
		}
	} else {
		return fmt.Errorf("некорректное число аргументов: %w", err)
	}
	compId = strings.TrimSpace(compId)

	compUuid, err := uuid.Parse(compId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	err = utils.PrintYearCollection("Отчёты", a.FinSvc.GetByCompany, compUuid, year)
	if err != nil {
		return fmt.Errorf("вывод отчетов с пагинацией: %w", err)
	}

	return nil
}

func GetUserFinReport(a *app.App, args ...any) (err error) {
	var username string
	var ok bool
	if len(args) > 0 {
		username, ok = args[0].(string)
		if !ok {
			return fmt.Errorf("приведение аргумента к string")
		}
	}

	user, err := a.UserSvc.GetByUsername(username)
	if err != nil {
		return fmt.Errorf("пользователь не найден")
	}

	var startYear, endYear, startQuarter, endQuarter int
	fmt.Printf("Укажите год начала периода: ")
	_, err = fmt.Scanf("%d", &startYear)
	if err != nil {
		return fmt.Errorf("ошибка ввода года: %w", err)
	}

	fmt.Printf("Укажите год конца периода: ")
	_, err = fmt.Scanf("%d", &endYear)
	if err != nil {
		return fmt.Errorf("ошибка ввода года: %w", err)
	}

	fmt.Printf("Укажите квартал начала периода (1-4): ")
	_, err = fmt.Scanf("%d", &startQuarter)
	if err != nil {
		return fmt.Errorf("ошибка ввода квартала: %w", err)
	}

	fmt.Printf("Укажите квартал конца периода (%d-4): ", startQuarter+1)
	_, err = fmt.Scanf("%d", &endQuarter)
	if err != nil {
		return fmt.Errorf("ошибка ввода квартала: %w", err)
	}

	period := &domain.Period{
		StartYear:    startYear,
		EndYear:      endYear,
		StartQuarter: startQuarter,
		EndQuarter:   endQuarter,
	}

	rep, err := a.Interactor.GetUserFinancialReport(user.ID, period)
	if err != nil {
		return fmt.Errorf("формирование отчёта предпринимателя: %w", err)
	}

	fmt.Printf("Прибыль=%f, Выручка=%f, Затраты=%f", rep.Profit(), rep.Revenue(), rep.Costs())

	return nil
}
