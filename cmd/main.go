package main

import (
	"net"

	"google.golang.org/grpc"

	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/Georgi-Progger/task-tracker-scheduler/internal/cron"
	pb "github.com/Georgi-Progger/task-tracker-scheduler/internal/grpc"
	"github.com/Georgi-Progger/task-tracker-scheduler/pkg/pb/scheduler"
)

func main() {
	logger := logger.NewLogger()

	c := cron.New(logger)
	c.Start()

	server := pb.New(c, logger)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Error(err, "Error lister server")
	}

	grpcServer := grpc.NewServer()
	scheduler.RegisterSchedulerSServiceServer(grpcServer, server)

	logger.Info("scheduler started on :8082")
	grpcServer.Serve(lis)
}
