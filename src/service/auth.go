package service

import (
	"time"

	"dev.chaiyapluek.cloud.final.frontend/src/entity"
	"dev.chaiyapluek.cloud.final.frontend/src/repository"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	GetUserInfo(userId string) (string, error)
	LoginAttempt(email, password string) (string, error)
	Login(attemptId, email, code string) (string, string, error)
	RegisterAttempt(name, email, password string) (string, error)
	Register(attemptId, email, code string) (string, string, error)
}

type authService struct {
	authRepo       repository.AuthRepository
	accessTokenKey string
}

func NewAuthService(authRepo repository.AuthRepository, accessToken string) AuthService {
	return &authService{
		authRepo:       authRepo,
		accessTokenKey: accessToken,
	}
}

func (s *authService) GetUserInfo(userId string) (string, error) {
	user, err := s.authRepo.GetUserInfo(userId)
	if err != nil {
		return "", err
	}

	return user.Name, nil
}

func (s *authService) generateToken(userId string) (string, error) {
	payload := entity.AccessToken{
		UserId: userId,
		RegisteredClaims: &jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(s.accessTokenKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *authService) LoginAttempt(email, password string) (string, error) {
	attempt, err := s.authRepo.LoginAttempt(email, password)
	if err != nil {
		return "", err
	}

	return attempt.ID, nil
}

func (s *authService) Login(attemptId, email, code string) (string, string, error) {
	user, err := s.authRepo.Login(attemptId, email, code)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.generateToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, user.Name, nil
}

func (s *authService) RegisterAttempt(name, email, password string) (string, error) {
	attempt, err := s.authRepo.RegisterAttempt(name, email, password)
	if err != nil {
		return "", err
	}
	return attempt.ID, nil
}

func (s *authService) Register(attemptId, email, code string) (string, string, error) {
	user, err := s.authRepo.Register(attemptId, email, code)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.generateToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, user.Name, nil
}
