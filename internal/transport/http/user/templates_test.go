package user

import (
	"bytes"
	"strings"
	"testing"
)

func TestTemplatesParse(t *testing.T) {
	if Templates() == nil {
		t.Fatal("Templates returned nil")
	}
}

func TestCabinetDeleteActionAvailability(t *testing.T) {
	tmpl := Templates()

	tests := []struct {
		name       string
		isActive   bool
		want       string
		wantAbsent string
	}{
		{
			name:       "active user can submit delete",
			isActive:   true,
			want:       `<button class="danger" type="submit">Удалить пользователя</button>`,
			wantAbsent: `class="danger unavailable"`,
		},
		{
			name:       "deleted user sees unavailable delete action",
			isActive:   false,
			want:       `data-unavailable-action="Пользователь уже удален."`,
			wantAbsent: `<button class="danger" type="submit">Удалить пользователя</button>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			err := tmpl.ExecuteTemplate(&output, "user_cabinet_page", map[string]any{
				"title": "Личный кабинет",
				"profile": CabinetView{
					ID:       "user-id",
					Email:    "user@example.com",
					Role:     "user",
					IsActive: tt.isActive,
				},
			})
			if err != nil {
				t.Fatalf("ExecuteTemplate() error = %v", err)
			}

			html := output.String()
			if !strings.Contains(html, tt.want) {
				t.Fatalf("rendered cabinet does not contain %q", tt.want)
			}
			if strings.Contains(html, tt.wantAbsent) {
				t.Fatalf("rendered cabinet contains %q", tt.wantAbsent)
			}
		})
	}
}
