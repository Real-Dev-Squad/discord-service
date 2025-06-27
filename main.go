package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Real-Dev-Squad/discord-service/commands/register"
	config "github.com/Real-Dev-Squad/discord-service/config"
	queue "github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/routes"
	"github.com/sirupsen/logrus"
)

func main() {
	register.SetupRegister()
	logrus.Info("Starting server on port " + config.AppConfig.Port)

	// Initialize queue
	_, err := queue.GetQueueInstance()
	if err != nil {
		logrus.Fatalf("Failed to initialize queue: %v", err)
	}

	// Create a channel to listen for OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		routes.Listen(":" + config.AppConfig.Port)
	}()

	// Wait for shutdown signal
	<-sigChan
	logrus.Info("Shutting down gracefully...")

	// Cleanup
	if err := queue.CloseQueueConnection(); err != nil {
		logrus.Errorf("Error during queue cleanup: %v", err)
	}

	logrus.Info("Shutdown complete")
}
