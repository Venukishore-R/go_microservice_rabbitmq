package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Mailer Mail
}

const webPort = "2000"

func main() {
	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail service on port: ", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error while serving mail-service: %v", err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	mail := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  strings.ToLower(os.Getenv("MAIL_ENCRYPTION")),
		FromName:    os.Getenv("MAIL_FROMNAME"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
	}

	return mail
}
