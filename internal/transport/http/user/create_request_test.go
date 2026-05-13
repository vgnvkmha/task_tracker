package user

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewCreateFormRequestMapsOptionalFields(t *testing.T) {
	form := url.Values{}
	form.Set("email", "person@example.com")
	form.Set("password", "Password1!")
	form.Set("role", "captain")
	form.Set("team_name", "Backend")
	form.Set("first_name", "Pavel")
	form.Set("last_name", "Pavlov")
	form.Set("age", "31")
	form.Set("birth_date", "1995-04-12")

	request := httptest.NewRequest("POST", "/ui/users", nil)
	request.PostForm = form

	input, err := NewCreateFormRequest(request)
	if err != nil {
		t.Fatalf("NewCreateFormRequest() error = %v", err)
	}

	serviceInput := input.ToServiceInput()

	if serviceInput.Email != "person@example.com" {
		t.Fatalf("Email = %q, want person@example.com", serviceInput.Email)
	}
	if serviceInput.Role != "captain" {
		t.Fatalf("Role = %q, want captain", serviceInput.Role)
	}
	if serviceInput.TeamName == nil || *serviceInput.TeamName != "Backend" {
		t.Fatalf("TeamName = %v, want Backend", serviceInput.TeamName)
	}
	if serviceInput.Age == nil || *serviceInput.Age != 31 {
		t.Fatalf("Age = %v, want 31", serviceInput.Age)
	}
	if serviceInput.BirthDate == nil || serviceInput.BirthDate.Format("2006-01-02") != "1995-04-12" {
		t.Fatalf("BirthDate = %v, want 1995-04-12", serviceInput.BirthDate)
	}
}

func TestNewCreateFormRequestRequiresFields(t *testing.T) {
	request := httptest.NewRequest("POST", "/ui/users", nil)
	request.PostForm = url.Values{
		"email": []string{"person@example.com"},
	}

	if _, err := NewCreateFormRequest(request); err == nil {
		t.Fatal("NewCreateFormRequest() error = nil, want required fields error")
	}
}
