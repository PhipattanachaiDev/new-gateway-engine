package database

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	dbEzvPool      *pgxpool.Pool
	dbRealtimePool *pgxpool.Pool
	dbHistoryPool  *pgxpool.Pool

	onceEzv      sync.Once
	onceRealtime sync.Once
	onceHistory  sync.Once
)

func ConnDBEzViewLiteConfig() (*pgxpool.Pool, error) {
	var err error
	onceEzv.Do(func() {
		dbEzvPool, err = createPgxPool("DB_EZVIEWLITE")
	})
	return dbEzvPool, err
}

func ConnDBRealtimeConfig() (*pgxpool.Pool, error) {
	var err error
	onceRealtime.Do(func() {
		dbRealtimePool, err = createPgxPool("DB_REALTIME")
	})
	return dbRealtimePool, err
}

func ConnDBHistoryConfig() (*pgxpool.Pool, error) {
	var err error
	onceHistory.Do(func() {
		dbHistoryPool, err = createPgxPool("DB_HISTORY")
	})
	return dbHistoryPool, err
}

func createPgxPool(prefix string) (*pgxpool.Pool, error) {
	dbHost := os.Getenv(prefix + "_PG_HOST")
	dbUser := os.Getenv(prefix + "_PG_USER")
	dbPass := os.Getenv(prefix + "_PG_PASS")
	dbDatabase := os.Getenv(prefix + "_PG_DATABASE")
	dbPort := os.Getenv(prefix + "_PG_PORT")


	// fmt.Println("ENV values:", dbHost, dbUser, dbPass, dbDatabase, dbPort)


	if dbHost == "" || dbUser == "" || dbPass == "" || dbDatabase == "" || dbPort == "" {
		return nil, fmt.Errorf("missing %s database config environment variables", prefix)
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbDatabase,
	)

	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 5; i++ {
		config, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			fmt.Printf("failed to parse %s config: %v\n", prefix, err)
			return nil, err
		}

		config.MaxConns = 120
		config.MinConns = 10
		config.MaxConnLifetime = 5 * time.Minute
		config.MaxConnIdleTime = 1 * time.Minute

		pool, err = pgxpool.ConnectConfig(context.Background(), config)
		if err != nil {
			fmt.Printf("failed to connect to %s (attempt %d/5): %v\n", prefix, i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err = pool.Ping(context.Background()); err != nil {
			fmt.Printf("failed to ping %s (attempt %d/5): %v\n", prefix, i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to %s after multiple attempts: %v", prefix, err)
}
