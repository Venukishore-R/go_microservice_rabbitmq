package main

import (
	"errors"
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage
	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, errors.New("invalid payload"+err.Error()), http.StatusBadRequest)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	if err := app.Mailer.SendSMTPMessage(msg); err != nil {
		log.Println("error", err)
		app.errorJSON(w, errors.New("error sending message"+err.Error()), http.StatusInternalServerError)
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "message sent successfully" + requestPayload.To

	app.writeJSON(w, http.StatusAccepted, resp)
}
