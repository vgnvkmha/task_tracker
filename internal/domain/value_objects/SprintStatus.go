package valueobjects

import err "task_tracker/internal/domain/errors"

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
		return err.ErrInvalidRights
	default:
		return nil
	}
}
