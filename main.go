package main

import (
	"fmt"
	"tabungan-api/api"
	"tabungan-api/app"
	"tabungan-api/repository"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logger := logrus.New()
	var database string
	var host string
	var port int
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if database = viper.GetString("DATABASE"); database == "" {
		database = "tabungan.db"
	}
	if host = viper.GetString("API_HOST"); host == "" {
		host = "0.0.0.0"
	}
	if port = viper.GetInt("API_PORT"); port == 0 {
		port = 8888
	}
	fmt.Print(host, port)
	repo := repository.InitDatabase(database, logger)
	app := app.NewTabunganApp(repo, logger)
	api := api.NewRESTAPI(host, port, app, logger)
	api.Start()
}
