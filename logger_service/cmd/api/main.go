package main

import (
	"context"
	"fmt"
	"log"
	"logger_service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "1900"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Fatalf("error while connecting to Mongo: %v", err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf(err.Error())
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	log.Println("starting service on port: ", webPort)

	// register the rpc server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Fatalf("error while registering RPC server: %v", err)
		return
	}

	// exposing server on rpc
	go app.RPCListen()

	// exposing server on gRPC
	go app.gRPCListen()

	// exposing server on REST
	app.serve()

}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error while listening on %s: %v", webPort, err)
	}
}

func (app *Config) RPCListen() error {
	log.Println("starting rpc server on port: ", rpcPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Println("error while listening: ", err)
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}
func connectToMongo() (*mongo.Client, error) {
	opts := options.Client().ApplyURI(mongoURL)
	opts.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Println("error connecting to Mongo: ", err)
		return nil, err
	}

	log.Println("mongo db connected")
	return c, nil
}