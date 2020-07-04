package server

import (
	"calories-counter/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"strings"
)

// SetVars middleware function which pass variables (like secrets, db connectors) to handlers
func SetVars(m map[string]interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		for k, v := range m {
			c.Set(k, v)
		}
	}
}

// AuthVerify middleware function which verifies jwt tokens
func AuthVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRepo := c.MustGet("userDatastore").(models.UserDatastore)
		secretKey := c.MustGet("secretKey").(string)

		auth := strings.Split(c.GetHeader("Authorization"), " ")
		if !(len(auth) == 2 && auth[0] == "Bearer") {
			c.Abort()
			handleErrorResponse(c, ErrMissingBearerToken)
			return
		}

		tokenString := auth[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			c.Abort()
			log.Info(err)
			handleErrorResponse(c, ErrUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Info(claims["UUID"], claims["AccountID"], claims["iat"], claims["exp"])
			user, err := userRepo.GetUserById(claims["AccountID"].(string), claims["UUID"].(string))
			if err != nil {
				c.Abort()
				handleErrorResponse(c, ErrUnauthorized)
				return
			}

			// set caller
			c.Set("caller", *user)
		} else {
			c.Abort()
			handleErrorResponse(c, ErrUnauthorized)
			return
		}
	}
}

func RoleAccessVerify(roleIDs ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		caller := c.MustGet("caller").(models.User)

		if roleIDs != nil {
			var hasAccess bool
			for _, roleID := range roleIDs {
				if roleID == caller.RoleID {
					hasAccess = true
				}
			}
			if !hasAccess {
				c.Abort()
				handleErrorResponse(c, ErrInsufficientPermissions)
				return
			}
		}
	}
}

func UserVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		caller := c.MustGet("caller").(models.User)
		userRepo := c.MustGet("userDatastore").(models.UserDatastore)

		user, err := userRepo.GetUserById(caller.AccountID, c.Param("user_id"))
		if err != nil {
			c.Abort()
			handleErrorResponse(c, err)
			return
		}
		c.Set("user", *user)
	}
}
