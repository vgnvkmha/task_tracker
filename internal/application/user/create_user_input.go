package user

type CreateUserInput struct {
	Email     string
	Password  string
	Role      string
	FirstName string
	LastName  string
	Age       *uint8
}
