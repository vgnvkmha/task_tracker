package models

import "time"

type Board struct {
	id         uint32
	teamId     uint32
	isPublic   bool
	name       string
	created_at time.Time
}
