package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostUpdate struct {
	Title string `json:"title" db:"title" validate:"required,min=8,max=256"`
	Desc  string `json:"desc" db:"description" validate:"required,min=8,max=256"`
}

type PostStatus int

const (
	POST_STATUS_ACTIVE PostStatus = iota
	POST_STATUS_DISABLED
	POST_STATUS_ARCHIVED
)

func (ps PostStatus) String() string {
	switch ps {
	case POST_STATUS_ACTIVE:
		return "active"
	case POST_STATUS_DISABLED:
		return "disabled"
	case POST_STATUS_ARCHIVED:
		return "archived"
	default:
		return "unknown"
	}
}

func (s *PostStatus) ScanText(value pgtype.Text) error {
	switch value.String {
	case "active":
		*s = POST_STATUS_ACTIVE
		return nil
	case "disabled":
		*s = POST_STATUS_DISABLED
		return nil
	case "archived":
		*s = POST_STATUS_ARCHIVED
		return nil
	default:
		return errors.New("Unsupported value scanning post status")
	}
}

func (s PostStatus) TextValue() (pgtype.Text, error) {
	val := pgtype.Text{}
	err := val.Scan(s.String())
	return val, err
}

type PostCreate struct {
	PostUpdate
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

type ApplicationStatus int

const (
	APPLICATION_STATUS_PENDING ApplicationStatus = iota
	APPLICATION_STATUS_ACCEPTED
	APPLICATION_STATUS_REJECTED
	APPLICATION_STATUS_WITHDRAWN
	APPLICATION_STATUS_COMPLETED
)

func (ps ApplicationStatus) String() string {
	switch ps {
	case APPLICATION_STATUS_PENDING:
		return "pending"
	case APPLICATION_STATUS_ACCEPTED:
		return "accepted"
	case APPLICATION_STATUS_REJECTED:
		return "rejected"
	case APPLICATION_STATUS_WITHDRAWN:
		return "withdrawn"
	case APPLICATION_STATUS_COMPLETED:
		return "completed"
	default:
		return "unknown"
	}
}

func (s *ApplicationStatus) ScanText(value pgtype.Text) error {
	switch value.String {
	case "pending":
		*s = APPLICATION_STATUS_PENDING
		return nil
	case "accepted":
		*s = APPLICATION_STATUS_ACCEPTED
		return nil
	case "rejected":
		*s = APPLICATION_STATUS_REJECTED
		return nil
	case "withdrawn":
		*s = APPLICATION_STATUS_WITHDRAWN
		return nil
	case "completed":
		*s = APPLICATION_STATUS_COMPLETED
		return nil
	default:
		return errors.New("Unsupported value scanning application status")
	}
}

func (s ApplicationStatus) TextValue() (pgtype.Text, error) {
	val := pgtype.Text{}
	err := val.Scan(s.String())
	return val, err
}

type PostApplicationData struct {
	User   UserOverview      `json:"user" db:"user"`
	Notes  string            `json:"notes" db:"notes"`
	Status ApplicationStatus `json:"status" db:"status"`
}

type PostApplications struct {
	BusinessId   uuid.UUID             `json:"business_id" db:"business_id"`
	PostId       int                   `json:"post_id" db:"post_id"`
	Applications []PostApplicationData `json:"applications" db:"applications"`
}

type ApplicationNoteUpdate struct {
	Data string `json:"data" db:"data"`
}
