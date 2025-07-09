package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	DBConn *pgxpool.Pool
}

func NewDBConfig() *DBConfig {
	return &DBConfig{getConn()}
}

func getConn() *pgxpool.Pool {
	ctx := context.Background()

	connStr := "postgresql://root:root@localhost:5432/rinha?sslmode=disable"

	db, err := pgxpool.New(ctx, connStr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	db.Config().MaxConnIdleTime = 10 * time.Minute
	db.Config().MaxConnLifetime = 2 * time.Hour
	db.Config().MaxConns = 200
	db.Config().MinConns = 30
	db.Config().HealthCheckPeriod = 10 * time.Minute

	if err := db.Ping(context.Background()); err != nil {
		panic("fail to connect in database")
	}

	return db
}
