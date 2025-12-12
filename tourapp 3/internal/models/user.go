package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New() 

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:100;not null;unique" json:"username" validate:"required,min=3,max=50"`
	Email     string    `gorm:"size:100;not null;unique" json:"email" validate:"required,email"`
	Password  string    `gorm:"not null" json:"password" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


func (u *User) Validate() error {
	return validate.Struct(u)
}
