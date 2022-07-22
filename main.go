package main

import (
	"tabungan-api/repository"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logger := logrus.New()
	var database string
	viper.SetConfigFile("./.env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if database = viper.GetString("DATABASE"); database == "" {
		database = "tabungan.db"
	}
	db := repository.InitDatabase(database, logger)
	db.StartTransaction()
}