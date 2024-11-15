package models

import (
	"time"

	"github.com/google/uuid"
)

type PostUpdate struct {
	Title   string  `json:"title" db:"title" validate:"required,min=8,max=256"`
	Desc    string  `json:"desc" db:"description" validate:"required,min=8,max=256"`
	Pay     float32 `json:"pay" db:"pay" validate:"required,gt=0,usd"`
	TimeEst int     `json:"time_est" db:"time_est" validate:"required,gt=0"`
}

type PostStatus string

const (
	POST_STATUS_ACTIVE   PostStatus = "active"
	POST_STATUS_DISABLED PostStatus = "disabled"
	POST_STATUS_ARCHIVED PostStatus = "archived"
)

type PostCreate struct {
	PostUpdate
}

type PostOverview struct {
	Id        int        `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Status    PostStatus `json:"status" db:"status"`
	Pay       float32    `json:"pay" db:"pay"`
	TimeEst   int        `json:"time_est" db:"time_est"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

type Post struct {
	PostCreate
	BusinessId uuid.UUID  `json:"business_id" db:"business_id"`
	Id         int        `json:"id" db:"id"`
	Status     PostStatus `json:"status" db:"status"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

type PostQueryParams struct {
	Status     *PostStatus
	BusinessId *uuid.UUID
	UserId     *uuid.UUID
}

type ApplicationStatus string

const (
	APPLICATION_STATUS_PENDING    ApplicationStatus = "pending"
	APPLICATION_STATUS_ACCEPTED   ApplicationStatus = "accepted"
	APPLICATION_STATUS_REJECTED   ApplicationStatus = "rejected"
	APPLICATION_STATUS_WITHDRAWN  ApplicationStatus = "withdrawn"
	APPLICATION_STATUS_COMPLETED  ApplicationStatus = "completed"
	APPLICATION_STATUS_INCOMPLETE ApplicationStatus = "incompleted"
)

type PostApplicationData struct {
	User      UserOverview      `json:"user" db:"user"`
	Notes     string            `json:"notes" db:"notes"`
	Status    ApplicationStatus `json:"status" db:"status"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
}

type PostApplications struct {
	BusinessId   uuid.UUID             `json:"business_id" db:"business_id"`
	PostId       int                   `json:"post_id" db:"post_id"`
	Applications []PostApplicationData `json:"applications" db:"applications"`
}

type UserApplication struct {
	Post      PostOverview      `json:"post" db:"post"`
	Business  BusinessOverview  `json:"business" db:"business"`
	Status    ApplicationStatus `json:"status" db:"status"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
}

type UserApplicationQueryParams struct {
	UserId            *uuid.UUID
	ApplicationStatus *ApplicationStatus
	PostStatus        *PostStatus
}

type ApplicationNoteUpdate struct {
	Data string `json:"data" db:"data"`
}
