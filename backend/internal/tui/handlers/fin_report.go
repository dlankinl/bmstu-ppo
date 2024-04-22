package handlers

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"os"
	"ppo/domain"
	"ppo/internal/app"
	"ppo/internal/tui/utils"
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

	_, err = fmt.Scanf("%f", &revenue)
	if err != nil {
		return fmt.Errorf("ошибка ввода выручки: %w", err)
	}

	_, err = fmt.Scanf("%f", &costs)
	if err != nil {
		return fmt.Errorf("ошибка ввода расходов: %w", err)
	}

	_, err = fmt.Scanf("%d", &year)
	if err != nil {
		return fmt.Errorf("ошибка ввода года: %w", err)
	}

	_, err = fmt.Scanf("%f", &quarter)
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

	err = a.FinSvc.DeleteById(compUuid)
	if err != nil {
		return fmt.Errorf("ошибка удаления финансового отчета: %w", err)
	}

	return nil
}

func GetCompanyReports(a *app.App, args ...any) (err error) {
	var compId string
	var ok bool
	if len(args) > 0 {
		compId, ok = args[0].(string)
		if !ok {
			return fmt.Errorf("приведение аргумента к string")
		}
	}
	compId = strings.TrimSpace(compId)

	compUuid, err := uuid.Parse(compId)
	if err != nil {
		return fmt.Errorf("парсинг uuid из строки: %w", err)
	}

	fmt.Println("You are:", user.Username, user.ID)
	err = utils.PrintPaginatedCollectionArgs("Отчёты", a.FinSvc.GetByCompany, compUuid)
	if err != nil {
		return fmt.Errorf("вывод навыков с пагинацией: %w", err)
	}

	return nil
}
