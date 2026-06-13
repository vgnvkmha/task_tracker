package user

import (
	userApplication "task_tracker/internal/application/user"
	domainUser "task_tracker/internal/domain/user"
	valueobjects "task_tracker/internal/domain/value_objects"

	"github.com/google/uuid"
)

type Response struct {
	ID             uuid.UUID          `json:"id"`
	Email          valueobjects.Email `json:"email"`
	Role           valueobjects.Role  `json:"role"`
	TeamID         *uuid.UUID         `json:"team_id"`
	PersonalDataID uuid.UUID          `json:"personal_data_id"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	FullName       string             `json:"full_name"`
}

func FromDomain(user *domainUser.User) Response {
	return Response{
		ID:             user.ID,
		Email:          user.Email,
		Role:           user.Role,
		TeamID:         user.TeamID,
		PersonalDataID: user.PersonalDataID,
	}
}

func FromProfile(profile *userApplication.Profile) Response {
	response := FromDomain(profile.User)
	response.FirstName = profile.PersonalData.FirstName
	response.LastName = profile.PersonalData.LastName
	response.FullName = profile.PersonalData.FirstName + " " + profile.PersonalData.LastName
	return response
}

func FromDomainReponses(users []*domainUser.User) []Response {
	var response []Response
	for _, v := range users {
		response = append(response, FromDomain(v))
	}
	return response
}
