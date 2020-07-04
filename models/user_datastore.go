package models

//go:generate mockgen -destination=../adapters/user_datastore/mock.go -package=user_datastore calories-counter/models UserDatastore

type User struct {
	ID        string `json:"id" db:"id"`
	AccountID string `json:"account_id" db:"account_id"`
	Username  string `json:"username" db:"username"`
	RoleID    int    `json:"role_id" db:"role_id"`
}

type UserSlice struct {
	Items []User `json:"items"`
	Total int    `json:"total"`
}

const (
	UserRole = iota
	UserManagerRole
	AdminRole
	OwnerRole
)

type Meal struct {
	ID              string `json:"id" db:"id"`
	Date            string `json:"date" db:"date"`
	Time            string `json:"time" db:"time"`
	Name            string `json:"name" db:"name"`
	Calories        int    `json:"calories" db:"calories"`
	CaloriesDeficit bool   `json:"calories_deficit" db:"calories_deficit"`
}

type MealSlice struct {
	Items []Meal `json:"items"`
	Total int    `json:"total"`
}

type Settings struct {
	ExpectedDailyCalories int `json:"expected_daily_calories" db:"expected_daily_calories"`
}

type UserDatastore interface {
	GetUserPassword(accountID, username string) (*string, error)
	SaveRootUser(username, pass string) (*User, error)

	GetUserById(accountID, userID string) (*User, error)
	GetUser(accountID, username string) (*User, error)
	GetUsers(accountID string, page, perPage int, filter string) (UserSlice, error)
	SaveUser(accountID, username, pass string, roleID int) (*User, error)
	UpdateUser(user User) (*User, error)
	DeleteUser(accountID, userID string) error

	SaveMeal(userID string, meal Meal) (*Meal, error)
	GetMeals(userID string, page, perPage int, filter string) (MealSlice, error)
	GetMeal(userID, mealID string) (*Meal, error)
	UpdateMeal(userID string, meal Meal) (*Meal, error)
	DeleteMeal(userID string, mealID string) error

	UpdateSettings(userID string, settings Settings) (*Settings, error)
	GetSettings(userID string) (*Settings, error)
}
