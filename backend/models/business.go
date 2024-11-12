package models

import (
	"time"

	"github.com/google/uuid"
)

type BusinessStatus string

const (
	BUSINESS_STATUS_PENDING  BusinessStatus = "pending"
	BUSINESS_STATUS_ACTIVE   BusinessStatus = "active"
	BUSINESS_STATUS_DISABLED BusinessStatus = "disabled"
)

type BusinessUpdate struct {
	Website string `json:"website" db:"website" validate:"required,http_url"`
	Desc    string `json:"desc" db:"description" validate:"required,min=8,max=256"`
}

type BusinessCreate struct {
	Name string `json:"name" db:"name" validate:"required,min=3,max=64"`
	BusinessUpdate
}

type BusinessOverview struct {
	Id        uuid.UUID      `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Status    BusinessStatus `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

type Business struct {
	BusinessCreate
	Id        uuid.UUID      `json:"id" db:"id"`
	UserId    uuid.UUID      `json:"-" db:"user_id"`
	Status    BusinessStatus `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

type BusinessQueryParams struct {
	Status *BusinessStatus
	UserId *uuid.UUID
}
