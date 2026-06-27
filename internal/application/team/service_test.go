package team

import (
	"context"
	"errors"
	"testing"

	"task_tracker/internal/common_errors"
	domainUser "task_tracker/internal/domain/user"
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

type fakeTeamRepo struct {
	teamsByID       map[uuid.UUID]*Team
	teamsByName     map[string]*Team
	teamsByLeaderID map[uuid.UUID]*Team

	createErr error
	updateErr error
	deleteErr error

	created *Team
	updated *Team
	deleted uuid.UUID
}

func newFakeTeamRepo() *fakeTeamRepo {
	return &fakeTeamRepo{
		teamsByID:       make(map[uuid.UUID]*Team),
		teamsByName:     make(map[string]*Team),
		teamsByLeaderID: make(map[uuid.UUID]*Team),
	}
}

func (r *fakeTeamRepo) Create(ctx context.Context, team teamRepo.Team) (*teamRepo.Team, error) {
	if r.createErr != nil {
		return nil, r.createErr
	}
	r.created = &team
	r.teamsByID[team.ID] = &team
	r.teamsByName[team.Name] = &team
	if team.LeaderID != nil {
		r.teamsByLeaderID[*team.LeaderID] = &team
	}
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
	team, ok := r.teamsByLeaderID[id]
	if !ok {
		return nil, common_errors.ErrNotFound
	}
	return team, nil
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
	if r.updateErr != nil {
		return nil, r.updateErr
	}
	team.ID = id
	r.updated = &team
	r.teamsByID[id] = &team
	r.teamsByName[team.Name] = &team
	return &team, nil
}

func (r *fakeTeamRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	r.deleted = id
	return nil
}

type fakeUserRepo struct {
	users map[uuid.UUID]*userRepo.User
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
		if !user.IsDeleted() {
			users = append(users, user)
		}
	}
	return users, nil
}

func (r *fakeUserRepo) ListActiveProfiles(ctx context.Context) ([]*userRepo.ActiveUserProfile, error) {
	var profiles []*userRepo.ActiveUserProfile
	for _, user := range r.users {
		if !user.IsDeleted() {
			profiles = append(profiles, &userRepo.ActiveUserProfile{User: user})
		}
	}
	return profiles, nil
}

func (r *fakeUserRepo) ListProfiles(ctx context.Context) ([]*userRepo.ActiveUserProfile, error) {
	var profiles []*userRepo.ActiveUserProfile
	for _, user := range r.users {
		profiles = append(profiles, &userRepo.ActiveUserProfile{User: user})
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
	user.ID = id
	r.users[id] = &user
	return &user, nil
}

func (r *fakeUserRepo) Restore(ctx context.Context, id uuid.UUID) (*userRepo.User, error) {
	user, ok := r.users[id]
	if !ok {
		return nil, common_errors.ErrNotFound
	}
	user.IsActive = true
	user.DeletedAt = nil
	return user, nil
}

func (r *fakeUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	delete(r.users, id)
	return nil
}

func newTestService(teamRepo *fakeTeamRepo, userRepo *fakeUserRepo, tx *fakeTxManager) TeamService {
	return New(teamRepo, userRepo, zap.NewNop().Sugar(), tx)
}

func TestServiceCreateCreatesTeamWithoutLeader(t *testing.T) {
	teamRepo := newFakeTeamRepo()
	userRepo := newFakeUserRepo()
	tx := &fakeTxManager{}
	service := newTestService(teamRepo, userRepo, tx)

	created, err := service.Create(context.Background(), CreateTeamInput{
		Name:     "backend",
		Timezone: stringPtr("Europe/Moscow"),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if created == nil {
		t.Fatal("Create returned nil team")
	}
	if created.Name != "backend" {
		t.Fatalf("team name = %q, want %q", created.Name, "backend")
	}
	if !created.IsActive {
		t.Fatal("created team should be active")
	}
	if teamRepo.created == nil {
		t.Fatal("team repo Create was not called")
	}
	if tx.calls != 1 {
		t.Fatalf("tx calls = %d, want 1", tx.calls)
	}
}

func TestServiceCreateReturnsAlreadyExistsWhenNameTaken(t *testing.T) {
	existing := &Team{
		ID:       uuid.New(),
		Name:     "backend",
		IsActive: true,
	}
	teamRepo := newFakeTeamRepo()
	teamRepo.teamsByName[existing.Name] = existing
	service := newTestService(teamRepo, newFakeUserRepo(), &fakeTxManager{})

	created, err := service.Create(context.Background(), CreateTeamInput{Name: "backend"})
	if !errors.Is(err, ErrTeamAlreadyExists) {
		t.Fatalf("Create error = %v, want ErrTeamAlreadyExists", err)
	}
	if created != nil {
		t.Fatalf("created = %#v, want nil", created)
	}
	if teamRepo.created != nil {
		t.Fatal("team repo Create should not be called")
	}
}

func TestServiceCreateReturnsLeaderInactive(t *testing.T) {
	leaderID := uuid.New()
	userRepo := newFakeUserRepo()
	userRepo.users[leaderID] = &domainUser.User{
		ID:       leaderID,
		IsActive: false,
	}
	service := newTestService(newFakeTeamRepo(), userRepo, &fakeTxManager{})

	created, err := service.Create(context.Background(), CreateTeamInput{
		Name:     "backend",
		LeaderID: &leaderID,
	})
	if !errors.Is(err, ErrLeaderInactive) {
		t.Fatalf("Create error = %v, want ErrLeaderInactive", err)
	}
	if created != nil {
		t.Fatalf("created = %#v, want nil", created)
	}
}

func TestServiceUpdateReturnsInvalidInputForInvalidTimezone(t *testing.T) {
	teamID := uuid.New()
	teamRepo := newFakeTeamRepo()
	teamRepo.teamsByID[teamID] = &Team{
		ID:       teamID,
		Name:     "backend",
		IsActive: true,
	}
	service := newTestService(teamRepo, newFakeUserRepo(), &fakeTxManager{})

	updated, err := service.Update(context.Background(), teamID, &UpdateTeamInput{
		Name:     stringPtr("backend"),
		Timezone: stringPtr("not-a-timezone"),
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("Update error = %v, want ErrInvalidInput", err)
	}
	if updated != nil {
		t.Fatalf("updated = %#v, want nil", updated)
	}
	if teamRepo.updated != nil {
		t.Fatal("team repo Update should not be called")
	}
}

func TestServiceDeleteByIDDeletesExistingTeam(t *testing.T) {
	teamID := uuid.New()
	teamRepo := newFakeTeamRepo()
	teamRepo.teamsByID[teamID] = &Team{
		ID:       teamID,
		Name:     "backend",
		IsActive: true,
	}
	service := newTestService(teamRepo, newFakeUserRepo(), &fakeTxManager{})

	err := service.DeleteByID(context.Background(), teamID)
	if err != nil {
		t.Fatalf("DeleteByID returned error: %v", err)
	}
	if teamRepo.deleted != teamID {
		t.Fatalf("deleted id = %s, want %s", teamRepo.deleted, teamID)
	}
}

func stringPtr(value string) *string {
	return &value
}
