package utils

import (
	"fmt"
	"ppo/domain"
)

func PrintUser(user *domain.User) {
	var gender string
	if user.Gender == "m" {
		gender = "Мужской"
	} else if user.Gender == "w" {
		gender = "Женский"
	}

	fmt.Printf("%s | %s | %s | %s", user.FullName, gender, user.Birthday.String(), user.City)
}

func PrintUsers(users []*domain.User) {
	for _, user := range users {
		PrintUser(user)
	}
}
