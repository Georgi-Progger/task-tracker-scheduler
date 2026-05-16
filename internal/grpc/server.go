package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Georgi-Progger/task-tracker-common/configurator"
	"github.com/Georgi-Progger/task-tracker-common/kafka/producer"
	"github.com/Georgi-Progger/task-tracker-common/logger"
	"github.com/Georgi-Progger/task-tracker-scheduler/internal/cron"
	"github.com/Georgi-Progger/task-tracker-scheduler/pkg/pb/scheduler"
	"github.com/google/uuid"
)

type Server struct {
	scheduler.UnimplementedSchedulerSServiceServer
	logger logger.Logger
	cron   *cron.Cron
	cfg    configurator.Config
}

func New(cron *cron.Cron, logger logger.Logger) *Server {
	return &Server{
		cron:   cron,
		logger: logger,
	}
}

func (s *Server) CreateDailyEvent(ctx context.Context, req *scheduler.CreateDailyEventRequest) (*scheduler.CreateDailyEventResponse, error) {
	jobID := uuid.New().String()

	s.cron.AddDailyJob(&cron.Job{
		ID:     jobID,
		Hour:   int(req.Hour),
		Minute: int(req.Minute),
		Run: func() {
			producer := producer.NewProducer(
				s.cfg.GetBrokers(),
				s.cfg.GetEventsTopic(),
				logger.NewLogger(),
			)

			event := map[string]string{
				"event_type": "daily_report",
				"timestamp":  time.Now().Format(time.RFC3339),
			}

			eventJSON, _ := json.Marshal(event)

			err := producer.Send(context.Background(), eventJSON)
			if err != nil {
				s.logger.Error(err, "Failed to send event")
			} else {
				s.logger.Info(fmt.Sprintf("Daily report event sent at %s", time.Now().Format("15:04:05")))
			}
		},
	})

	return &scheduler.CreateDailyEventResponse{JobId: jobID}, nil
}
