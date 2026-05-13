package user

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	userApplication "task_tracker/internal/application/user"
)

var (
	errInvalidForm           = errors.New("invalid form")
	errInvalidAge            = errors.New("invalid age")
	errInvalidBirthDate      = errors.New("invalid birth date")
	errRequiredFieldsMissing = errors.New("required fields are missing")
)

type CreateRequest struct {
	Email     string     `json:"email" binding:"required"`
	Password  string     `json:"password" binding:"required"`
	Role      string     `json:"role" binding:"required"`
	TeamName  *string    `json:"team_name"`
	FirstName string     `json:"first_name" binding:"required"`
	LastName  string     `json:"last_name" binding:"required"`
	Age       *uint8     `json:"age"`
	BirthDate *time.Time `json:"birth_date"`
}

func (r CreateRequest) ToServiceInput() userApplication.CreateUserInput {
	return userApplication.CreateUserInput{
		Email:     r.Email,
		Password:  r.Password,
		Role:      r.Role,
		TeamName:  r.TeamName,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Age:       r.Age,
		BirthDate: r.BirthDate,
	}
}

type CreateFormRequest struct {
	Email     string
	Password  string
	Role      string
	TeamName  *string
	FirstName string
	LastName  string
	Age       *uint8
	BirthDate *time.Time
}

func NewCreateFormRequest(r *http.Request) (CreateFormRequest, error) {
	if err := r.ParseForm(); err != nil {
		return CreateFormRequest{}, errInvalidForm
	}

	request := CreateFormRequest{
		Email:     strings.TrimSpace(r.PostFormValue("email")),
		Password:  r.PostFormValue("password"),
		Role:      strings.TrimSpace(r.PostFormValue("role")),
		FirstName: strings.TrimSpace(r.PostFormValue("first_name")),
		LastName:  strings.TrimSpace(r.PostFormValue("last_name")),
	}

	teamName := strings.TrimSpace(r.PostFormValue("team_name"))
	if teamName != "" {
		request.TeamName = &teamName
	}

	ageValue := strings.TrimSpace(r.PostFormValue("age"))
	if ageValue != "" {
		age, err := strconv.ParseUint(ageValue, 10, 8)
		if err != nil {
			return CreateFormRequest{}, errInvalidAge
		}
		ageUint8 := uint8(age)
		request.Age = &ageUint8
	}

	birthDateValue := strings.TrimSpace(r.PostFormValue("birth_date"))
	if birthDateValue != "" {
		birthDate, err := time.Parse("2006-01-02", birthDateValue)
		if err != nil {
			return CreateFormRequest{}, errInvalidBirthDate
		}
		request.BirthDate = &birthDate
	}

	if request.Email == "" || request.Password == "" || request.Role == "" ||
		request.FirstName == "" || request.LastName == "" {
		return CreateFormRequest{}, errRequiredFieldsMissing
	}

	return request, nil
}

func (r CreateFormRequest) ToServiceInput() userApplication.CreateUserInput {
	return userApplication.CreateUserInput{
		Email:     r.Email,
		Password:  r.Password,
		Role:      r.Role,
		TeamName:  r.TeamName,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Age:       r.Age,
		BirthDate: r.BirthDate,
	}
}
