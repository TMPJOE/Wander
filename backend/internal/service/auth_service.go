package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"wander/backend/internal/middleware"
	"wander/backend/internal/models"
	"wander/backend/internal/repository"
)

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtHours  int
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtHours int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtHours:  jwtHours,
	}
}

func (s *AuthService) Register(ctx context.Context, req models.UserCreateRequest) (*models.AuthResponse, error) {
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, fmt.Errorf("email already in use: %w", models.ErrConflict)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, req, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("register user: %w", err)
	}

	token, err := s.generateToken(*user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, models.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, models.ErrUnauthorized
	}

	token, err := s.generateToken(*user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &models.AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) generateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.jwtHours) * time.Hour)
	claims := &middleware.JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
