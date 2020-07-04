package server

import (
	"calories-counter/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateMeal(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	caloriesDatastore := c.MustGet("caloriesDatastore").(models.CaloriesDatastore)
	user := c.MustGet("caller").(models.User)
	if _, ok := c.Get("user"); ok {
		user = c.MustGet("user").(models.User)
	}

	var body MealPostBody
	if err := c.ShouldBindJSON(&body); err != nil {
		handleErrorResponse(c, ErrInvalidJSON)
		return
	}

	if err := body.Validate(); err != nil {
		handleErrorResponse(c, err)
		return
	}

	if body.Calories == nil {
		calories, err := caloriesDatastore.GetCalories(body.Name)
		if err != nil {
			log.Printf("couldn not get calories for meal: %s", body.Name)
		}
		if calories == nil {
			v := 0
			body.Calories = &v
		} else {
			body.Calories = calories
		}
	}

	newMeal, err := userRepo.SaveMeal(user.ID, models.Meal{
		Date:     body.Date.String(),
		Time:     body.Time.String(),
		Name:     body.Name,
		Calories: *body.Calories,
	})
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, newMeal)
}

func GetMeals(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	user := c.MustGet("caller").(models.User)
	if _, ok := c.Get("user"); ok {
		user = c.MustGet("user").(models.User)
	}

	page, perPage, filter := PageParams(c)
	meals, err := userRepo.GetMeals(user.ID, page, perPage, filter)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	result := struct {
		models.MealSlice
		Links []Link `json:"links"`
	}{
		MealSlice: meals,
		Links:     CreateLinks(c, meals.Total, page, perPage),
	}

	c.PureJSON(http.StatusOK, result)
}

func GetMeal(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	user := c.MustGet("caller").(models.User)
	if _, ok := c.Get("user"); ok {
		user = c.MustGet("user").(models.User)
	}

	meal, err := userRepo.GetMeal(user.ID, c.Param("meal_id"))
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, meal)
}

func UpdateMeal(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	user := c.MustGet("caller").(models.User)
	if _, ok := c.Get("user"); ok {
		user = c.MustGet("user").(models.User)
	}

	var body MealPutBody
	if err := c.ShouldBindJSON(&body); err != nil {
		handleErrorResponse(c, ErrInvalidJSON)
		return
	}

	if err := body.Validate(); err != nil {
		handleErrorResponse(c, err)
		return
	}

	meal, err := userRepo.GetMeal(user.ID, c.Param("meal_id"))
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	if body.Name != "" {
		meal.Name = body.Name
	}
	if body.Time != nil {
		meal.Time = body.Time.String()
	}
	if body.Date != nil {
		meal.Date = body.Date.String()
	}
	if body.Calories != nil {
		meal.Calories = *body.Calories
	}
	updateMeal, err := userRepo.UpdateMeal(user.ID, *meal)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, updateMeal)
}

func DeleteMeal(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	user := c.MustGet("caller").(models.User)
	if _, ok := c.Get("user"); ok {
		user = c.MustGet("user").(models.User)
	}

	err := userRepo.DeleteMeal(user.ID, c.Param("meal_id"))
	if err != nil {
		handleErrorResponse(c, err)
	}

	c.Status(http.StatusNoContent)
}

func UpdateSettings(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	user := c.MustGet("caller").(models.User)
	if _, ok := c.Get("user"); ok {
		user = c.MustGet("user").(models.User)
	}
	var body SettingsPutBody
	if err := c.ShouldBindJSON(&body); err != nil {
		handleErrorResponse(c, ErrInvalidJSON)
		return
	}
	if err := body.Validate(); err != nil {
		handleErrorResponse(c, err)
		return
	}

	setting, err := userRepo.UpdateSettings(user.ID, models.Settings{ExpectedDailyCalories: *body.ExpectedDailyCalories})
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, setting)
}

func GetSettings(c *gin.Context) {
	userRepo := c.MustGet("userDatastore").(models.UserDatastore)
	user := c.MustGet("caller").(models.User)
	if _, ok := c.Get("user"); ok {
		user = c.MustGet("user").(models.User)
	}

	setting, err := userRepo.GetSettings(user.ID)
	if err != nil {
		handleErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, setting)
}
