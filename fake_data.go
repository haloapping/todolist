package main

import (
	"context"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateUserFakeData(nData int) error {
	for i := 0; i < nData; i++ {
		var u User
		u.Id = uuid.New()
		hashPass, err := bcrypt.GenerateFromPassword([]byte(gofakeit.Word()), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashPass)
		err = gofakeit.Struct(&u)
		if err != nil {
			return err
		}
		RegisterUserService(context.Background(), u.Name, u.Email, u.Password)
	}

	return nil
}

func GenerateTodoFakeData(nData int) error {
	userIds, err := GetAllIdUsersService(context.Background())
	if err != nil {
		return err
	}

	for i := 0; i < nData; i++ {
		userId := userIds[rand.Intn(len(userIds))]
		title := gofakeit.Word()
		description := gofakeit.Noun()
		if err != nil {
			return err
		}
		err = CreateTodoService(context.Background(), userId, title, description)
		if err != nil {
			return err
		}
	}

	return nil
}
