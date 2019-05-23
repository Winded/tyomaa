package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	"github.com/winded/tyomaa/backend/util"
	"github.com/winded/tyomaa/shared/api"
)

const (
	PASSWORD_HASH_COST = 14
)

type User struct {
	gorm.Model
	Name     string      `gorm:"NOT NULL"`
	Password string      `gorm:"NOT NULL"`
	Entries  []TimeEntry `gorm:"association_autoupdate:false"`
}

func (this *User) BeforeSave() error {
	if !util.ValidateNameIdentifier(this.Name) {
		return errors.New("Name must only contain alphabetic characters, numbers and dashes")
	}

	return nil
}

// CheckPassword verifies that the given raw password matches the user's hashed password
func (this *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(this.Password), []byte(password))
}

// SetPassword hashes the given password and adds the result as the user's password
func (this *User) SetPassword(password string) error {
	data, err := bcrypt.GenerateFromPassword([]byte(password), PASSWORD_HASH_COST)
	if err != nil {
		return err
	}

	this.Password = string(data)
	return nil
}

func (this *User) ToApiFormat() api.User {
	return api.User{
		ID:   this.ID,
		Name: this.Name,
	}
}

type TimeEntry struct {
	gorm.Model
	UserID  uint       `gorm:"NOT NULL"`
	Project string     `gorm:"NOT NULL"`
	Start   time.Time  `gorm:"NOT NULL"`
	End     *time.Time `gorm:"NULL"`
}

func (this *TimeEntry) BeforeSave() error {
	if !util.ValidateNameIdentifier(this.Project) {
		return errors.New("Project name must only contain alphabetic characters, numbers and dashes")
	}

	return nil
}

func (this *TimeEntry) ToApiFormat() api.TimeEntry {
	return api.TimeEntry{
		ID:      this.ID,
		Project: this.Project,
		Start:   this.Start,
		End:     this.End,
	}
}
