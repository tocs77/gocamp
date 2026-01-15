package models

import "database/sql"

type Exec struct {
	ID                     int            `json:"id,omitempty"`
	FirstName              string         `json:"first_name,omitempty"`
	LastName               string         `json:"last_name,omitempty"`
	Email                  string         `json:"email,omitempty"`
	Username               string         `json:"username,omitempty"`
	Password               string         `json:"password,omitempty"`
	PasswordChangedAt      sql.NullString `json:"passwordChangedAt,omitempty"`
	UserCreatedAt          sql.NullString `json:"userCreatedAt,omitempty"`
	PasswordResetToken     sql.NullString `json:"passwordResetToken,omitempty"`
	PasswordResetExpiresAt sql.NullString `json:"passwordResetExpiresAt,omitempty"`
	InactiveStatus         bool           `json:"inactiveStatus,omitempty"`
	Role                   string         `json:"role,omitempty"`
}
