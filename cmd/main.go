package main

import (
	"os"

	inkassback "github.com/Husenjon/InkassBack"
	"github.com/Husenjon/InkassBack/pkg/handler"
	"github.com/Husenjon/InkassBack/pkg/repository"
	"github.com/Husenjon/InkassBack/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// log := logger.New("", "inkass")
	if err := initConfig(); err != nil {
		logrus.Fatalf("erro config file %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("erro env file %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DBPASSWORD"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialized db %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)
	// go handler.GetDatasFrom1C()
	srv := new(inkassback.Server)
	if err := srv.Run(os.Getenv("PORT"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("error occuerd while run http server %s", err.Error())
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
