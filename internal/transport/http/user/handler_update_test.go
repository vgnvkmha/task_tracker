package user

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	userApplication "task_tracker/internal/application/user"
	"task_tracker/internal/domain/auth"
	domainUser "task_tracker/internal/domain/user"
	valueobjects "task_tracker/internal/domain/value_objects"
	"task_tracker/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fakeUserService struct {
	updateCalled bool
	updateInput  userApplication.UpdateUserInput
	updateActor  auth.Actor
}

func (s *fakeUserService) CreateRegister(ctx context.Context, input userApplication.CreateUserInput) (*userApplication.User, error) {
	return nil, nil
}

func (s *fakeUserService) CreateByActor(ctx context.Context, actor auth.Actor, input userApplication.CreateUserInput) (*userApplication.User, error) {
	return nil, nil
}

func (s *fakeUserService) Update(ctx context.Context, actor auth.Actor, input userApplication.UpdateUserInput) (*userApplication.User, error) {
	s.updateCalled = true
	s.updateActor = actor
	s.updateInput = input

	role := valueobjects.User
	if input.Role != nil {
		role = valueobjects.Role(*input.Role)
	}
	email := valueobjects.Email("updated@example.com")
	if input.Email != nil {
		email = valueobjects.Email(*input.Email)
	}

	return &domainUser.User{
		ID:             input.UserID,
		TeamID:         input.TeamId,
		Email:          email,
		Password:       mustHandlerPassword(nil),
		Role:           role,
		PersonalDataID: uuid.New(),
		IsActive:       true,
	}, nil
}

func (s *fakeUserService) Login(ctx context.Context, email string, password string) (*userApplication.User, error) {
	return nil, nil
}

func (s *fakeUserService) Restore(ctx context.Context, email string, password string) (*userApplication.User, error) {
	return nil, nil
}

func (s *fakeUserService) GetByID(ctx context.Context, id uuid.UUID) (*userApplication.User, error) {
	return nil, nil
}

func (s *fakeUserService) GetProfileByID(ctx context.Context, id uuid.UUID) (*userApplication.Profile, error) {
	return nil, nil
}

func (s *fakeUserService) ListActive(ctx context.Context) ([]*userApplication.User, error) {
	return nil, nil
}

func (s *fakeUserService) ListActiveProfiles(ctx context.Context) ([]*userApplication.Profile, error) {
	return nil, nil
}

func (s *fakeUserService) ListProfiles(ctx context.Context) ([]*userApplication.Profile, error) {
	return nil, nil
}

func (s *fakeUserService) List(ctx context.Context) ([]*userApplication.User, error) {
	return nil, nil
}

func (s *fakeUserService) DeleteByID(ctx context.Context, actor auth.Actor, id uuid.UUID) error {
	return nil
}

func TestHandlerUpdateAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	managerID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()
	teamID := uuid.New()

	tests := []struct {
		name        string
		actorID     uuid.UUID
		actorRole   valueobjects.Role
		body        string
		withActor   bool
		wantStatus  int
		wantService bool
	}{
		{
			name:        "manager updates any user with role and team_id",
			actorID:     managerID,
			actorRole:   valueobjects.Admin,
			body:        fmt.Sprintf(`{"user_id":%q,"role":"user","team_id":%q}`, userID.String(), teamID.String()),
			withActor:   true,
			wantStatus:  http.StatusOK,
			wantService: true,
		},
		{
			name:        "manager updates any user without role and team_id",
			actorID:     managerID,
			actorRole:   valueobjects.Admin,
			body:        fmt.Sprintf(`{"user_id":%q,"first_name":"Ivan"}`, userID.String()),
			withActor:   true,
			wantStatus:  http.StatusOK,
			wantService: true,
		},
		{
			name:        "user updates self without role and team_id",
			actorID:     userID,
			actorRole:   valueobjects.User,
			body:        fmt.Sprintf(`{"user_id":%q,"last_name":"Petrov"}`, userID.String()),
			withActor:   true,
			wantStatus:  http.StatusOK,
			wantService: true,
		},
		{
			name:        "user updates self with role",
			actorID:     userID,
			actorRole:   valueobjects.User,
			body:        fmt.Sprintf(`{"user_id":%q,"role":"admin"}`, userID.String()),
			withActor:   true,
			wantStatus:  http.StatusForbidden,
			wantService: false,
		},
		{
			name:        "user updates self with team_id",
			actorID:     userID,
			actorRole:   valueobjects.User,
			body:        fmt.Sprintf(`{"user_id":%q,"team_id":%q}`, userID.String(), teamID.String()),
			withActor:   true,
			wantStatus:  http.StatusForbidden,
			wantService: false,
		},
		{
			name:        "user updates another user",
			actorID:     otherUserID,
			actorRole:   valueobjects.User,
			body:        fmt.Sprintf(`{"user_id":%q,"email":"new@example.com"}`, userID.String()),
			withActor:   true,
			wantStatus:  http.StatusForbidden,
			wantService: false,
		},
		{
			name:        "actor missing",
			body:        fmt.Sprintf(`{"user_id":%q,"email":"new@example.com"}`, userID.String()),
			withActor:   false,
			wantStatus:  http.StatusForbidden,
			wantService: false,
		},
		{
			name:        "invalid birth_date format",
			actorID:     userID,
			actorRole:   valueobjects.User,
			body:        fmt.Sprintf(`{"user_id":%q,"birth_date":"1999-05-10"}`, userID.String()),
			withActor:   true,
			wantStatus:  http.StatusUnprocessableEntity,
			wantService: false,
		},
		{
			name:        "partial update with only email",
			actorID:     userID,
			actorRole:   valueobjects.User,
			body:        fmt.Sprintf(`{"user_id":%q,"email":"new@example.com"}`, userID.String()),
			withActor:   true,
			wantStatus:  http.StatusOK,
			wantService: true,
		},
		{
			name:        "partial update with only password",
			actorID:     userID,
			actorRole:   valueobjects.User,
			body:        fmt.Sprintf(`{"user_id":%q,"password":"Strong1!"}`, userID.String()),
			withActor:   true,
			wantStatus:  http.StatusOK,
			wantService: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &fakeUserService{}
			handler := New(service)

			router := gin.New()
			if tt.withActor {
				router.PATCH("/user/update", middleware.ActorMiddleware(nil, true), handler.Update)
			} else {
				router.PATCH("/user/update", handler.Update)
			}

			req := httptest.NewRequest(http.MethodPatch, "/user/update", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.withActor {
				req.Header.Set("X-User-ID", tt.actorID.String())
				req.Header.Set("X-User-Role", string(tt.actorRole))
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if service.updateCalled != tt.wantService {
				t.Fatalf("service Update called = %t, want %t", service.updateCalled, tt.wantService)
			}
		})
	}
}

func mustHandlerPassword(t *testing.T) valueobjects.Password {
	password, err := valueobjects.NewPassword("Strong1!")
	if err != nil {
		if t != nil {
			t.Fatalf("NewPassword returned error: %v", err)
		}
		panic(err)
	}
	return password
}
