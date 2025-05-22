package entity

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Caknoooo/go-gin-clean-starter/helpers"
)

// User represents a system user with authentication and profile information.
type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=2,max=100"`
	Email      string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	TelpNumber string    `gorm:"type:varchar(20);index" json:"telp_number" validate:"omitempty,required,min=8,max=20"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-" validate:"required,min=8"`
	Role       string    `gorm:"type:varchar(50);not null;default:'user'" json:"role" validate:"required,oneof=user admin"`
	ImageUrl   string    `gorm:"type:varchar(255)" json:"image_url" validate:"omitempty,url"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`

	Timestamp
}

// validate is an instance of a validator used to validate structs based on defined tags.
var validate = validator.New()

// BeforeCreate is a GORM hook executed before creating a User record to hash the password, set default role, and validate the struct.
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	if u.Password != "" {
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}

	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if u.Role == "" {
		u.Role = "user"
	}

	if err := validate.Struct(u); err != nil {
		return err
	}

	return nil
}

// BeforeUpdate is a GORM hook executed before updating a User record to hash the password if a new password is provided.
func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	if u.Password != "" {
		u.Password, err = helpers.HashPassword(u.Password)
		if err != nil {
			return err
		}
	}
	return nil
}
