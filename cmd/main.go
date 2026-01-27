package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/Georgi-Progger/task-tracker-sheduler/internal/cron"
	pb "github.com/Georgi-Progger/task-tracker-sheduler/internal/grpc"
	"github.com/Georgi-Progger/task-tracker-sheduler/pkg/pb/scheduler"
)

func main() {
	c := cron.New()
	c.Start()

	server := pb.New(c)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	scheduler.RegisterSchedulerSServiceServer(grpcServer, server)

	log.Println("scheduler started on :8082")
	grpcServer.Serve(lis)
}
