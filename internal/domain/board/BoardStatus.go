package board

type BoardStatus string

const (
	BoardActive   BoardStatus = "active"
	BoardArchived BoardStatus = "archived"
)

func (s BoardStatus) IsValid() bool {
	switch s {
	case BoardActive, BoardArchived:
		return true
	default:
		return false
	}
}
