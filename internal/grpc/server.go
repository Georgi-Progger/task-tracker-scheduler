package grpc

import (
	"context"

	"github.com/Georgi-Progger/task-tracker-common/kafka/producer"
	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/Georgi-Progger/task-tracker-scheduler/internal/cron"
	"github.com/Georgi-Progger/task-tracker-scheduler/pkg/pb/scheduler"
	"github.com/google/uuid"
)

type Server struct {
	scheduler.UnimplementedSchedulerSServiceServer
	cron *cron.Cron
}

func New(cron *cron.Cron) *Server {
	return &Server{cron: cron}
}

func (s *Server) CreateDailyReportJob(ctx context.Context, req *scheduler.CreateJobRequest) (*scheduler.CreateJobResponse, error) {

	jobID := uuid.New().String()
	producer := producer.NewProducer(
		"kafka:29092",
		"EVENTS_NOTIFICATIONS",
		logger.NewLogger(),
	)

	s.cron.AddDailyJob(&cron.Job{
		ID:     jobID,
		Hour:   int(req.Hour),
		Minute: int(req.Minute),
		Run: func() {
			_ = producer
		},
	})

	return &scheduler.CreateJobResponse{JobId: jobID}, nil
}
