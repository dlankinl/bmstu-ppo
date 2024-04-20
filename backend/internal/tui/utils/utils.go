package utils

import (
	"fmt"
	"github.com/google/uuid"
	"ppo/domain"
	"ppo/internal/config"
	"reflect"
)

func PrintHeader(val any) {
	var str string
	t := reflect.TypeOf(val)

	for i := 0; i < t.NumField(); i++ {
		str = fmt.Sprintf("%s | %s", str, t.Field(i).Name)
	}
	fmt.Println(str)
}

func PrintStruct(val any) {
	var str string
	v := reflect.ValueOf(val)

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch v := f.Interface().(type) {
		case int:
			str = fmt.Sprintf("%s | %d", str, v)
		case float32:
			str = fmt.Sprintf("%s | %f", str, v)
		case string, uuid.UUID:
			str = fmt.Sprintf("%s | %s", str, v)
		}
	}

	fmt.Printf("%s\n", str)
}

func PrintCollection[T any](collName string, collection []*T) {
	fmt.Println(collName)

	for i, val := range collection {
		if i == 0 {
			PrintHeader(*val)
		}

		if reflect.TypeOf(val).Kind() == reflect.Ptr {
			PrintStruct(reflect.ValueOf(val).Elem().Interface())
		} else {
			PrintStruct(val)
		}
	}
}

func PrintActivityField(field *domain.ActivityField) {
	fmt.Printf("%s | %s | %s | %f\n", field.ID, field.Name, field.Description, field.Cost)
}

func PrintActivityFields(fields []*domain.ActivityField) {
	fmt.Println("Сферы деятельности:")
	for _, field := range fields {
		PrintActivityField(field)
	}
}

func PrintPaginatedCollection[T any](collectionName string, fn func(int) ([]*T, error)) (err error) {
	page := 1
	for {
		tmp, err := fn(page)
		if err != nil {
			return fmt.Errorf("получение пагинированных данных: %w", err)
		}

		PrintCollection(collectionName, tmp)

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
