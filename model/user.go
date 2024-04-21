package model

import (
	"errors"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"unicode"
)

type RegisterUserReq struct {
	Phone    string `json:"phone" validate:"required,min=10,max=13,phone_prefix=+62"`
	Name     string `json:"name" validate:"required,min=3,max=60"`
	Password string `json:"password" validate:"required,min=6,max=64,password"`
}

func validatePhonePrefix(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return strings.HasPrefix(phone, "+62")
}

func registerCustomValidators(validate *validator.Validate) {
	err := validate.RegisterValidation("phone_prefix", validatePhonePrefix)
	if err != nil {
		return
	}
	err = validate.RegisterValidation("password", validatePassword)
	if err != nil {
		return
	}
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasSpecialChar bool
		hasCapital     bool
		hasNumeric     bool
	)

	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasNumeric = true
		case unicode.IsUpper(char):
			hasCapital = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	return hasSpecialChar && hasCapital && hasNumeric
}

func (r *RegisterUserReq) Validate() error {
	validate := validator.New()
	registerCustomValidators(validate)
	return validate.Struct(r)
}

func (r *RegisterUserReq) ToDAO() (repository.RegisterUser, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return repository.RegisterUser{}, err
	}
	return repository.RegisterUser{
		ID:       uuid.New().String(),
		Phone:    r.Phone,
		Name:     r.Name,
		Password: string(hashedPassword),
	}, nil
}

type User struct {
	UserID    string    `json:"userID" db:"id"`
	Phone     string    `json:"phone" db:"phone"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (u *User) CheckLogin(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return errors.New("password does not match")
	}
	return nil
}

func (u *User) ToProfileResp() GetProfileResp {
	return GetProfileResp{
		Name:  u.Name,
		Phone: u.Phone,
	}
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func FromRepoUser(repoUser repository.User) User {
	return User{
		UserID:    repoUser.UserID,
		Phone:     repoUser.Phone,
		Name:      repoUser.Name,
		Password:  repoUser.Password,
		CreatedAt: repoUser.CreatedAt,
		UpdatedAt: repoUser.UpdatedAt.Time,
	}
}

type GetProfileResp struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type UpdateUserReq struct {
	Phone string `json:"phone,omitempty"`
	Name  string `json:"name,omitempty"`
}

type ValidateUpdateUserPhone struct {
	Phone string `json:"phone" validate:"required,min=10,max=13,phone_prefix=+62"`
}

type ValidateUpdateUserName struct {
	Name string `json:"name" validate:"required,min=3,max=60"`
}

func (u *UpdateUserReq) ToDAO() repository.UpdateUser {
	return repository.UpdateUser{
		Phone: u.Phone,
		Name:  u.Name,
	}
}

func (u *UpdateUserReq) Validate() error {
	validate := validator.New()
	registerCustomValidators(validate)
	var eitherExists bool
	if len(strings.TrimSpace(u.Phone)) > 1 {
		eitherExists = true
		phone := ValidateUpdateUserPhone{Phone: u.Phone}
		err := validate.Struct(phone)
		if err != nil {
			return err
		}
	}
	if len(strings.TrimSpace(u.Name)) > 1 {
		eitherExists = true
		name := ValidateUpdateUserName{Name: u.Name}
		err := validate.Struct(name)
		if err != nil {
			return err
		}
	}

	if !eitherExists {
		return errors.New("phone or name should exists")
	}
	return nil
}
