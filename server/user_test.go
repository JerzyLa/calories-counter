package server

import (
	"calories-counter/adapters/user_datastore"
	"calories-counter/models"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	testUser := &models.User{
		ID:        "1",
		AccountID: "1",
		Username:  "test",
		RoleID:    1,
	}
	jsonTestUser, _ := json.Marshal(testUser)

	testCases := []testCase{
		// error tests
		{
			name:          "UserNotFound",
			expectedCode:  http.StatusNotFound,
			expectedError: models.ErrUserNotFound,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(nil, models.ErrUserNotFound)
			},
		},

		// success tests
		{
			name:         "UserReturnedSuccessfully",
			expectedCode: http.StatusOK,
			expectedBody: string(jsonTestUser),
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(testUser, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "GET", "/v1/users/1", tc)
		})
	}
}

func TestGetUsers(t *testing.T) {
	testCases := []testCase{
		// success tests
		{
			name:         "UserReturnedSuccessfully",
			expectedCode: http.StatusOK,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUsers("", 0, 10, "").Return(models.UserSlice{
					Items: []models.User{
						{AccountID: "1", RoleID: models.AdminRole, Username: "admin1"},
						{AccountID: "1", RoleID: models.UserRole, Username: "user1"},
					},
					Total: 12,
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "GET", "/v1/users", tc)
		})
	}
}

func TestCreateUser(t *testing.T) {
	testCases := []testCase{
		// error tests
		{
			name:          "InvalidBody",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidJSON,
		},
		{
			name:          "MissingPassword",
			body:          `{"username":"testuser", "password":""}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrMissingPassword,
		},
		{
			name:          "InsufficientPermission",
			body:          `{"username":"testuser", "password": "Xyz123", "role_id": 2}`,
			caller:        models.User{RoleID: models.UserManagerRole},
			expectedCode:  http.StatusForbidden,
			expectedError: ErrInsufficientPermissions,
		},

		// success tests
		{
			name:         "UserManagerCreatedUser",
			body:         `{"username":"testuser", "password": "Xyz123", "role_id": 0}`,
			caller:       models.User{RoleID: models.UserManagerRole},
			expectedCode: http.StatusCreated,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().SaveUser("", "testuser", "Xyz123", 0)
			},
		},
		{
			name:         "AdminCreatedUserManager",
			body:         `{"username":"usermanager", "password": "Xyz123", "role_id": 1}`,
			caller:       models.User{RoleID: models.AdminRole},
			expectedCode: http.StatusCreated,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().SaveUser("", "usermanager", "Xyz123", 1)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "POST", "/v1/users", tc)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	testCases := []testCase{
		// error tests
		{
			name:          "InvalidBody",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidJSON,
		},
		{
			name:          "MissingUsername",
			body:          `{"username":""}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrMissingUsername,
		},
		{
			name:          "InvalidRoleID",
			body:          `{"username":"testuser", "role_id":4}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidRoleID,
		},
		{
			name:          "UserNotFound",
			body:          `{"username":"testuser"}`,
			expectedCode:  http.StatusNotFound,
			expectedError: models.ErrUserNotFound,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(nil, models.ErrUserNotFound)
			},
		},
		{
			name:          "InsufficientPermissionForUserManagerToChangeRole",
			body:          `{"username":"testuser", "role_id":1}`,
			caller:        models.User{RoleID: models.UserManagerRole},
			expectedCode:  http.StatusForbidden,
			expectedError: ErrInsufficientPermissions,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(&models.User{}, nil)
			},
		},
		{
			name:          "InsufficientPermissionForUserManagerToMakeChangeForOtherUserManager",
			body:          `{"username":"testuser"}`,
			caller:        models.User{RoleID: models.UserManagerRole},
			expectedCode:  http.StatusForbidden,
			expectedError: ErrInsufficientPermissions,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(&models.User{RoleID: models.UserManagerRole}, nil)
			},
		},

		// success tests
		{
			name:         "UserManagerUpdatedUser",
			body:         `{"username":"testuser"}`,
			caller:       models.User{RoleID: models.UserManagerRole},
			expectedCode: http.StatusOK,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(&models.User{}, nil)
				m.EXPECT().UpdateUser(models.User{Username: "testuser"})
			},
		},
		{
			name:         "AdminUpdatedUserManager",
			body:         `{"username":"testuser", "role_id":2}`,
			caller:       models.User{RoleID: models.AdminRole},
			expectedCode: http.StatusOK,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(&models.User{RoleID: models.UserManagerRole}, nil)
				m.EXPECT().UpdateUser(models.User{Username: "testuser", RoleID: 2})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "PUT", "/v1/users/1", tc)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []testCase{
		// error tests
		{
			name:          "UserNotFound",
			expectedCode:  http.StatusNotFound,
			expectedError: models.ErrUserNotFound,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(nil, models.ErrUserNotFound)
			},
		},
		{
			name:          "InsufficientPermissionForUserManagerToDeleteUserManager",
			caller:        models.User{RoleID: models.UserManagerRole},
			expectedCode:  http.StatusForbidden,
			expectedError: ErrInsufficientPermissions,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(&models.User{RoleID: models.UserManagerRole}, nil)
			},
		},

		// success tests
		{
			name:         "UserDeletedByUserManager",
			caller:       models.User{RoleID: models.UserManagerRole},
			expectedCode: http.StatusNoContent,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(&models.User{RoleID: models.UserRole}, nil)
				m.EXPECT().DeleteUser("", "1")
			},
		},
		{
			name:         "UserManagerDeletedByAdmin",
			caller:       models.User{RoleID: models.AdminRole},
			expectedCode: http.StatusNoContent,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetUserById("", "1").Return(&models.User{RoleID: models.UserManagerRole}, nil)
				m.EXPECT().DeleteUser("", "1")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "DELETE", "/v1/users/1", tc)
		})
	}
}
