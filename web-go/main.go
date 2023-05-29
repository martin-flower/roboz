package main

import (
	"os"
	"os/signal"

	"github.com/martin-flower/roboz-web-go/config"
	"github.com/martin-flower/roboz-web-go/database"
	swagger "github.com/martin-flower/roboz-web-go/docs"
	"github.com/martin-flower/roboz-web-go/httpserver"
	"github.com/martin-flower/roboz-web-go/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Main
// @title roboz cleaner
// @description example restful implementation of developer test task
// @version 1.0
// @contact.email gokonsulten@gmail.com
// @contact.name Martin Flower
// @host localhost:5000
// @BasePath /
func main() {

	logger.Setup()

	cfg := config.Read()

	err := database.Setup()
	if err != nil {
		zap.S().Fatalf("failed to start database %w", err)
	}
	defer database.DB.Close()

	swagger.SwaggerInfo.Host = viper.GetString("host")

	server := httpserver.New(cfg)

	server.Setup()

	// fiber graceful-shutdown using a channel
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if err := server.Shutdown(); err != nil {
			zap.S().Fatalf("failed to shutdown gracefully %+v", err)
		}
	}()

	if err := server.Listen(); err != nil {
		zap.S().Fatalf("failed serving connections %+v", err)
	}
}
