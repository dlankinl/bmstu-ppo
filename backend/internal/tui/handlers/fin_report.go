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

	//compUuid, err := uuid.Parse(compId)
	//if err != nil {
	//	return fmt.Errorf("парсинг uuid из строки: %w", err)
	//}

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

	//compUuid, err := uuid.Parse(compId)
	//if err != nil {
	//	return fmt.Errorf("парсинг uuid из строки: %w", err)
	//}

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

// TODO: incorrect login or password
