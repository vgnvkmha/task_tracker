package personaldata

import "time"

func (data *PersonalData) Update(firstName, lastName *string, age *uint8, birthDate *time.Time) error {
	if firstName != nil {
		data.FirstName = *firstName
	}
	if lastName != nil {
		data.LastName = *lastName
	}
	if birthDate != nil {
		data.BirthDate = birthDate
	}
	if age != nil {
		data.Age = age
	}

	return data.Validate()
}
