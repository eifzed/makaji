package users

import (
	"time"
)

type UserDetail struct {
	UserID   string `json:"-" xorm:"user_id pk autoincr"`
	Email    string `json:"email" xorm:"email"`
	Username string `json:"username" xorm:"username"`
	Password string `json:"password" xorm:"-"`
}

type UserRegistration struct {
	Username string `json:"username" xorm:"username"`
	Email    string `json:"email" xorm:"email"`
	Password string `json:"password" xorm:"-"`
}

type UserLogin struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type UserAuth struct {
	Token      string
	ValidUntil time.Time
}

type UserRole struct {
	ID   int64  `yaml:"id" json:"id" xorm:"role_id"`
	Name string `yaml:"name" json:"name" xorm:"role_name"`
}
