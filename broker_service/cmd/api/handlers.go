package main

import (
	"broker_service/events"
	"broker_service/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown call"))
		return
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	request, err := http.NewRequest("POST", "http://auth-service:1800/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("error while calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	if err := json.NewDecoder(response.Body).Decode(&jsonFromService); err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "authenticated"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, 200, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service:1900/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, errors.New("error from logger-service"+err.Error()))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error while calling logger-service"))
		return
	}

	var payload jsonResponse

	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, 200, payload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsondata, _ := json.MarshalIndent(msg, "", "\t")

	mailerURL := "http://mailer-service:2000/send"

	request, err := http.NewRequest("POST", mailerURL, bytes.NewBuffer(jsondata))
	if err != nil {
		app.errorJSON(w, errors.New("error while creating request"+err.Error()))
		return
	}

	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, errors.New("error while sending mail "+err.Error()))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error from mailer service"))
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "message sent",
	}

	app.writeJSON(w, 200, payload)
}

func (app *Config) logEventViaRabbit(w http.ResponseWriter, payload LogPayload) {
	if err := app.pushToQueue(payload.Name, payload.Data); err != nil {
		app.errorJSON(w, err)
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "logged via rabbitmq"

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := events.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")
	if err = emitter.Push(string(j), "log.INFO"); err != nil {
		return err
	}

	return nil
}

func (app *Config) logViaRPC(w http.ResponseWriter, logPayload LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := struct {
		Name string
		Data string
	}{
		Name: logPayload.Name,
		Data: logPayload.Data,
	}

	var result string
	if err = client.Call("RPCServer.LogInfo", payload, &result); err != nil {
		app.errorJSON(w, err)
		return
	}

	var resp jsonResponse
	resp.Error = false
	resp.Message = "logged via RPC"

	app.writeJSON(w, 200, resp)
}

func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var reqPayload RequestPayload

	if err := app.readJSON(w, r, &reqPayload); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	client := logs.NewLogServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: reqPayload.Log.Name,
			Data: reqPayload.Log.Data,
		},
	})

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged via gRPC"

	app.writeJSON(w, 200, payload)
}
