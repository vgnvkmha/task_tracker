package valueobjects

type BoardStatus string

const (
	BoardActive   BoardStatus = "active"
	BoardArchived BoardStatus = "archived"
)
