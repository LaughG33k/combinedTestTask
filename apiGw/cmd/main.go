package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apigw "github.com/LaughG33k/apiGW"

	"gopkg.in/yaml.v2"
)

func main() {

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	apigw.InitLogrus("./logs.json")

	var cfg apigw.Cfg

	if err := initConfig("./config.yaml", &cfg); err != nil {
		apigw.Logger.Error(err)
		return
	}

	parser, err := initParser(fmt.Sprintf("http://%s/publicKey", cfg.AuthApi), 25*time.Second)

	if err != nil {
		apigw.Logger.Errorf("Failed to get public key. Error: %s", err)
		return
	}

	proxy := apigw.InitProxy(cfg, parser)

	httpServer := http.Server{
		Addr:         cfg.Addr,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutInSec) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutInSec) * time.Second,
	}

	apigw.GenerateHandlers(proxy, cfg.EndPoint)

	go httpServer.ListenAndServe()

	<-mainCtx.Done()

	httpServer.Shutdown(context.Background())

}

func initParser(authApiUrl string, timeout time.Duration) (*apigw.JwtParser, error) {

	var parser *apigw.JwtParser
	err := apigw.Rerty(func() error {

		key, err := apigw.GetKey(authApiUrl, timeout)

		if err != nil {
			return err
		}

		p, err := apigw.NewParser(key.Key, key.Method)

		if err != nil {
			return err
		}

		parser = p

		return nil

	}, 5, 1*time.Second)

	if err != nil {
		return nil, err
	}

	return parser, nil

}

func initConfig(path string, dest any) error {

	bytes, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(bytes, dest); err != nil {
		return err
	}

	return nil

}
