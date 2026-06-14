package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"task_tracker/internal/common_errors"
	"task_tracker/internal/domain/auth"
	personaldata "task_tracker/internal/domain/personal_data"
	domainUser "task_tracker/internal/domain/user"
	valueobjects "task_tracker/internal/domain/value_objects"
	teamRepo "task_tracker/internal/repo/team"
	userRepo "task_tracker/internal/repo/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type fakeTxManager struct {
	calls int
}

func (m *fakeTxManager) WithTx(ctx context.Context, fn func(context.Context) error) error {
	m.calls++
	return fn(ctx)
}

type fakeUserRepo struct {
	users     map[uuid.UUID]*userRepo.User
	updateErr error
	updated   *userRepo.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		users: make(map[uuid.UUID]*userRepo.User),
	}
}

func (r *fakeUserRepo) Create(ctx context.Context, user userRepo.User) (*userRepo.User, error) {
	r.users[user.ID] = &user
	return &user, nil
}

func (r *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*userRepo.User, error) {
	for _, user := range r.users {
		if string(user.Email) == email {
			return user, nil
		}
	}
	return nil, common_errors.ErrNotFound
}

func (r *fakeUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*userRepo.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, common_errors.ErrNotFound
	}
	return user, nil
}

func (r *fakeUserRepo) ListActive(ctx context.Context) ([]*userRepo.User, error) {
	var users []*userRepo.User
	for _, user := range r.users {
		if user.IsActive {
			users = append(users, user)
		}
	}
	return users, nil
}

func (r *fakeUserRepo) ListActiveProfiles(ctx context.Context) ([]*userRepo.ActiveUserProfile, error) {
	var profiles []*userRepo.ActiveUserProfile
	for _, user := range r.users {
		if user.IsActive {
			profiles = append(profiles, &userRepo.ActiveUserProfile{
				User:         user,
				PersonalData: personaldata.PersonalData{Id: user.PersonalDataID},
			})
		}
	}
	return profiles, nil
}

func (r *fakeUserRepo) List(ctx context.Context) ([]*userRepo.User, error) {
	var users []*userRepo.User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *fakeUserRepo) Update(ctx context.Context, id uuid.UUID, user userRepo.User) (*userRepo.User, error) {
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	user.ID = id
	r.updated = &user
	r.users[id] = &user
	return &user, nil
}

func (r *fakeUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(r.users, id)
	return nil
}

type fakePersonalDataRepo struct {
	data    map[uuid.UUID]userRepo.PersonalData
	updated *userRepo.PersonalData
}

func newFakePersonalDataRepo() *fakePersonalDataRepo {
	return &fakePersonalDataRepo{
		data: make(map[uuid.UUID]userRepo.PersonalData),
	}
}

func (r *fakePersonalDataRepo) Create(ctx context.Context, data userRepo.PersonalData) (userRepo.PersonalData, error) {
	r.data[data.Id] = data
	return data, nil
}

func (r *fakePersonalDataRepo) Get(ctx context.Context, dataID uuid.UUID) (userRepo.PersonalData, error) {
	data, ok := r.data[dataID]
	if !ok {
		return userRepo.PersonalData{}, common_errors.ErrNotFound
	}
	return data, nil
}

func (r *fakePersonalDataRepo) Update(ctx context.Context, data userRepo.PersonalData) (userRepo.PersonalData, error) {
	r.updated = &data
	r.data[data.Id] = data
	return data, nil
}

type fakeTeamRepo struct {
	teamsByID   map[uuid.UUID]*teamRepo.Team
	teamsByName map[string]*teamRepo.Team
}

func newFakeTeamRepo() *fakeTeamRepo {
	return &fakeTeamRepo{
		teamsByID:   make(map[uuid.UUID]*teamRepo.Team),
		teamsByName: make(map[string]*teamRepo.Team),
	}
}

func (r *fakeTeamRepo) Create(ctx context.Context, team teamRepo.Team) (*teamRepo.Team, error) {
	r.teamsByID[team.ID] = &team
	r.teamsByName[team.Name] = &team
	return &team, nil
}

func (r *fakeTeamRepo) GetByID(ctx context.Context, id uuid.UUID) (*teamRepo.Team, error) {
	team, ok := r.teamsByID[id]
	if !ok {
		return nil, common_errors.ErrNotFound
	}
	return team, nil
}

func (r *fakeTeamRepo) GetByName(ctx context.Context, name string) (*teamRepo.Team, error) {
	team, ok := r.teamsByName[name]
	if !ok {
		return nil, common_errors.ErrNotFound
	}
	return team, nil
}

func (r *fakeTeamRepo) GetByLeaderID(ctx context.Context, id uuid.UUID) (*teamRepo.Team, error) {
	return nil, common_errors.ErrNotFound
}

