package controllers

import (
	"bytes"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/goldsheva/go-fanout/internal/configs"
	"github.com/sirupsen/logrus"
)

func FanoutHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	ctx := r.Context()
	config := configs.GetEnvConfig()

	if strings.Contains(r.URL.Path, "favicon.ico") {
		return
	}

	logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Infof("<<< %s: %s", r.Method, r.URL.Path)

	if !slices.Contains(config.AllowedMethods, r.Method) {
		logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Infof("Method '%s' is not allowed! skipping...\n", r.Method)
		w.Write([]byte("OK"))
		return
	}

	var bodyBytes []byte
	var err error
	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Errorf("Failed to read request body: %v", err)
			return
		}

	}

	logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Debugf("HEADERS: %v", r.Header)
	logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Debugf("BODY: %s", string(bodyBytes))

	for _, url := range config.TargetUrls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return

			default:

				req, err := http.NewRequest(r.Method, url, io.NopCloser(bytes.NewReader(bodyBytes)))
				if err != nil {
					logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Errorf(">>> %s: %s (%v)", r.Method, u, err)
					return
				}

				req.Header = r.Header.Clone()

				client := &http.Client{
					Timeout: 10 * time.Second,
				}

				resp, err := client.Do(req)
				if err != nil {
					logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Warningf(">>> %s: %s (%v)", r.Method, u, err)
					return
				}
				defer resp.Body.Close()

				logrus.WithFields(logrus.Fields{"gopher": "http_server"}).Infof(">>> %s: %s (Status %s)", r.Method, u, resp.Status)
			}
		}(url)
	}
	wg.Wait()

	w.Write([]byte("OK"))
}
