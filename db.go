package main

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

func SetPool(newPool *pgxpool.Pool) error {
	if newPool == nil {
		return errors.New("cannot assign nil pool")
	}

	pool = newPool

	return nil
}
