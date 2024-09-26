package datamodels

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is our User example model.
// Keep note that the tags for public-use (for our web app)
// should be kept in other file like "web/viewmodels/user.go"
// which could wrap by embedding the datamodels.User or
// define completely new fields instead but for the sake
// of the example, we will use this datamodel
// as the only one User model in our application.
type User struct {
	gorm.Model
	UserId         string `json:"user_id" form:"user_id" gorm:"unique"`
	Nickname       string `json:"nickname" form:"nickname"`
	Username       string `json:"username" form:"username" gorm:"unique"`
	HashedPassword []byte `json:"-" form:"-"`
}

// IsValid can do some very very simple "low-level" data validations.
func (u User) IsValid() bool {
	return u.ID > 0
}

// GeneratePassword will generate a hashed password for us based on the
// user's input.
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword will check if passwords are matched.
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}

func GenerateUserId() string {
	return uuid.New().String()
}
