package server

import (
	"calories-counter/models"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, secretKey string, userDatastore models.UserDatastore, caloriesDatastore models.CaloriesDatastore) {
	r.Use(SetVars(map[string]interface{}{
		"userDatastore":     userDatastore,
		"caloriesDatastore": caloriesDatastore,
		"secretKey":         secretKey,
	}))
	r.POST("/v1/signup", SignUp)
	r.POST("/v1/account/:account_id/signin", SignIn)

	authorized := r.Group("/v1")
	authorized.Use(AuthVerify())
	{
		users := authorized.Group("/users")
		users.Use(RoleAccessVerify(models.AdminRole, models.UserManagerRole, models.OwnerRole))
		{
			users.POST("/", CreateUser)
			users.GET("/", GetUsers)
			users.GET("/:user_id", GetUser)
			users.PUT("/:user_id", UpdateUser)
			users.DELETE("/:user_id", DeleteUser)
		}

		meals := authorized.Group("/meals")
		meals.Use(RoleAccessVerify(models.UserRole))
		{
			meals.POST("/", CreateMeal)
			meals.GET("/", GetMeals)
			meals.GET("/:meal_id", GetMeal)
			meals.PUT("/:meal_id", UpdateMeal)
			meals.DELETE("/:meal_id", DeleteMeal)
		}
		settings := authorized.Group("/settings")
		settings.Use(RoleAccessVerify(models.UserRole))
		{
			settings.PUT("/", UpdateSettings)
			settings.GET("/", GetSettings)
		}

		adminMeals := authorized.Group("/users/:user_id/meals")
		adminMeals.Use(RoleAccessVerify(models.AdminRole), UserVerify())
		{
			adminMeals.POST("/", CreateMeal)
			adminMeals.GET("/", GetMeals)
			adminMeals.GET("/:meal_id", GetMeal)
			adminMeals.PUT("/:meal_id", UpdateMeal)
			adminMeals.DELETE("/:meal_id", DeleteMeal)
		}
		adminSettings := authorized.Group("/users/:user_id/settings")
		adminSettings.Use(RoleAccessVerify(models.AdminRole), UserVerify())
		{
			adminSettings.PUT("/", UpdateSettings)
			adminSettings.GET("/", GetSettings)
		}
	}
}
