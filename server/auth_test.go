package server

import (
	"calories-counter/adapters/user_datastore"
	"calories-counter/common"
	"calories-counter/models"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
)

func TestSignUp(t *testing.T) {
	testUser := &models.User{
		ID:        "1",
		AccountID: "2",
		Username:  "test",
		RoleID:    1,
	}
	jsonTestUser, _ := json.Marshal(testUser)

	testCases := []testCase{
		// error tests
		{
			name:          "MissingUsername",
			postForm:      true,
			body:          postForm("", ""),
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrMissingUsername,
		},
		{
			name:          "MissingPassword",
			postForm:      true,
			body:          postForm("username", ""),
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrMissingPassword,
		},
		{
			name:          "InvalidUsername",
			postForm:      true,
			body:          postForm("???", "Xyz123"),
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidUsername,
		},
		{
			name:          "InvalidPassword",
			postForm:      true,
			body:          postForm("xyz@xyz.xyz", "xyzxyz"),
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidPassword,
		},
		{
			name:          "ErrWhenGetRootUser",
			postForm:      true,
			body:          postForm("xyz@xyz.xyz", "Xyz123"),
			expectedCode:  http.StatusInternalServerError,
			expectedError: ErrInternalServerError,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().SaveRootUser("xyz@xyz.xyz", "Xyz123").Return(nil, errors.New("err")).Times(1)
			},
		},
		{
			name:          "AccountAlreadyExists",
			postForm:      true,
			body:          postForm("test@test.com", "Xyz123"),
			expectedCode:  http.StatusConflict,
			expectedError: models.ErrAccountAlreadyExists,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().SaveRootUser(gomock.Any(), gomock.Any()).Return(nil, models.ErrAccountAlreadyExists).Times(1)
			},
		},

		// success tests
		{
			name:         "AccountCreated",
			postForm:     true,
			body:         postForm("test@test.com", "Xyz123"),
			expectedBody: string(jsonTestUser),
			expectedCode: http.StatusCreated,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().SaveRootUser(gomock.Any(), gomock.Any()).Return(testUser, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "POST", "/v1/signup", tc)
		})
	}
}

func TestSignIn(t *testing.T) {
	testCases := []testCase{
		// errors
		{
			name:          "MissingUsername",
			postForm:      true,
			body:          postForm("", ""),
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrMissingUsername,
		},
		{
			name:          "MissingPassword",
			postForm:      true,
			body:          postForm("username", ""),
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrMissingPassword,
		},
		{
			name:          "UnauthorizedWrongPassword",
			postForm:      true,
			body:          postForm("user", "Xyz123"),
			expectedCode:  http.StatusUnauthorized,
			expectedError: ErrUnauthorized,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserPassword("1", "user").Return(common.String("pass"), nil)
			},
		},

		// success tests
		{
			name:         "LoggedSuccessfully",
			postForm:     true,
			body:         postForm("user", "Xyz123"),
			expectedCode: http.StatusCreated,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserPassword("1", "user").Return(common.String("Xyz123"), nil)
				m.EXPECT().GetUser("1", "user").Return(&models.User{}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "POST", "/v1/account/1/signin", tc)
		})
	}
}
