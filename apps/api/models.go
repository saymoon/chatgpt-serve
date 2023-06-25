package main

import (
	"time"
)

const (
	TASK_STATUS_PENDING    = "pending"
	TASK_STATUS_PROCESSING = "processing"
	TASK_STATUS_COMPLETED  = "completed"
	TASK_STATUS_ERROR      = "error"
)

type Task struct {
	Id             string    `json:"id"`
	InstanceId     string    `json:"instance_id"`
	ConversationId string    `json:"conversation_id"`
	Model          string    `json:"model"`
	Prompt         string    `json:"prompt"`
	Response       string    `json:"response"`
	Status         string    `json:"status"`
	ErrorMessage   string    `json:"error_message"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Token struct {
	Token   string `json:"token"`
	IsAdmin bool   `json:"is_admin"`
}

