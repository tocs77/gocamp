package models

import (
	"encoding/json"

	"rest-srv/utility"
)

type Exec struct {
	ID                 int                `json:"id,omitempty" db:"id,primary_key,auto_increment"`
	FirstName          string             `json:"first_name,omitempty" db:"first_name,not_null"`
	LastName           string             `json:"last_name,omitempty" db:"last_name,not_null"`
	Email              string             `json:"email,omitempty" db:"email,not_null,unique"`
	Username           string             `json:"username,omitempty" db:"username,not_null,unique"`
	Password           string             `json:"password,omitempty" db:"password,not_null"`
	PasswordChangedAt  utility.NullString `json:"password_changed_at,omitempty" db:"password_changed_at"`
	UserCreatedAt      utility.NullString `json:"user_created_at,omitempty" db:"user_created_at"`
	PasswordResetToken utility.NullString `json:"password_reset_token,omitempty" db:"password_reset_token"`
	InactiveStatus     bool               `json:"inactive_status,omitempty" db:"inactive_status,not_null"`
	Role               string             `json:"role,omitempty" db:"role,not_null"`
}

func (e Exec) MarshalJSON() ([]byte, error) {
	type ExecAlias Exec
	return json.Marshal(&struct {
		ExecAlias
		Password string `json:"password,omitempty"`
	}{
		ExecAlias: ExecAlias(e),
		Password:  "",
	})
}

func (e *Exec) Validate() error {
	return utility.ValidateBlank(e)
}
