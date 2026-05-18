package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/primosec/pulse/internal/domain"
	"github.com/primosec/pulse/pkg/apperror"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, email, password, name string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

type AuthService struct {
	users     UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewAuthService(users UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (*domain.User, error) {
	if _, err := s.users.GetByEmail(ctx, email); err == nil {
		return nil, apperror.ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.ErrInternal
	}

	return s.users.Create(ctx, email, string(hash), name)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", apperror.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", apperror.ErrUnauthorized
	}

	return s.generateToken(user.ID)
}

func (s *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.jwtExpiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", apperror.ErrInternal
	}
	return signed, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return "", apperror.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", apperror.ErrUnauthorized
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", apperror.ErrUnauthorized
	}

	return sub, nil
}
