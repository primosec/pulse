package service

import (
	"context"
	"testing"
	"time"

	"github.com/primosec/pulse/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	users map[string]*domain.User
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[string]*domain.User)}
}

func (f *fakeUserRepo) Create(ctx context.Context, email, password, name string) (*domain.User, error) {
	u := &domain.User{ID: "fake-id", Email: email, Password: password, Name: name}
	f.users[email] = u
	return u, nil
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, ok := f.users[email]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return u, nil
}

func (f *fakeUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func TestRegisterAndLogin(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewAuthService(repo, "test-secret", time.Hour)
	ctx := context.Background()

	_, err := svc.Register(ctx, "a@a.com", "senha123456", "Ana")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	stored := repo.users["a@a.com"]
	if stored.Password == "senha123456" {
		t.Fatal("password was not hashed")
	}
	if bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte("senha123456")) != nil {
		t.Fatal("stored hash does not match original password")
	}

	token, err := svc.Login(ctx, "a@a.com", "senha123456")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	userID, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("token validation failed: %v", err)
	}
	if userID != "fake-id" {
		t.Fatalf("expected userID 'fake-id', got %q", userID)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewAuthService(repo, "test-secret", time.Hour)
	ctx := context.Background()

	_, _ = svc.Register(ctx, "b@b.com", "senhacerta123", "Bia")

	_, err := svc.Login(ctx, "b@b.com", "senhaerrada123")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}
