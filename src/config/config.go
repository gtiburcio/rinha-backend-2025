package config

import (
	"context"
	"errors"
	"fmt"
	"log"
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

	db, err := pgxpool.New(ctx, getStrConnection())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	db.Config().MaxConnIdleTime = 10 * time.Minute
	db.Config().MaxConnLifetime = 2 * time.Hour
	db.Config().MaxConns = 200
	db.Config().MinConns = 30
	db.Config().HealthCheckPeriod = 10 * time.Minute

	if err := pingDB(db); err != nil {
		log.Fatalf("fail to connect in database: %v", err)
	}

	return db
}

func pingDB(db *pgxpool.Pool) error {
	for i := 0; i < 5; i++ {
		if err := db.Ping(context.Background()); err != nil {
			time.Sleep(time.Second * 2)
		} else {
			return nil
		}
	}

	return errors.New("fail to connect on DB")
}

func getStrConnection() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
}
