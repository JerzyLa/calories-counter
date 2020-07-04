package test

import (
	"calories-counter/adapters/calories_datastore"
	"calories-counter/adapters/user_datastore"
	"calories-counter/models"
	"calories-counter/server"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

var (
	secretKey = os.Getenv("TOKEN_SECRET")
	dbSource  = os.Getenv("MYSQL_DB_SOURCE")
	appID     = os.Getenv("API_APP_ID")
	apiKey    = os.Getenv("API_KEY")
)

// this test uses original database and api instead of mocks
func TestUserCRUD(t *testing.T) {
	caloriesDatastore := calories_datastore.NewNutritionixApi(appID, apiKey)
	userDatastore, _ := user_datastore.NewMySQLStore(dbSource)
	defer func() { _ = userDatastore.Close() }()
	r := gin.Default()
	server.SetupRouter(r, secretKey, userDatastore, caloriesDatastore)

	ownerUsername := "userAdmin"
	pass := "userAdmin1"
	form := url.Values{}
	form.Set("username", ownerUsername)
	form.Set("password", pass)

	t.Log("1. Signup")
	req, _ := http.NewRequest("POST", "/v1/signup", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var userAdmin models.User
	_ = json.NewDecoder(w.Body).Decode(&userAdmin)

	t.Log("2. Signin")
	req, _ = http.NewRequest("POST", "/v1/account/"+userAdmin.AccountID+"/signin", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	token := struct {
		Token string `json:"token"`
	}{}
	_ = json.NewDecoder(w.Body).Decode(&token)

	t.Log("3. Create user")
	testUsername := "usertest"
	req, _ = http.NewRequest("POST", "/v1/users/", strings.NewReader(fmt.Sprintf(`{"username":"%s","password":"%s"}`, testUsername, pass)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var user models.User
	_ = json.NewDecoder(w.Body).Decode(&user)
	if !(user.AccountID == userAdmin.AccountID && user.Username == testUsername) {
		t.Errorf("create user failed")
	}

	t.Log("4. Update user")
	newTestUsername := "newTestUsername"
	req, _ = http.NewRequest("PUT", "/v1/users/"+user.ID, strings.NewReader(fmt.Sprintf(`{"username":"%s"}`, newTestUsername)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var updatedUser models.User
	_ = json.NewDecoder(w.Body).Decode(&updatedUser)
	if !(updatedUser.AccountID == user.AccountID && updatedUser.Username == newTestUsername) {
		t.Errorf("update user failed")
	}

	t.Log("5. Delete user")
	req, _ = http.NewRequest("DELETE", "/v1/users/"+user.ID, nil)
	req.Header.Set("Authorization", "Bearer "+token.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	t.Log("6. Delete myself")
	req, _ = http.NewRequest("DELETE", "/v1/users/"+userAdmin.ID, nil)
	req.Header.Set("Authorization", "Bearer "+token.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/v1/users/", nil)
	req.Header.Set("Authorization", "Bearer "+token.Token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("delete user failed")
	}
}
