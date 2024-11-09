package main

import (
	"context"
	"fmt"
	"log"
	"logger_service/data"
	"logger_service/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write log in mongodb

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	if err := l.Models.LogEntry.Insert(logEntry); err != nil {
		res := &logs.LogResponse{
			Result: "failed to insert log entry",
		}

		return res, err
	}

	return &logs.LogResponse{
		Result: "log entry inserted successfully",
	}, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("error while listening on port %s : %v", grpcPort, err)
		return
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{
		Models: app.Models,
	})

	log.Printf("gRPC server listening on port %s\n", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("error while listening on port %s : %v", grpcPort, err)
		return
	}
}
