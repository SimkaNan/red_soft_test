package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	lib "golibrary"
	"golibrary/cmd/logConfig"
	_ "golibrary/docs"
	"golibrary/init"
	cntrlr "golibrary/internal/controller"
	"golibrary/internal/repository"
	"golibrary/internal/service"
	"golibrary/pkg/debug"
	"golibrary/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// @title users-library
// @version 1.0
// @description API server for library

// @host localhost:8080
// @BasePath /

var (
	version = ""
	log     *zap.SugaredLogger
	wg      sync.WaitGroup
)

func main() {
	loggerConfig, err := logConfig.ReadConfig()
	if err != nil {
		panic("failed to load readConfig: " + err.Error())
	}

	logger.BuildLogger(loggerConfig.LogLevel)
	log = logger.Logger().Named("main").Sugar()

	if err := InitConfig(); err != nil {
		log.Fatal("error initalizing configs: ", zap.Error(err))
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("error to open env file:", zap.Error(err))
	}

	db, err := repository.NewPostgresDB(
		repository.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: viper.GetString("db.username"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   viper.GetString("db.dbname"),
			SSLMode:  viper.GetString("db.sslmode"),
		})
	if err != nil {
		log.Fatal("Error to init db:", zap.Error(err))
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	controller := cntrlr.NewController(services, log)

	check := repos.Init.CheckEmpty()

	err = repos.Init.InitDB()
	if err != nil {
		log.Fatal("Error to init db:", zap.Error(err))
	}

	time.Sleep(time.Second * 3)

	if check {
		err = initData.Init(services)
		if err != nil {
			log.Error("error initalizing tables:", zap.Error(err))
		}
		if err != nil {
			panic(err)
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	wg.Add(1)

	go debug.Run(":40000")

	srv := new(lib.Server)

	go func() {
		defer wg.Done()

		if err := srv.Run(viper.GetString("port"), controller.InitRoutes()); err != nil {
			log.Fatal("error occuring while running http server:", zap.Error(err))
		}
	}()

	<-stop

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error shutting down server: %s", err))
	}
	wg.Wait()

	log.Info(fmt.Sprintf("Server stopped"))
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
