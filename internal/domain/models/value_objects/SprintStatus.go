package valueobjects

type SprintStatus string

const (
	SprintPlanned   SprintStatus = "planned"
	SprintActive    SprintStatus = "active"
	SprintCompleted SprintStatus = "completed"
)
