package server

import (
	"calories-counter/common"
	"calories-counter/models"
	"unicode"
)

type UserPostBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id"`
}

type UserPutBody struct {
	Username string `json:"username"`
	RoleID   *int   `json:"role_id"`
}

type Validator interface {
	Validate() error
}

func (body *UserPostBody) Validate() error {
	if err := ValidateUsername(body.Username); err != nil {
		return err
	}
	if err := ValidatePassword(body.Password); err != nil {
		return err
	}
	if body.RoleID < models.UserRole || body.RoleID > models.AdminRole {
		return ErrInvalidRoleID
	}

	return nil
}

// Validate username length and chars
func ValidateUsername(username string) error {
	if username == "" {
		return ErrMissingUsername
	}
	if len(username) > MaxUsernameLength || len(username) < MinUsernameLength {
		return ErrInvalidUsernameLength
	}
	for _, c := range username {
		if !(unicode.IsLetter(c) || unicode.IsNumber(c) || unicode.IsPunct(c)) {
			return ErrInvalidUsername
		}
	}

	return nil
}

func (body *UserPutBody) Validate() error {
	if err := ValidateUsername(body.Username); err != nil {
		return err
	}
	if body.RoleID != nil && (*body.RoleID < models.UserRole || *body.RoleID > models.AdminRole) {
		return ErrInvalidRoleID
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return ErrMissingPassword
	}
	if len(password) > MaxPasswordLength || len(password) < MinPasswordLength {
		return ErrInvalidPasswordLength
	}
	var digits, uppercase bool
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			digits = true
		case unicode.IsUpper(c):
			uppercase = true
		case unicode.IsPunct(c) || unicode.IsLetter(c):
		default:
			return ErrInvalidPassword
		}
	}

	if !digits || !uppercase {
		return ErrInvalidPassword
	}

	return nil
}

type MealPostBody struct {
	Date     *common.Date `json:"date"`
	Time     *common.Time `json:"time"`
	Name     string       `json:"name"`
	Calories *int         `json:"calories"`
}

func (body *MealPostBody) Validate() error {
	if body.Name == "" {
		return ErrMissingName
	}
	if body.Date == nil {
		return ErrMissingDate
	}
	if body.Time == nil {
		return ErrMissingTime
	}
	if len(body.Name) > MaxNameLength {
		return ErrInvalidNameLength
	}
	return nil
}

type MealPutBody struct {
	MealPostBody
}

func (body *MealPutBody) Validate() error {
	if len(body.Name) > MaxNameLength {
		return ErrInvalidNameLength
	}
	if body.Name == "" && body.Date == nil && body.Time == nil && body.Calories == nil {
		return ErrInvalidJSON
	}
	return nil
}

type SettingsPutBody struct {
	ExpectedDailyCalories *int `json:"expected_daily_calories"`
}

func (body *SettingsPutBody) Validate() error {
	if body.ExpectedDailyCalories == nil || *body.ExpectedDailyCalories < 0 {
		return ErrInvalidJSON
	}
	return nil
}
