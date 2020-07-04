package server

import (
	"calories-counter/common"
	"errors"
	"fmt"
	"net/http"
)

const (
	MaxUsernameLength = 50
	MinUsernameLength = 5
	MaxPasswordLength = 50
	MinPasswordLength = 5
	MaxNameLength     = 50
)

var (
	ErrInternalServerError = common.ApiErr{
		Code: http.StatusInternalServerError,
		Err:  errors.New(http.StatusText(http.StatusInternalServerError)),
	}

	ErrUnauthorized = common.ApiErr{
		Code: http.StatusUnauthorized,
		Err:  errors.New(http.StatusText(http.StatusUnauthorized)),
	}

	ErrInsufficientPermissions = common.ApiErr{
		Code: http.StatusForbidden,
		Err:  errors.New("insufficient permissions"),
	}

	ErrInvalidRoleID = common.ApiErr{
		Code: http.StatusBadRequest,
		Err:  errors.New("invalid roleID"),
	}

	ErrMissingBearerToken = common.ApiErr{
		Code: http.StatusUnauthorized,
		Err:  errors.New("missing Bearer token"),
	}

	ErrUnexpectedSigningMethod = common.ApiErr{
		Code: http.StatusUnauthorized,
		Err:  errors.New("unexpected signing method"),
	}

	ErrMissingName = common.ApiErr{
		Code: http.StatusBadRequest,
		Err:  errors.New("missing name"),
	}

	ErrMissingDate = common.ApiErr{
		Code: http.StatusBadRequest,
		Err:  errors.New("missing date"),
	}

	ErrMissingTime = common.ApiErr{
		Code: http.StatusBadRequest,
		Err:  errors.New("missing time"),
	}

	ErrInvalidNameLength = common.ApiErr{
		Code: http.StatusBadRequest,
		Err:  fmt.Errorf("invalid name length, name can not be larger than %d character", MaxNameLength),
	}
)

var ErrMissingPassword = common.ApiErr{
	Code: http.StatusBadRequest,
	Err:  errors.New("missing password"),
}

var ErrMissingUsername = common.ApiErr{
	Code: http.StatusBadRequest,
	Err:  errors.New("missing username"),
}

var ErrMissingAccountID = common.ApiErr{
	Code: http.StatusBadRequest,
	Err:  errors.New("missing accountID"),
}

var ErrInvalidPasswordLength = common.ApiErr{
	Code: http.StatusBadRequest,
	Err: fmt.Errorf("invalid password length, "+
		"password can not be shorter than %d and larger than %d characters", MinPasswordLength, MaxPasswordLength),
}

var ErrInvalidUsernameLength = common.ApiErr{
	Code: http.StatusBadRequest,
	Err: fmt.Errorf("invalid username length, "+
		"username can not be shorter than %d and larger than %d character", MinUsernameLength, MaxUsernameLength),
}

var ErrInvalidUsername = common.ApiErr{
	Code: http.StatusBadRequest,
	Err:  errors.New("invalid username"),
}

var ErrInvalidPassword = common.ApiErr{
	Code: http.StatusBadRequest,
	Err:  errors.New("invalid password, password has to contain at one number and one uppercase"),
}

var ErrInvalidJSON = common.ApiErr{
	Code: http.StatusBadRequest,
	Err:  errors.New("invalid JSON"),
}
