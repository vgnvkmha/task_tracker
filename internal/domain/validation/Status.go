package validation

import (
	"fmt"
	"strings"
	valueobjects "task_tracker/internal/domain/models/value_objects"
)

func ParseStatus(input string) (valueobjects.Status, error) {
	s := valueobjects.Status(strings.ToLower(strings.TrimSpace(input)))

	if !s.IsValid() {
		return "", fmt.Errorf("invalid status: %s.Should be In Progress, To do, etc.", input)
	}

	return s, nil
}
