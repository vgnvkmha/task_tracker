package models

import "time"

type Board struct {
	id         int
	team       Team
	isPublic   bool
	name       string
	created_at time.Time
}
