package users

import (
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/eifzed/joona/lib/common/commonerr"
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

func (s *UserRegistration) ValidateInput() error {
	if s.Username == "" || s.Password == "" || s.Email == "" {
		return commonerr.ErrorBadRequest("invalid input", "email, username, and email cannot be empty")
	}
	if _, err := mail.ParseAddress(s.Email); err != nil {
		return commonerr.ErrorBadRequest("invalid email", "invalid email format")
	}
	if !s.CheckIsValidUsername() {
		return commonerr.ErrorBadRequest("invalid username", "invalid username")
	}
	if !s.CheckIsValidPassword() {
		return commonerr.ErrorBadRequest("invalid password", "invalid password")
	}

	return nil
}

func (s *UserRegistration) CheckIsValidPassword() bool {
	letters := 0
	var number, upper, special, sevenOrMore bool
	for _, c := range s.Password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false
		}
	}
	sevenOrMore = letters >= 7
	return number && upper && special && sevenOrMore
}

func (s *UserRegistration) CheckIsValidUsername() bool {
	if len(s.Username) < 4 || strings.Contains(s.Username, " ") {
		return false
	}
	if s.Username[0] == '.' || s.Username[len(s.Username)-1] == '.' || s.Username[0] == '_' || s.Username[len(s.Username)-1] == '_' {
		return false
	}
	for _, chr := range s.Username {
		if (chr < 'a' || chr > 'z') && (chr < 'A' || chr > 'Z') && (chr != '_' && chr != '.') {
			return false
		}
	}
	return true
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
