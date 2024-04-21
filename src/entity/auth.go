package entity

import "time"

type User struct {
	ID   string `json:"userId"`
	Name string `json:"name"`
}

type UserResponse struct {
	SuccessResponse
	Data *User `json:"data"`
}

type LoginRequest struct {
	AttemptID string     `json:"attemptId"`
	Email     string     `json:"email"`
	Code      string     `json:"code"`
	RequestAt *time.Time `json:"requestAt"`
}

type LoginAttemptRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginAttempt struct {
	ID string `json:"attemptId"`
}

type LoginAttemptResponse struct {
	SuccessResponse
	Data *LoginAttempt `json:"data"`
}

type RegisterAttemptRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterAttempt struct {
	ID string `json:"attemptId"`
}

type RegisterAttemptResponse struct {
	SuccessResponse
	Data *RegisterAttempt `json:"data"`
}
