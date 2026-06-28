package user

import (
	"testing"

	userApplication "task_tracker/internal/application/user"
)

func TestMapUIErrorUserAlreadyExists(t *testing.T) {
	got := mapUIError(userApplication.ErrUserAlreadyExists)
	want := "Пользователь с такой почтой уже существует."

	if got != want {
		t.Fatalf("mapUIError() = %q, want %q", got, want)
	}
}

func TestMapUIFormErrorRequiredFieldsMissing(t *testing.T) {
	got := mapUIFormError(errRequiredFieldsMissing)
	want := "Заполните обязательные поля: электронную почту, роль, имя, фамилию и оба поля пароля."

	if got != want {
		t.Fatalf("mapUIFormError() = %q, want %q", got, want)
	}
}

func TestMapUIFormErrorPasswordMismatch(t *testing.T) {
	got := mapUIFormError(errPasswordMismatch)
	want := "Пароли не совпадают."

	if got != want {
		t.Fatalf("mapUIFormError() = %q, want %q", got, want)
	}
}
