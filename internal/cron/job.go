package cron

import "time"

type Job struct {
	ID      string
	Hour    int
	Minute  int
	LastRun time.Time
	Run     func()
}
