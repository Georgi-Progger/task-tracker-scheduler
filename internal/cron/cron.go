package cron

import (
	"log"
	"sync"
	"time"
)

type Cron struct {
	mu   sync.Mutex
	jobs map[string]*Job
}

func New() *Cron {
	return &Cron{
		jobs: make(map[string]*Job),
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
	c.mu.Lock()
	defer c.mu.Unlock()
	c.jobs[job.ID] = job
}

func (c *Cron) tick() {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, job := range c.jobs {
		if job.Hour == now.Hour() &&
			job.Minute == now.Minute() &&
			!sameDay(job.LastRun, now) {

			log.Println("cron: executing job", job.ID)
			job.LastRun = now
			go job.Run()
		}
	}
}

func sameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
