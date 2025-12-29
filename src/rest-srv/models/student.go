package models

import "rest-srv/utility"

type Student struct {
	ID        int    `json:"id,omitempty" db:"id,primary_key,auto_increment"`
	FirstName string `json:"first_name,omitempty" db:"first_name,not_null"`
	LastName  string `json:"last_name,omitempty" db:"last_name,not_null"`
	Email     string `json:"email,omitempty" db:"email,not_null,unique"`
	Class     string `json:"class,omitempty" db:"class,not_null"`
}

func (s *Student) Validate() error {
	return utility.ValidateBlank(s)
}
