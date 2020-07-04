package server

import (
	"calories-counter/adapters/calories_datastore"
	"calories-counter/adapters/user_datastore"
	"calories-counter/models"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type testCase struct {
	name              string
	postForm          bool
	body              string
	expectedBody      string
	expectedCode      int
	expectedError     error
	caller            models.User
	user              models.User
	setupMockUser     func(m *user_datastore.MockUserDatastore)
	setupMockCalories func(m *calories_datastore.MockCaloriesDatastore)
}

func runTest(t *testing.T, method, path string, tc testCase) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockUD := user_datastore.NewMockUserDatastore(controller)
	if tc.setupMockUser != nil {
		tc.setupMockUser(mockUD)
	}
	mockCD := calories_datastore.NewMockCaloriesDatastore(controller)
	if tc.setupMockCalories != nil {
		tc.setupMockCalories(mockCD)
	}

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(SetVars(map[string]interface{}{
		"caller": tc.caller,
		"user":   tc.user,
	}))
	setupTestRouter(r, "test_secret_key", mockUD, mockCD)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(tc.body))
	if tc.postForm {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	r.ServeHTTP(w, req)

	if w.Code != tc.expectedCode {
		t.Errorf("Expected status code to be %d but was %d", tc.expectedCode, w.Code)
	}
	if tc.expectedError != nil {
		if !strings.Contains(w.Body.String(), tc.expectedError.Error()) {
			t.Errorf("Expected error message to have `%s`, msg: `%s`", tc.expectedError.Error(), w.Body.String())
		}
	} else if tc.expectedBody != "" {
		if w.Body.String() != tc.expectedBody {
			t.Errorf("Expected body to be `%s` but was `%s`", tc.expectedBody, w.Body.String())
		}
	}
}

func postForm(username, password string) string {
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	return form.Encode()
}

func setupTestRouter(r *gin.Engine, secretKey string, datastore models.UserDatastore, caloriesDatastore models.CaloriesDatastore) {
	r.Use(SetVars(map[string]interface{}{
		"userDatastore":     datastore,
		"secretKey":         secretKey,
		"caloriesDatastore": caloriesDatastore,
	}))
	r.POST("/v1/signup", SignUp)
	r.POST("/v1/account/:account_id/signin", SignIn)

	r.POST("/v1/users", CreateUser)
	r.GET("/v1/users", GetUsers)
	r.GET("/v1/users/:user_id", GetUser)
	r.PUT("/v1/users/:user_id", UpdateUser)
	r.DELETE("/v1/users/:user_id", DeleteUser)

	r.POST("/v1/meals", CreateMeal)
	r.GET("/v1/meals", GetMeals)
	r.GET("/v1/meals/:meal_id", GetMeal)
	r.PUT("/v1/meals/:meal_id", UpdateMeal)
	r.DELETE("/v1/meals/:meal_id", DeleteMeal)

	r.PUT("/v1/settings", UpdateSettings)
	r.GET("/v1/settings", GetSettings)
}
