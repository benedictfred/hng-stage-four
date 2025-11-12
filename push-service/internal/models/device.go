package models

import (
	"time"
)

type Device struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	Platform  string    `json:"platform" db:"platform"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateDeviceRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=ios android web"`
}

type DeviceResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
	Platform string `json:"platform"`
	IsActive bool   `json:"is_active"`
}