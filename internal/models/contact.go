package models

import (
	"time"
)

type Contact struct {
	ID             int        `json:"id" db:"id"`
	PhoneNumber    *string    `json:"phoneNumber" db:"phone_number"`
	Email          *string    `json:"email" db:"email"`
	LinkedID       *int       `json:"linkedId" db:"linked_id"`
	LinkPrecedence string     `json:"linkPrecedence" db:"link_precedence"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt      *time.Time `json:"deletedAt" db:"deleted_at"`
}

type IdentifyRequest struct {
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
}

type IdentifyResponse struct {
	Contact ContactInfo `json:"contact"`
}

type ContactInfo struct {
	PrimaryContactID    int      `json:"primaryContatctId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIDs []int    `json:"secondaryContactIds"`
}
