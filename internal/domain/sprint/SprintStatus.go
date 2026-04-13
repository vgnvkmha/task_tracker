package sprint

import (
	"errors"
)

type SprintStatus string

const (
	SprintDraft     SprintStatus = "draft"
	SprintPlanned   SprintStatus = "planned"
	SprintActive    SprintStatus = "active"
	SprintCompleted SprintStatus = "completed"
	SprintCancelled SprintStatus = "cancelled"
)

func (s SprintStatus) IsImmutable() error {
	switch s {
	case SprintCompleted, SprintCancelled:
		return errors.New("immutable sprint status")
	default:
		return nil
	}
}
