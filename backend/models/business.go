package models

import (
	"net/url"
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
	Name    string `json:"name" db:"name" validate:"required,min=3,max=64"`
}

type BusinessCreate struct {
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
	UserId    uuid.UUID      `json:"user_id" db:"user_id"`
	Status    BusinessStatus `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}

type BusinessQueryParams struct {
	Status *BusinessStatus
	UserId *uuid.UUID
}

func (b *Business) URI(baseURL string) (string, error) {
	// TODO: Point to a specific business
	return url.JoinPath(baseURL, "account", "businesses")
}
