package main

import (
	"auth_service/config"
	"auth_service/data"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	webPort = "1800"
)

var counts int64

type Auth struct {
	// DB     *gorm.DB
	// Models *data.Models
	Repo   data.Repository
	Client *http.Client
}

func openDB(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB(config *config.Config) (*gorm.DB, error) {
	for {
		db, err := openDB(config)
		if err != nil {
			log.Printf("error while connecting to database: %v", err)
			counts++
		} else {
			if err := db.AutoMigrate(&data.User{}); err != nil {
				log.Fatalf("error during database migration: %v", err)
			}

			slog.Info("connected to database", "success", "")
			return db, nil
		}

		if counts > 10 {
			log.Println("error while connecting to database", err)
			return nil, err
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}
func main() {
	config, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	db, err := connectToDB(config)
	if err != nil {
		log.Fatalf("error while connecting db %v", err)
	}

	app := Auth{
		Client: &http.Client{},
	}

	app.setupRepo(db)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	slog.Info("Server starting", "port", webPort)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error while serving auth-service: %v", err)
	}
}

func (app *Auth) setupRepo(conn *gorm.DB) {
	db := data.NewPostgresRepository(conn)
	app.Repo = db
}
