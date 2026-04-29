package team

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID       uuid.UUID
	Name     string
	Timezone *string
	LeaderID *uuid.UUID
	IsActive bool
}

func New(name string, timezone *string, leaderID *uuid.UUID) (*Team, error) {

	if err := validateName(&name); err != nil {
		return nil, err
	}

	if err := validateTimezone(timezone); err != nil {
		return nil, err
	}

	return &Team{
		ID:       uuid.New(),
		Name:     name,
		Timezone: timezone,
		LeaderID: leaderID,
		IsActive: true,
	}, nil
}
func (t *Team) ApplyChanges(name, tz *string, leader_id *uuid.UUID, is_active *bool) error {
	if err := t.setName(name); err != nil {
		return err
	}
	if err := t.setTimezone(tz); err != nil {
		return err
	}
	if err := t.setLeaderID(leader_id); err != nil {
		return err
	}
	if is_active != nil && *is_active == false {
		t.deactivate()
	}
	return nil
}

func (t *Team) setName(name *string) error {
	if err := validateName(name); err != nil {
		return err
	}
	t.Name = *name
	return nil
}

func (t *Team) setTimezone(tz *string) error {
	if err := validateTimezone(tz); err != nil {
		return err
	}
	t.Timezone = tz
	return nil
}

func (t *Team) setLeaderID(id *uuid.UUID) error {
	t.LeaderID = id
	return nil
}

func validateName(name *string) error {
	if name == nil || *name == "" {
		return ErrEmptyName
	}
	return nil
}

func (t *Team) deactivate() {
	t.IsActive = false
}

func validateTimezone(tz *string) error {
	if tz != nil && *tz != "" {
		if _, err := time.LoadLocation(*tz); err != nil {
			return ErrInvalidTZ
		}
	}
	return nil
}
