package main

import (
	"context"

	"github.com/google/uuid"
)

func RegisterUserService(ctx context.Context, name, email, password string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}

	err = RegisterUserRepo(ctx, tx, name, email, password)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)

	return nil
}

func LoginUserService(ctx context.Context, email string) (string, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return "", err
	}

	chiperPassword := LoginUserRepo(ctx, tx, email)
	tx.Commit(ctx)

	return chiperPassword, nil
}

func GetAllIdUsersService(ctx context.Context) ([]uuid.UUID, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return []uuid.UUID{}, err
	}

	userIds, err := GetAllIdUsersRepo(ctx, tx)
	if err != nil {
		tx.Rollback(ctx)
		return []uuid.UUID{}, err
	}

	return userIds, nil
}

func CreateTodoService(ctx context.Context, userId uuid.UUID, title string, description string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = CreateTodoRepo(ctx, tx, userId, title, description)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)

	return nil
}

func UpdateTodoService(ctx context.Context, id uuid.UUID, title string, description string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = UpdateTodoRepo(ctx, tx, id, title, description)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)

	return nil
}

func DeleteTodoService(ctx context.Context, id uuid.UUID) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = DeleteTodoRepo(ctx, tx, id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	tx.Commit(ctx)

	return nil
}

func GetTodosService(ctx context.Context, page int, limit int) ([]Todo, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return []Todo{}, nil
	}

	todos, err := GetTodosRepo(ctx, tx, page, limit)
	if err != nil {
		tx.Rollback(ctx)
		return []Todo{}, nil
	}

	tx.Commit(ctx)

	return todos, nil
}
