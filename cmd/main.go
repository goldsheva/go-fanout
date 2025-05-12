package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/goldsheva/go-fanout/internal/configs"
	"github.com/goldsheva/go-fanout/internal/workers"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("No .env file found")
	}

	config := configs.GetEnvConfig()

	switch config.LogLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	wg.Add(1)
	go workers.GoHTTPServer(ctx, wg)

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan
	logrus.WithFields(logrus.Fields{"gopher": "main"}).Warn("Initiating shutdown...")
	cancelFunc()

	wg.Wait()
	logrus.WithFields(logrus.Fields{"gopher": "main"}).Warn("Shutdown complete. All processes stopped!")
}
