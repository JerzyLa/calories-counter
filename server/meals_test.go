package server

import (
	"calories-counter/adapters/user_datastore"
	"calories-counter/models"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
)

func TestCreateMeal(t *testing.T) {
	meal := models.Meal{
		ID:              "1",
		Date:            "2020-01-01",
		Time:            "11:11:11",
		Name:            "test",
		Calories:        100,
		CaloriesDeficit: false,
	}

	testCases := []testCase{
		// error tests
		{
			name:          "InvalidBody",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidJSON,
		},
		{
			name:          "MissingDate",
			body:          `{"name":"chicken"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrMissingDate,
		},

		// success tests
		{
			name:         "SavedMeal",
			body:         `{"name":"chicken", "date":"2020-01-01", "time":"10:10:10", "calories":100}`,
			user:         models.User{ID: "1"},
			expectedCode: http.StatusCreated,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().SaveMeal("1", gomock.Any()).Return(&meal, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "POST", "/v1/meals", tc)
		})
	}
}

func TestGetMeal(t *testing.T) {
	testMeal := &models.Meal{
		ID:              "1",
		Date:            "2020-01-01",
		Time:            "11:11:11",
		Name:            "meal",
		Calories:        10,
		CaloriesDeficit: false,
	}
	jsonTestMeal, _ := json.Marshal(testMeal)

	testCases := []testCase{
		// error tests
		{
			name:          "MealNotFound",
			expectedCode:  http.StatusNotFound,
			expectedError: models.ErrMealNotFound,
			user:          models.User{ID: "1"},
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetMeal("1", "1").Return(nil, models.ErrMealNotFound)
			},
		},

		// success tests
		{
			name:         "MealReturnedSuccessfully",
			expectedCode: http.StatusOK,
			expectedBody: string(jsonTestMeal),
			user:         models.User{ID: "1"},
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetMeal("1", "1").Return(testMeal, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "GET", "/v1/meals/1", tc)
		})
	}
}

func TestUpdateMeal(t *testing.T) {
	testMeal := &models.Meal{
		ID:              "1",
		Date:            "2020-01-01",
		Time:            "11:11:11",
		Name:            "meal",
		Calories:        10,
		CaloriesDeficit: false,
	}

	testCases := []testCase{
		// error tests
		{
			name:          "InvalidBody",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidJSON,
		},
		{
			name:          "MealNotFound",
			body:          `{"name":"chicken", "date":"2020-01-01", "time":"10:10:10", "calories":100}`,
			expectedCode:  http.StatusNotFound,
			expectedError: models.ErrMealNotFound,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetMeal("", "1").Return(nil, models.ErrMealNotFound)
			},
		},

		// success tests
		{
			name:         "UpdateMeal",
			body:         `{"name":"chicken", "calories":100}`,
			expectedCode: http.StatusOK,
			setupMockUser: func(m *user_datastore.MockUserDatastore) {
				m.EXPECT().GetMeal("", "1").Return(testMeal, nil)
				testMeal.Name = "chicken"
				testMeal.Calories = 100
				m.EXPECT().UpdateMeal("", *testMeal)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runTest(t, "PUT", "/v1/meals/1", tc)
		})
	}
}
