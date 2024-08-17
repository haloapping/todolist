package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	host = "localhost"
	port = "8000"

	dbDriver   = "postgres"
	dbUserName = "postgres"
	dbPassword = "appingganteng"
	dbHost     = "localhost"
	dbPort     = "5432"
	dbName     = "todo"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf("%s://%s:%s@%s:%s/%s", dbDriver, dbUserName, dbPassword, dbHost, dbPort, dbName),
	)
	if err != nil {
		panic(err)
	}
	SetPool(pool)

	r := Route()

	server := http.Server{
		Addr:    host + ":" + port,
		Handler: Logger(r),
	}

	slog.Info(fmt.Sprintf("Server started at %s:%s...", host, port))
	// GenerateUserFakeData(50)
	// GenerateTodoFakeData(100)
	server.ListenAndServe()

	slog.Info("Server stopped")
}
