package dto

import "time"

type Status int

const (
	Todo Status = iota
	InProgress
	Done
)

type Task struct {
	ID          int
	description string
	created_at  time.Time
	deadline    time.Time
	done        bool
	Type        Status
}
