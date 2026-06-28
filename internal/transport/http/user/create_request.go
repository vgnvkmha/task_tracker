package user

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	userApplication "task_tracker/internal/application/user"

	"github.com/google/uuid"
)

var (
	errInvalidForm           = errors.New("invalid form")
	errInvalidAge            = errors.New("invalid age")
	errInvalidBirthDate      = errors.New("invalid birth date")
	errRequiredFieldsMissing = errors.New("required fields are missing")
	errLoginFieldsMissing    = errors.New("login fields are missing")
	errPasswordMismatch      = errors.New("passwords do not match")
)

type CreateRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Age       *uint8 `json:"age"`
}

func (r CreateRequest) ToServiceInput() userApplication.CreateUserInput {
	return userApplication.CreateUserInput{
		Email:     r.Email,
		Password:  r.Password,
		Role:      r.Role,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Age:       r.Age,
	}
}

type CreateFormRequest struct {
	Email                string
	Password             string
	PasswordConfirmation string
	Role                 string
	FirstName            string
	LastName             string
	Age                  *uint8
}

func NewCreateFormRequest(r *http.Request) (CreateFormRequest, error) {
	if err := r.ParseForm(); err != nil {
		return CreateFormRequest{}, errInvalidForm
	}

	request := CreateFormRequest{
		Email:                strings.TrimSpace(r.PostFormValue("email")),
		Password:             r.PostFormValue("password"),
		PasswordConfirmation: r.PostFormValue("password_confirmation"),
		Role:                 strings.TrimSpace(r.PostFormValue("role")),
		FirstName:            strings.TrimSpace(r.PostFormValue("first_name")),
		LastName:             strings.TrimSpace(r.PostFormValue("last_name")),
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

	if request.Email == "" || request.Password == "" || request.PasswordConfirmation == "" || request.Role == "" ||
		request.FirstName == "" || request.LastName == "" {
		return CreateFormRequest{}, errRequiredFieldsMissing
	}

	if request.Password != request.PasswordConfirmation {
		return CreateFormRequest{}, errPasswordMismatch
	}

	return request, nil
}

func (r CreateFormRequest) ToServiceInput() userApplication.CreateUserInput {
	return userApplication.CreateUserInput{
		Email:     r.Email,
		Password:  r.Password,
		Role:      r.Role,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Age:       r.Age,
	}
}

type LoginFormRequest struct {
	Email    string
	Password string
}

func NewLoginFormRequest(r *http.Request) (LoginFormRequest, error) {
	if err := r.ParseForm(); err != nil {
		return LoginFormRequest{}, errInvalidForm
	}

	request := LoginFormRequest{
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: r.FormValue("password"),
	}

	if request.Email == "" || request.Password == "" {
		return LoginFormRequest{}, errLoginFieldsMissing
	}

	return request, nil
}

type UpdateProfileFormRequest struct {
	UserID    uuid.UUID
	TeamID    *uuid.UUID
	TeamName  *string
	Email     *string
	Password  *string
	FirstName *string
	LastName  *string
	Age       *uint8
	BirthDate *time.Time
}

func NewUpdateProfileFormRequest(r *http.Request) (UpdateProfileFormRequest, error) {
	if err := r.ParseForm(); err != nil {
		return UpdateProfileFormRequest{}, errInvalidForm
	}

	userID, err := uuid.Parse(strings.TrimSpace(r.PostFormValue("user_id")))
	if err != nil {
		return UpdateProfileFormRequest{}, userApplication.ErrInvalidUserID
	}

	request := UpdateProfileFormRequest{UserID: userID}
	if value := strings.TrimSpace(r.PostFormValue("email")); value != "" {
		request.Email = &value
	}
	if value := r.PostFormValue("password"); value != "" {
		request.Password = &value
	}
	if value := strings.TrimSpace(r.PostFormValue("team_id")); value != "" {
		teamID, err := uuid.Parse(value)
		if err != nil {
			return UpdateProfileFormRequest{}, userApplication.ErrInvalidTeamID
		}
		request.TeamID = &teamID
	}
	if value := strings.TrimSpace(r.PostFormValue("team_name")); value != "" {
		request.TeamName = &value
	}
	if value := strings.TrimSpace(r.PostFormValue("first_name")); value != "" {
		request.FirstName = &value
	}
	if value := strings.TrimSpace(r.PostFormValue("last_name")); value != "" {
		request.LastName = &value
	}

	ageValue := strings.TrimSpace(r.PostFormValue("age"))
	if ageValue != "" {
		age, err := strconv.ParseUint(ageValue, 10, 8)
		if err != nil {
			return UpdateProfileFormRequest{}, errInvalidAge
		}
		ageUint8 := uint8(age)
		request.Age = &ageUint8
	}

	birthDateValue := strings.TrimSpace(r.PostFormValue("birth_date"))
	if birthDateValue != "" {
		birthDate, err := time.Parse("2006-01-02", birthDateValue)
		if err != nil {
			return UpdateProfileFormRequest{}, errInvalidBirthDate
		}
		request.BirthDate = &birthDate
	}

	return request, nil
}

func (r UpdateProfileFormRequest) ToServiceInput() userApplication.UpdateUserInput {
	return userApplication.UpdateUserInput{
		UserID:    r.UserID,
		TeamId:    r.TeamID,
		TeamName:  r.TeamName,
		Email:     r.Email,
		Password:  r.Password,
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Age:       r.Age,
		BirthDate: r.BirthDate,
	}
}
