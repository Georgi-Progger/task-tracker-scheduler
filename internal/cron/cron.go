package cron

import (
	"sync"
	"time"

	"github.com/Georgi-Progger/task-tracker-common/logger"
)

type Cron struct {
	logger logger.Logger
	jobs   sync.Map // I use this for insurance purposes
}

func New(logger logger.Logger) *Cron {
	return &Cron{
		logger: logger,
	}
}

func (c *Cron) Start() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			c.tick()
		}
	}()
}

func (c *Cron) AddDailyJob(job *Job) {
	c.jobs.Store(job.ID, job)
}

func (c *Cron) tick() {
	now := time.Now()

	c.jobs.Range(func(key, value interface{}) bool {
		job := value.(*Job)

		if job.Hour == now.Hour() &&
			job.Minute == now.Minute() &&
			!(job.LastRun.Truncate(24 * time.Hour).Equal(now.Truncate(24 * time.Hour))) {

			job.LastRun = now
			go job.Run()
		}
		return true
	})
}
