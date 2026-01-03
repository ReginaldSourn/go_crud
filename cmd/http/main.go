package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/reginaldsourn/go-crud/config"
	"github.com/reginaldsourn/go-crud/internal/adapters/primary/http"
	dbadapter "github.com/reginaldsourn/go-crud/internal/adapters/secondary/db"
	"github.com/reginaldsourn/go-crud/internal/adapters/secondary/db/migrations"
	"github.com/reginaldsourn/go-crud/internal/core/ports"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found; using existing environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	var db *gorm.DB
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("db connect failed: %v", err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("db handle failed: %v", err)
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("db ping failed: %v", err)
		}
		if err := migrations.Run(db); err != nil {
			log.Fatalf("db migrate failed: %v", err)
		}
		defer sqlDB.Close()
	} else {
		log.Printf("DATABASE_URL not set; running without a database")
	}

	var usersStore ports.UserStore
	if db != nil {
		usersStore = dbadapter.NewGormUserStore(db)
	}

	router := http.NewRouter(http.RouterDependencies{
		UserStore: usersStore,
		JWTSecret: []byte(cfg.JWTSecret),
		JWTTTL:    cfg.JWTTTL,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := ":" + cfg.Port
	if err := http.Serve(ctx, addr, router); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
