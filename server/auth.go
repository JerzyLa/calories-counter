package server

import (
	"calories-counter/common"
	"calories-counter/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SignUp(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if err := ValidateUsername(username); err != nil {
		handleErrorResponse(c, err)
		return
	}
	if err := ValidatePassword(password); err != nil {
		handleErrorResponse(c, err)
		return
	}

	newAccount, err := userRepo.SaveRootUser(username, password)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, newAccount)
}

func SignIn(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	secretKey := c.MustGet("secretKey").(string)

	accountID := c.Param("account_id")
	username := c.PostForm("username")
	password := c.PostForm("password")
	if accountID == "" {
		handleErrorResponse(c, ErrMissingAccountID)
		return
	}
	if username == "" {
		handleErrorResponse(c, ErrMissingUsername)
		return
	}
	if password == "" {
		handleErrorResponse(c, ErrMissingPassword)
		return
	}

	pass, err := userRepo.GetUserPassword(accountID, username)
	if err != nil {
		handleErrorResponse(c, err)
		return
	} else if pass == nil || password != *pass {
		handleErrorResponse(c, ErrUnauthorized)
		return
	}

	user, err := userRepo.GetUser(accountID, username)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	claims := common.JWTClaims{
		UUID:      user.ID,
		AccountID: user.AccountID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": tokenString})
}
