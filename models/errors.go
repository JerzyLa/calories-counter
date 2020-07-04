package models

import (
	"calories-counter/common"
	"errors"
	"net/http"
)

var (
	ErrAccountAlreadyExists = common.ApiErr{
		Code: http.StatusConflict,
		Err:  errors.New("account already exist"),
	}

	ErrUserAlreadyExists = common.ApiErr{
		Code: http.StatusConflict,
		Err:  errors.New("user already exist"),
	}

	ErrUserNotFound = common.ApiErr{
		Code: http.StatusNotFound,
		Err:  errors.New("user not found"),
	}

	ErrMealNotFound = common.ApiErr{
		Code: http.StatusNotFound,
		Err:  errors.New("meal not found"),
	}

	ErrInvalidFilter = common.ApiErr{
		Code: http.StatusBadRequest,
		Err:  errors.New("invalid filter"),
	}

	ErrInvalidQuery = common.ApiErr{
		Code: http.StatusBadRequest,
		Err:  errors.New("invalid query"),
	}

	ErrMealCaloriesNotFound = common.ApiErr{
		Code: http.StatusNotFound,
		Err:  errors.New("meal calories not found"),
	}
)
