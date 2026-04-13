package valueobjects

type Role string

const (
	Admin   Role = "admin"
	Captain Role = "captain"
	User    Role = "user"
	Guest   Role = "guest"
)

func (r Role) IsManagerRole() bool {
	switch r {
	case Admin, Captain:
		return true
	default:
		return false
	}
}
