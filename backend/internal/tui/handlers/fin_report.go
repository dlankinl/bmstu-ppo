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
	"strconv"
	"strings"
)

func AddReport(a *app.App, args ...any) (err error) {
	ctx := context.Background()
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

	err = a.FinSvc.Create(ctx, &report)
	if err != nil {
		return fmt.Errorf("ошибка добавления финансового отчета: %w", err)
	}

	return nil
}

func DeleteFinReport(a *app.App, args ...any) (err error) {
	ctx := context.Background()
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

	err = a.FinSvc.DeleteById(ctx, repUuid)
	if err != nil {
		return fmt.Errorf("ошибка удаления финансового отчета: %w", err)
	}

	return nil
}

func UpdateFinReport(a *app.App, args ...any) (err error) {
	ctx := context.Background()
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

	rep, err := a.FinSvc.GetById(ctx, repUuid)
	if err != nil {
		return fmt.Errorf("получение отчета по id: %w", err)
	}

	var revenue, costs float64
	var yearUpd, quarter int
	var revenueStr, costsStr, yearStr, quarterStr string

	fmt.Printf("Введите выручку (%f): ", rep.Revenue)
	revenueStr, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода выручки: %w", err)
	}
	revenueStr = strings.TrimSpace(revenueStr)
	if revenueStr != "" {
		revenue, err = strconv.ParseFloat(revenueStr, 32)
		if err != nil {
			return fmt.Errorf("парсинг строки в revenue: %w", err)
		}

		rep.Revenue = float32(revenue)
	}

	fmt.Printf("Введите расходы (%f): ", rep.Costs)
	costsStr, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода расходов: %w", err)
	}
	costsStr = strings.TrimSpace(costsStr)
	if costsStr != "" {
		costs, err = strconv.ParseFloat(costsStr, 32)
		if err != nil {
			return fmt.Errorf("парсинг строки в cost: %w", err)
		}

		rep.Costs = float32(costs)
	}

	fmt.Printf("Введите год (%d): ", rep.Year)
	yearStr, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода года: %w", err)
	}
	yearStr = strings.TrimSpace(yearStr)
	if yearStr != "" {
		yearUpd, err = strconv.Atoi(yearStr)
		if err != nil {
			return fmt.Errorf("парсинг строки в year: %w", err)
		}

		rep.Year = yearUpd
	}

	fmt.Printf("Введите квартал (%d): ", rep.Quarter)
	quarterStr, err = reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("ошибка ввода квартала: %w", err)
	}
	quarterStr = strings.TrimSpace(quarterStr)
	if quarterStr != "" {
		quarter, err = strconv.Atoi(quarterStr)
		if err != nil {
			return fmt.Errorf("парсинг строки в quarter: %w", err)
		}

		rep.Quarter = quarter
	}

	err = a.FinSvc.Update(ctx, rep)
	if err != nil {
		return fmt.Errorf("ошибка обновления финансового отчета: %w", err)
	}

	return nil
}

func GetCompanyReports(a *app.App, args ...any) (err error) {
	ctx := context.Background()
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

	err = utils.PrintYearCollection("Отчёты", a.FinSvc.GetByCompany, ctx, compUuid, year)
	if err != nil {
		return fmt.Errorf("вывод отчетов с пагинацией: %w", err)
	}

	return nil
}

func GetUserFinReport(a *app.App, args ...any) (err error) {
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

	fmt.Printf("Укажите квартал конца периода (1-4): ")
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

	rep, err := a.Interactor.GetUserFinancialReport(ctx, user.ID, period)
	if err != nil {
		return fmt.Errorf("формирование отчёта предпринимателя: %w", err)
	}

	fmt.Printf("Прибыль=%f, Выручка=%f, Затраты=%f", rep.Profit(), rep.Revenue(), rep.Costs())

	return nil
}
