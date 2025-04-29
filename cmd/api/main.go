package main

import (
	"github.com/gabriel-assis7/gophersocial/internal/db"
	"github.com/gabriel-assis7/gophersocial/internal/env"
	"go.uber.org/zap"
)

func main() {
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:9009@localhost:5432/gopher-social?sslmode=disable"),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		},
		env: env.GetString("ENV", "development"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxIdleTime,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
	)
	if err != nil {
		logger.Fatalf("Could not connect with the database: %v", err)
	}

	defer db.Close()
	logger.Info("Successfully connected with database")

	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
