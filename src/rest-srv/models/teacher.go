package models

type Teacher struct {
	ID        int    `json:"id,omitempty" db:"id,primary_key,auto_increment"`
	FirstName string `json:"first_name,omitempty" db:"first_name,not_null"`
	LastName  string `json:"last_name,omitempty" db:"last_name,not_null"`
	Email     string `json:"email,omitempty" db:"email,not_null,unique"`
	Class     string `json:"class,omitempty" db:"class,not_null"`
	Subject   string `json:"subject,omitempty" db:"subject,not_null"`
}
