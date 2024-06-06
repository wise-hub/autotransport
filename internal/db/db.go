package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func InitDB() *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"postgres")

	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to the database: %v", err))
	}

	var exists bool
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", "autotransport").Scan(&exists)
	if err != nil {
		panic(fmt.Sprintf("Unable to check if database exists: %v", err))
	}

	if !exists {
		_, err = db.Exec(context.Background(), "CREATE DATABASE autotransport")
		if err != nil {
			panic(fmt.Sprintf("Unable to create database: %v", err))
		}
	}

	db.Close()

	dsn = fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		"autotransport")

	db, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to the database: %v", err))
	}

	sqlFile, err := os.ReadFile("./ddl.sql")
	if err != nil {
		panic(fmt.Sprintf("Unable to read SQL file: %v", err))
	}

	_, err = db.Exec(context.Background(), string(sqlFile))
	if err != nil {
		panic(fmt.Sprintf("Failed to execute DDL statements: %v", err))
	}

	return db
}
