package main

import (
	"calories-counter/adapters/calories_datastore"
	"calories-counter/adapters/user_datastore"
	"calories-counter/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	secretKey = os.Getenv("TOKEN_SECRET")
	dbSource  = os.Getenv("MYSQL_DB_SOURCE")
	appID     = os.Getenv("API_APP_ID")
	apiKey    = os.Getenv("API_KEY")
)

func main() {
	caloriesDatastore := calories_datastore.NewNutritionixApi(appID, apiKey)
	userDatastore, err := user_datastore.NewMySQLStore(dbSource)
	if err != nil {
		log.Error(err)
	}
	defer func() { _ = userDatastore.Close() }()

	r := gin.Default()
	server.SetupRouter(r, secretKey, userDatastore, caloriesDatastore)
	_ = r.Run(":8000")
}
