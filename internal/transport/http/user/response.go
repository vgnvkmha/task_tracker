package user

import (
	"time"

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
	TeamName       string             `json:"team_name"`
	PersonalDataID uuid.UUID          `json:"personal_data_id"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	FullName       string             `json:"full_name"`
	Age            *uint8             `json:"age"`
	BirthDate      string             `json:"birth_date"`
	IsActive       bool               `json:"is_active"`
	IsDeleted      bool               `json:"is_deleted"`
	DeletedAt      string             `json:"deleted_at"`
}

func FromDomain(user *domainUser.User) Response {
	response := Response{
		ID:             user.ID,
		Email:          user.Email,
		Role:           user.Role,
		TeamID:         user.TeamID,
		PersonalDataID: user.PersonalDataID,
		IsActive:       user.IsActive,
		IsDeleted:      user.IsDeleted(),
	}
	if user.DeletedAt != nil {
		response.DeletedAt = user.DeletedAt.Format(time.RFC3339)
	}
	return response
}

func FromProfile(profile *userApplication.Profile) Response {
	response := FromDomain(profile.User)
	response.FirstName = profile.PersonalData.FirstName
	response.LastName = profile.PersonalData.LastName
	response.FullName = profile.PersonalData.FirstName + " " + profile.PersonalData.LastName
	response.TeamName = profile.TeamName
	response.Age = profile.PersonalData.Age
	if profile.PersonalData.BirthDate != nil {
		response.BirthDate = profile.PersonalData.BirthDate.Format(time.DateOnly)
	}
	return response
}

func FromProfiles(profiles []*userApplication.Profile) []Response {
	response := make([]Response, 0, len(profiles))
	for _, profile := range profiles {
		if profile != nil {
			response = append(response, FromProfile(profile))
		}
	}
	return response
}

func FromDomainReponses(users []*domainUser.User) []Response {
	var response []Response
	for _, v := range users {
		response = append(response, FromDomain(v))
	}
	return response
}
