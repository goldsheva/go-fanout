package configs

import (
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var config *Config
var once sync.Once

type Config struct {
	LogLevel       string
	AllowedMethods []string
	TargetUrls     []string
}

func GetEnvConfig() *Config {
	once.Do(func() {
		config = &Config{
			LogLevel: os.Getenv("LOG_LEVEL"),
		}

		if config.LogLevel == "" {
			config.LogLevel = "info"
		}

		allowedMethodsStr := os.Getenv("ALLOWED_METHODS")
		if allowedMethodsStr == "" {
			config.AllowedMethods = []string{}
		} else {
			config.AllowedMethods = strings.Split(allowedMethodsStr, ",")
		}

		urlsStr := os.Getenv("TARGET_URLS")
		if urlsStr == "" {
			config.TargetUrls = []string{}
		} else {
			urls := strings.Split(urlsStr, ",")
			for i := range urls {
				urls[i] = strings.TrimSpace(urls[i])

				if _, err := url.ParseRequestURI(urls[i]); err != nil {
					logrus.WithFields(logrus.Fields{"gopher": "main"}).Errorf("Invalid URL '%s' detected, skipping...\n", urls[i])
					continue
				}

				config.TargetUrls = append(config.TargetUrls, urls[i])
			}
		}
	})

	return config
}
