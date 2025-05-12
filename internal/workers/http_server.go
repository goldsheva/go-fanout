package workers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/goldsheva/go-fanout/internal/controllers"
	"github.com/sirupsen/logrus"
)

func GoHTTPServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	mux := http.NewServeMux()
	mux.HandleFunc("/", controllers.FanoutHandler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	go func() {
		<-ctx.Done()

		logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Warn("HTTP server stopped...")

		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Errorf("HTTP server shutdown error: %v\n", err)
		}
	}()

	logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Infof("Starting HTTP server on :%s", os.Getenv("HTTP_PORT"))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Errorf("HTTP server error: %v\n", err)
	}
}
