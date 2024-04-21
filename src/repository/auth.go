package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"dev.chaiyapluek.cloud.final.frontend/src/entity"
)

type AuthRepository interface {
	GetUserInfo(userId string) (*entity.User, error)
	Login(attemptId, email, code string) (*entity.User, error)
	LoginAttempt(email, password string) (*entity.LoginAttempt, error)
	Register(attemptId, email, code string) (*entity.User, error)
	RegisterAttempt(name, email, password string) (*entity.RegisterAttempt, error)
}

type authRepository struct {
	backendURL string
}

func NewAuthRepository(backendURL string) AuthRepository {
	return &authRepository{
		backendURL: backendURL,
	}
}

func (r *authRepository) GetUserInfo(userId string) (*entity.User, error) {
	req, _ := http.NewRequest("GET", r.backendURL+"/users/"+userId, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("auth repository get /users/"+userId+" error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("auth repository get /users/"+userId+" status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("auth repository get /users/"+userId+" decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var userResp entity.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		log.Println("auth repository get /users/"+userId+" decode UserResponse error", err)
		return nil, err
	}

	return userResp.Data, nil
}

func (r *authRepository) Login(attemptId, email, code string) (*entity.User, error) {
	now := time.Now()
	req := entity.LoginRequest{
		AttemptID: attemptId,
		Email:     email,
		Code:      code,
		RequestAt: &now,
	}
	reqBody, _ := json.Marshal(req)
	resp, err := http.Post(r.backendURL+"/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("auth repository post /auth/login error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("auth repository post /auth/login status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("auth repository post /auth/login decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var loginResp entity.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		log.Println("auth repository post /auth/login decode LocationResponse error", err)
		return nil, err
	}

	return loginResp.Data, nil
}

func (r *authRepository) LoginAttempt(email, password string) (*entity.LoginAttempt, error) {
	req := entity.LoginAttemptRequest{
		Email:    email,
		Password: password,
	}
	reqBody, _ := json.Marshal(req)
	resp, err := http.Post(r.backendURL+"/auth/login-attempt", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("auth repository post /auth/login-attempt error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("auth repository post /auth/login-attempt status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("auth repository post /auth/login-attempt decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var loginAttemptResp entity.LoginAttemptResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginAttemptResp); err != nil {
		log.Println("auth repository post /auth/login-attempt decode LocationResponse error", err)
		return nil, err
	}

	return loginAttemptResp.Data, nil
}

func (r *authRepository) Register(attemptId, email, code string) (*entity.User, error) {
	now := time.Now()
	req := entity.LoginRequest{
		AttemptID: attemptId,
		Email:     email,
		Code:      code,
		RequestAt: &now,
	}
	reqBody, _ := json.Marshal(req)
	resp, err := http.Post(r.backendURL+"/auth/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("auth repository post /auth/register error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("auth repository post /auth/register status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("auth repository post /auth/register decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var registerResp entity.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		log.Println("auth repository post /auth/register decode LocationResponse error", err)
		return nil, err
	}

	return registerResp.Data, nil
}

func (r *authRepository) RegisterAttempt(name, email, password string) (*entity.RegisterAttempt, error) {
	req := entity.RegisterAttemptRequest{
		Name:     name,
		Email:    email,
		Password: password,
	}
	reqBody, _ := json.Marshal(req)
	resp, err := http.Post(r.backendURL+"/auth/register-attempt", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println("auth repository post /auth/register-attempt error", err)
		return nil, err
	}

	if resp.StatusCode > 299 {
		log.Println("auth repository post /auth/register-attempt status code error", resp.StatusCode)
		var errorResp entity.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			log.Println("auth repository post /auth/register-attempt decode ErrorResponse error", err)
			return nil, err
		}
		return nil, errors.New(errorResp.Message)
	}

	var registerAttemptResp entity.RegisterAttemptResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerAttemptResp); err != nil {
		log.Println("auth repository post /auth/register-attempt decode LocationResponse error", err)
		return nil, err
	}

	return registerAttemptResp.Data, nil
}
