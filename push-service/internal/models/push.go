package models

import "time"

type PushNotification struct {
	ID           string         `json:"id" db:"id"`
	DeviceID     *string        `json:"device_id,omitempty" db:"device_id"`
	UserID       string         `json:"user_id" db:"user_id"`
	Title        string         `json:"title" db:"title"`
	Body         string         `json:"body" db:"body"`
	Image        *string        `json:"image,omitempty" db:"image"`
	Link         *string        `json:"link,omitempty" db:"link"`
	Data         map[string]any `json:"data,omitempty" db:"data"`
	Status       string         `json:"status" db:"status"`
	ErrorMessage *string        `json:"error_message,omitempty" db:"error_message"`
	SentAt       *time.Time     `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}

type SendPushRequest struct {
	UserID    string         `json:"user_id" binding:"required"`
	Title     string         `json:"title" binding:"required"`
	Body      string         `json:"body" binding:"required"`
	Image     *string        `json:"image,omitempty"`
	Link      *string        `json:"link,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
	Platforms []string       `json:"platforms,omitempty"` // Filter by specific platforms
}

type BulkPushRequest struct {
	UserIDs []string       `json:"user_ids" binding:"required"`
	Title   string         `json:"title" binding:"required"`
	Body    string         `json:"body" binding:"required"`
	Data    map[string]any `json:"data,omitempty"`
}
