package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username    string `gorm:"varchar(50), unique"`
	Email       string `gorm:"varchar(50), unique"`
	Password    string `gorm:"vrchar(50)"`
	Firstname   string `gorm:"vrchar(50)"`
	Lastname    string `gorm:"vrchar(50)"`
	PhoneNumber string `gorm:"varchar(50), unique"`
	Gender      string `gorm:"vrchar(50)"`
}

func (gdb *GormDB) CreateNewUser(u *User) error {
	// Check duplicate username
	var count int64
	gdb.db.Model(&User{}).Where(&User{Username: u.Username}).Count(&count)
	if count > 0 {
		return errors.New("Username already exists")
	}
	// Check duplicate username
	gdb.db.Model(&User{}).Where(&User{Email: u.Email}).Count(&count)
	if count > 0 {
		return errors.New("Email already exists")
	}

	// Check duplicate PhoneNumber
	gdb.db.Model(&User{}).Where(&User{PhoneNumber: u.PhoneNumber}).Count(&count)
	if count > 0 {
		return errors.New("PhoneNumber already exists")
	}

	// Encrypt user's password
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	if err != nil {
		return err
	}
	u.Password = string(pwd)

	// Create the user
	if err := gdb.db.Create(u).Error; err != nil {
		return err
	}
	return nil
}

func (gdb *GormDB) GetUserByUsername(username string) (*User, error) {
	var user User
	err := gdb.db.Where(&User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
