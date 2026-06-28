package user

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewCreateFormRequestMapsFields(t *testing.T) {
	form := url.Values{}
	form.Set("email", "person@example.com")
	form.Set("password", "Password1!")
	form.Set("password_confirmation", "Password1!")
	form.Set("role", "captain")
	form.Set("first_name", "Pavel")
	form.Set("last_name", "Pavlov")
	form.Set("age", "31")

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
	if serviceInput.Age == nil || *serviceInput.Age != 31 {
		t.Fatalf("Age = %v, want 31", serviceInput.Age)
	}
	if serviceInput.Password != "Password1!" {
		t.Fatalf("Password = %q, want Password1!", serviceInput.Password)
	}
}

func TestNewCreateFormRequestRejectsPasswordMismatch(t *testing.T) {
	form := url.Values{}
	form.Set("email", "person@example.com")
	form.Set("password", "Password1!")
	form.Set("password_confirmation", "Password2!")
	form.Set("role", "user")
	form.Set("first_name", "Pavel")
	form.Set("last_name", "Pavlov")

	request := httptest.NewRequest("POST", "/ui/users", nil)
	request.PostForm = form

	if _, err := NewCreateFormRequest(request); err != errPasswordMismatch {
		t.Fatalf("NewCreateFormRequest() error = %v, want %v", err, errPasswordMismatch)
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