func (r *fakeTeamRepo) ListActive(ctx context.Context) ([]*teamRepo.Team, error) {
	var teams []*teamRepo.Team
	for _, team := range r.teamsByID {
		if team.IsActive {
			teams = append(teams, team)
		}
	}
	return teams, nil
}

func (r *fakeTeamRepo) List(ctx context.Context) ([]*teamRepo.Team, error) {
	var teams []*teamRepo.Team
	for _, team := range r.teamsByID {
		teams = append(teams, team)
	}
	return teams, nil
}

func (r *fakeTeamRepo) Update(ctx context.Context, id uuid.UUID, team teamRepo.Team) (*teamRepo.Team, error) {
	team.ID = id
	r.teamsByID[id] = &team
	r.teamsByName[team.Name] = &team
	return &team, nil
}

func (r *fakeTeamRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(r.teamsByID, id)
	return nil
}

func TestServiceUpdateAllowsPartialUserPatch(t *testing.T) {
	svc, users, data, _ := newUpdateTestService()
	existing := addUserFixture(t, users, data, nil, valueobjects.User)

	firstName := "Petr"
	updated, err := svc.Update(context.Background(), managerActor(), UpdateUserInput{
		UserID:    existing.ID,
		FirstName: &firstName,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Email != existing.Email {
		t.Fatalf("email = %q, want preserved %q", updated.Email, existing.Email)
	}
	if updated.Role != existing.Role {
		t.Fatalf("role = %q, want preserved %q", updated.Role, existing.Role)
	}
	if data.updated == nil || data.updated.FirstName != firstName {
		t.Fatalf("personal data first name was not updated")
	}
}

func TestServiceUpdateAppliesTeamID(t *testing.T) {
	svc, users, data, teams := newUpdateTestService()
	existing := addUserFixture(t, users, data, nil, valueobjects.User)
	teamID := uuid.New()
	teams.teamsByID[teamID] = &teamRepo.Team{
		ID:       teamID,
		Name:     "backend",
		IsActive: true,
	}

	updated, err := svc.Update(context.Background(), managerActor(), UpdateUserInput{
		UserID: existing.ID,
		TeamId: &teamID,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.TeamID == nil || *updated.TeamID != teamID {
		t.Fatalf("team_id = %v, want %s", updated.TeamID, teamID)
	}
}

func TestServiceUpdateAllowsSelfUpdateForNonManager(t *testing.T) {
	svc, users, data, _ := newUpdateTestService()
	existing := addUserFixture(t, users, data, nil, valueobjects.User)
	firstName := "Self"

	updated, err := svc.Update(context.Background(), auth.Actor{
		ID:   existing.ID,
		Role: valueobjects.User,
	}, UpdateUserInput{
		UserID:    existing.ID,
		FirstName: &firstName,
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.ID != existing.ID {
		t.Fatalf("updated id = %s, want %s", updated.ID, existing.ID)
	}
	if data.updated == nil || data.updated.FirstName != firstName {
		t.Fatal("self update did not update personal data")
	}
}

func TestServiceUpdateRejectsNonManagerUpdatingAnotherUser(t *testing.T) {
	svc, users, data, _ := newUpdateTestService()
	existing := addUserFixture(t, users, data, nil, valueobjects.User)
	firstName := "Other"

	updated, err := svc.Update(context.Background(), auth.Actor{
		ID:   uuid.New(),
		Role: valueobjects.User,
	}, UpdateUserInput{
		UserID:    existing.ID,
		FirstName: &firstName,
	})
	if !errors.Is(err, ErrOnlyManagersCanModify) {
		t.Fatalf("Update error = %v, want ErrOnlyManagersCanModify", err)
	}
	if updated != nil {
		t.Fatalf("updated = %#v, want nil", updated)
	}
	if data.updated != nil {
		t.Fatal("personal data should not be updated")
	}
	if users.updated != nil {
		t.Fatal("user repo Update should not be called")
	}
}

func TestServiceUpdateRejectsNonManagerSelfRoleUpdate(t *testing.T) {
	svc, users, data, _ := newUpdateTestService()
	existing := addUserFixture(t, users, data, nil, valueobjects.User)
	role := string(valueobjects.Admin)

	updated, err := svc.Update(context.Background(), auth.Actor{
		ID:   existing.ID,
		Role: valueobjects.User,
	}, UpdateUserInput{
		UserID: existing.ID,
		Role:   &role,
	})
	if !errors.Is(err, ErrOnlyManagersCanModify) {
		t.Fatalf("Update error = %v, want ErrOnlyManagersCanModify", err)
	}
	if updated != nil {
		t.Fatalf("updated = %#v, want nil", updated)
	}
	if data.updated != nil {
		t.Fatal("personal data should not be updated")
	}
	if users.updated != nil {
		t.Fatal("user repo Update should not be called")
	}
}

func TestServiceUpdateRejectsNonManagerSelfTeamIDUpdate(t *testing.T) {
	svc, users, data, _ := newUpdateTestService()
	existing := addUserFixture(t, users, data, nil, valueobjects.User)
	teamID := uuid.New()

	updated, err := svc.Update(context.Background(), auth.Actor{
		ID:   existing.ID,
		Role: valueobjects.User,
	}, UpdateUserInput{
		UserID: existing.ID,
		TeamId: &teamID,
	})
	if !errors.Is(err, ErrOnlyManagersCanModify) {
		t.Fatalf("Update error = %v, want ErrOnlyManagersCanModify", err)
	}
	if updated != nil {
		t.Fatalf("updated = %#v, want nil", updated)
	}
	if data.updated != nil {
		t.Fatal("personal data should not be updated")
	}
	if users.updated != nil {
		t.Fatal("user repo Update should not be called")
	}
}

func TestServiceUpdateReturnsAlreadyExistsForDuplicateEmail(t *testing.T) {
	svc, users, data, _ := newUpdateTestService()
	existing := addUserFixture(t, users, data, nil, valueobjects.User)
	users.updateErr = common_errors.ErrAlreadyExists
	email := "used@example.com"

	updated, err := svc.Update(context.Background(), managerActor(), UpdateUserInput{
		UserID: existing.ID,
		Email:  &email,
	})
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("Update error = %v, want ErrUserAlreadyExists", err)
	}
	if updated != nil {
		t.Fatalf("updated = %#v, want nil", updated)
	}
}

func TestServiceUpdateRejectsManagerWithoutTeam(t *testing.T) {
	svc, users, data, _ := newUpdateTestService()
	existing := addUserFixture(t, users, data, nil, valueobjects.User)
	role := string(valueobjects.Admin)

	updated, err := svc.Update(context.Background(), managerActor(), UpdateUserInput{
		UserID: existing.ID,
		Role:   &role,
	})
	if !errors.Is(err, domainUser.ErrManagerMustHaveTeam) {
		t.Fatalf("Update error = %v, want ErrManagerMustHaveTeam", err)
	}
	if updated != nil {
		t.Fatalf("updated = %#v, want nil", updated)
	}
	if users.updated != nil {
		t.Fatal("user repo Update should not be called")
	}
}

func TestServiceCreateRegisterRejectsMissingRoleBeforePersistingPersonalData(t *testing.T) {
	svc, users, data, _ := newUpdateTestService()

	created, err := svc.CreateRegister(context.Background(), CreateUserInput{
		Email:     "new@example.com",
		Password:  "Strong1!",
		FirstName: "Ivan",
		LastName:  "Petrov",
	})
	if !errors.Is(err, domainUser.ErrInvalidRole) {
		t.Fatalf("CreateRegister error = %v, want ErrInvalidRole", err)
	}
	if created != nil {
		t.Fatalf("created = %#v, want nil", created)
	}
	if len(data.data) != 0 {
		t.Fatal("personal data should not be persisted when user domain validation fails")
	}
	if len(users.users) != 0 {
		t.Fatal("user should not be persisted when domain validation fails")
	}
}

func newUpdateTestService() (UserService, *fakeUserRepo, *fakePersonalDataRepo, *fakeTeamRepo) {
	users := newFakeUserRepo()
	data := newFakePersonalDataRepo()
	teams := newFakeTeamRepo()
	service := New(users, data, teams, zap.NewNop().Sugar(), &fakeTxManager{})
	return service, users, data, teams
}

func addUserFixture(t *testing.T, users *fakeUserRepo, data *fakePersonalDataRepo, teamID *uuid.UUID, role valueobjects.Role) *domainUser.User {
	t.Helper()

	personalDataID := uuid.New()
	data.data[personalDataID] = personaldata.PersonalData{
		Id:        personalDataID,
		FirstName: "Ivan",
		LastName:  "Petrov",
		Age:       uint8Ptr(25),
		BirthDate: timePtr(time.Date(1999, 5, 10, 0, 0, 0, 0, time.UTC)),
	}

	user := &domainUser.User{
		ID:             uuid.New(),
		TeamID:         teamID,
		Email:          valueobjects.Email("ivan@example.com"),
		Password:       mustPassword(t, "Strong1!"),
		Role:           role,
		PersonalDataID: personalDataID,
		IsActive:       true,
	}
	users.users[user.ID] = user
	return user
}

func managerActor() auth.Actor {
	return auth.Actor{
		ID:   uuid.New(),
		Role: valueobjects.Admin,
	}
}

func mustPassword(t *testing.T, raw string) valueobjects.Password {
	t.Helper()

	password, err := valueobjects.NewPassword(raw)
	if err != nil {
		t.Fatalf("NewPassword returned error: %v", err)
	}
	return password
}

func uint8Ptr(value uint8) *uint8 {
	return &value
}

func timePtr(value time.Time) *time.Time {
	return &value
}
