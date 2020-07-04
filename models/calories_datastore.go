package models

//go:generate mockgen -destination=../adapters/calories_datastore/mock.go -package=calories_datastore calories-counter/models CaloriesDatastore

type CaloriesDatastore interface {
	GetCalories(mealName string) (*int, error)
}
