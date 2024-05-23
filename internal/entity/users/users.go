package users

import (
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/eifzed/joona/lib/common/commonerr"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDetail struct {
	UserID         primitive.ObjectID `json:"_id" bson:"_id"`
	Email          string             `json:"email" bson:"email"`
	Username       string             `json:"username" bson:"username"`
	Password       string             `json:"password" bson:"password"`
	FullName       string             `json:"full_name" bson:"full_name"`
	ProfilePicture string             `json:"profile_picture" bson:"profile_picture"`
	Nationality    string             `json:"nationality" bson:"nationality"`
	Occupation     string             `json:"occupation" bson:"occupation"`
	Bio            string             `json:"bio" bson:"bio"`
}

type UserProfile struct {
	UserID         primitive.ObjectID `json:"_id" bson:"_id"`
	Email          string             `json:"email" bson:"email"`
	Username       string             `json:"username" bson:"username"`
	FullName       string             `json:"full_name" bson:"full_name"`
	ProfilePicture string             `json:"profile_picture" bson:"profile_picture"`
	Nationality    string             `json:"nationality" bson:"nationality"`
	Occupation     string             `json:"occupation" bson:"occupation"`
	Bio            string             `json:"bio" bson:"bio"`
}

type UserRegistration struct {
	Username string `json:"username" conform:"trim"`
	Email    string `json:"email" conform:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name" conform:"name"`
}

func (s *UserRegistration) ValidateInput() error {
	if s.Username == "" || s.Password == "" || s.Email == "" {
		return commonerr.ErrorBadRequest("invalid input", "email, username, and email cannot be empty")
	}
	if len(s.FullName) < 4 {
		return commonerr.ErrorBadRequest("invalid input", "full name must be at least 4 characters")
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
	var number, upper, special, sevenOrMore bool
	for _, c := range s.Password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
		default:
			return false
		}
	}
	sevenOrMore = len(s.Password) >= 7
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
	Token      string    `json:"token"`
	ValidUntil time.Time `json:"valid_until"`
}

type UserRole struct {
	ID   int64  `yaml:"id" json:"id" xorm:"role_id"`
	Name string `yaml:"name" json:"name" xorm:"role_name"`
}

type UserItem struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	ProfilePicture string `json:"profile_picture"`
	Nationality    string `json:"nationality"`
	Bio            string `json:"bio"`
}
