package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/Georgi-Progger/task-tracker-common/configurator"
	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/Georgi-Progger/task-tracker-scheduler/internal/cron"
	pb "github.com/Georgi-Progger/task-tracker-scheduler/internal/grpc"
	"github.com/Georgi-Progger/task-tracker-scheduler/pkg/pb/scheduler"
)

func main() {
	logger := logger.NewLogger()

	c := cron.New(logger)
	c.Start()

	cfg, err := configurator.LoadConfig()
	if err != nil {
		logger.Error(err, "Failed to load config")
		os.Exit(1)
	}

	server := pb.New(c, logger)

	port := cfg.GetSchedulerPort()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Error(err, "Error lister server")
	}

	grpcServer := grpc.NewServer()
	scheduler.RegisterSchedulerSServiceServer(grpcServer, server)

	logger.Info(fmt.Sprintf("scheduler started on :%s", port))
	grpcServer.Serve(lis)
}
