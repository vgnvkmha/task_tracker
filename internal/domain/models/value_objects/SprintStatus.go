package valueobjects

type SprintStatus string

const (
	SprintDraft     SprintStatus = "draft"
	SprintPlanned   SprintStatus = "planned"
	SprintActive    SprintStatus = "active"
	SprintCompleted SprintStatus = "completed"
	SprintCancelled SprintStatus = "cancelled"
)

func (s SprintStatus) IsImmutable() bool {
	switch s {
	case SprintCompleted, SprintCancelled:
		return true
	default:
		return false
	}
}
