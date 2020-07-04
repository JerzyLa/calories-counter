package calories_datastore

import (
	"calories-counter/models"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const apiURL = "https://trackapi.nutritionix.com/v2/search/instant?query="
const kcalAttrID = 208

// NutritionixAPI implements models.Calories_Datastore
type NutritionixApi struct {
	id  string
	key string
}

type Nutrient struct {
	Value  float64 `json:"value"`
	AttrId int     `json:"attr_id"`
}

type apiResponse struct {
	Common []struct {
		FoodName      string     `json:"food_name"`
		FullNutrients []Nutrient `json:"full_nutrients"`
	} `json:"common"`

	Branded []struct {
		FoodName      string     `json:"food_name"`
		FullNutrients []Nutrient `json:"full_nutrients"`
	} `json:"branded"`
}

func NewNutritionixApi(appId, apiKey string) models.CaloriesDatastore {
	return &NutritionixApi{
		id:  appId,
		key: apiKey,
	}
}

func (a *NutritionixApi) GetCalories(mealName string) (*int, error) {
	log.Info("nutritionix api called")
	mealName = strings.ReplaceAll(mealName, " ", "%20")
	client := http.Client{}
	request, err := http.NewRequest("GET", apiURL+mealName+"&self=false&detailed=true", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("x-app-id", a.id)
	request.Header.Set("x-app-key", a.key)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result apiResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	for _, e := range result.Common {
		for _, attr := range e.FullNutrients {
			if attr.AttrId == kcalAttrID {
				res := int(attr.Value)
				return &res, nil
			}
		}

	}
	for _, e := range result.Branded {
		for _, attr := range e.FullNutrients {
			if attr.AttrId == kcalAttrID {
				res := int(attr.Value)
				return &res, nil
			}
		}
	}

	return nil, models.ErrMealCaloriesNotFound
}
