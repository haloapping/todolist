package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserRepo(ctx context.Context, tx pgx.Tx, name string, email string, password string) error {
	q := `
		INSERT INTO users
		VALUES($1, $2, $3, $4);
	`
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}
	hashPass := string(bytes)

	_, err = tx.Exec(ctx, q, uuid.New(), name, email, hashPass)
	if err != nil {
		return err
	}

	return nil
}

func LoginUserRepo(ctx context.Context, tx pgx.Tx, email string) string {
	q := `
		SELECT password
		FROM users
		WHERE email = $1;
	`

	row := tx.QueryRow(ctx, q, email)
	var chiperPassword string
	row.Scan(&chiperPassword)

	return chiperPassword
}

func GetAllIdUsersRepo(ctx context.Context, tx pgx.Tx) ([]uuid.UUID, error) {
	q := `
		SELECT id FROM users;
	`

	rows, err := tx.Query(ctx, q)
	if err != nil {
		return []uuid.UUID{}, err
	}

	var id uuid.UUID
	var ids []uuid.UUID
	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids, nil
}

func CreateTodoRepo(ctx context.Context, tx pgx.Tx, userId uuid.UUID, title string, description string) error {
	q := `
		INSERT INTO todos(id, user_id, title, description)
		VALUES($1, $2, $3, $4)
	`
	_, err := tx.Exec(ctx, q, uuid.New(), userId, title, description)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTodoRepo(ctx context.Context, tx pgx.Tx, id uuid.UUID, title string, description string) error {
	q := `
		UPDATE todos
		SET title = $1,
			description = $2
		WHERE id = $3;
	`

	_, err := tx.Exec(ctx, q, title, description, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTodoRepo(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	q := `
		DELETE FROM todos
		WHERE id = $1;
	`

	_, err := tx.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func GetTodosRepo(ctx context.Context, tx pgx.Tx, page int, limit int) ([]Todo, error) {
	q := `
		SELECT * FROM todos LIMIT $1 OFFSET $2;
	`

	rows, err := tx.Query(ctx, q, limit, page)
	if err != nil {
		return []Todo{}, err
	}

	var todos []Todo
	for rows.Next() {
		var t Todo
		rows.Scan(&t.Id, &t.UserId, &t.Title, &t.Description)

		todos = append(todos, t)
	}

	return todos, nil
}
