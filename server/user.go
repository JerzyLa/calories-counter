package server

import (
	"calories-counter/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	caller := c.MustGet("caller").(models.User)
	var body UserPostBody
	if err := c.ShouldBindJSON(&body); err != nil {
		handleErrorResponse(c, ErrInvalidJSON)
		return
	}

	if err := body.Validate(); err != nil {
		handleErrorResponse(c, err)
		return
	}

	// user manager can create only standard users
	if caller.RoleID == models.UserManagerRole && body.RoleID != models.UserRole {
		handleErrorResponse(c, ErrInsufficientPermissions)
		return
	}

	newUser, err := userRepo.SaveUser(caller.AccountID, body.Username, body.Password, body.RoleID)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func GetUsers(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	caller := c.MustGet("caller").(models.User)
	page, perPage, filter := PageParams(c)
	users, err := userRepo.GetUsers(caller.AccountID, page, perPage, filter)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	result := struct {
		models.UserSlice
		Links []Link `json:"links"`
	}{
		UserSlice: users,
		Links:     CreateLinks(c, users.Total, page, perPage),
	}

	c.PureJSON(http.StatusOK, result)
}

func GetUser(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	caller := c.MustGet("caller").(models.User)

	user, err := userRepo.GetUserById(caller.AccountID, c.Param("user_id"))
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	caller := c.MustGet("caller").(models.User)
	var body UserPutBody
	if err := c.ShouldBindJSON(&body); err != nil {
		handleErrorResponse(c, ErrInvalidJSON)
		return
	}

	if err := body.Validate(); err != nil {
		handleErrorResponse(c, err)
		return
	}

	user, err := userRepo.GetUserById(caller.AccountID, c.Param("user_id"))
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	// user manager cannot change role or modify non standard users
	if caller.RoleID == models.UserManagerRole && (user.RoleID != models.UserRole || body.RoleID != nil) {
		handleErrorResponse(c, ErrInsufficientPermissions)
		return
	}

	user.Username = body.Username
	if body.RoleID != nil {
		user.RoleID = *body.RoleID
	}
	newUser, err := userRepo.UpdateUser(*user)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func DeleteUser(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	caller := c.MustGet("caller").(models.User)

	user, err := userRepo.GetUserById(caller.AccountID, c.Param("user_id"))
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	// user manager cannot delete non standard users
	if caller.RoleID == models.UserManagerRole && user.RoleID != models.UserRole {
		handleErrorResponse(c, ErrInsufficientPermissions)
		return
	}

	err = userRepo.DeleteUser(caller.AccountID, c.Param("user_id"))
	if err != nil {
		handleErrorResponse(c, err)
	}

	c.Status(http.StatusNoContent)
}
