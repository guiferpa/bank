package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github/guiferpa/bank/domain/account"
	logd "github/guiferpa/bank/domain/log"
	"github/guiferpa/bank/handler/http/api"
	"github/guiferpa/bank/infra/logger/log"
	"github/guiferpa/bank/infra/storage/postgres"
)

func main() {
	value := logd.LoggerContext{
		RequestID: fmt.Sprintf("%v", time.Now().Unix()),
	}
	ctx := context.WithValue(context.Background(), logd.LoggerContextKey, &value)

	logger := log.NewLogger()
	storage, err := postgres.NewStorage(postgres.NewStorageOptions{
		Host:         os.Getenv("DATABASE_HOST"),
		User:         os.Getenv("DATABASE_USER"),
		Password:     os.Getenv("DATABASE_PASSWORD"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		Port:         os.Getenv("DATABASE_PORT"),
		Logger:       logger,
	})
	if err != nil {
		logger.Error(ctx, err.Error())
		return
	}
	service := account.NewUseCaseService(storage, logger)
	handler := api.NewHTTPHandler(service, logger)

	// Reason which make me to do seed on my own hand: https://github.com/go-gorm/gorm/issues/5339
	if err := storage.RunSeed(); err != nil {
		logger.Error(ctx, err.Error())
		return
	}
	logger.Info(ctx, "seed done successful")

	port := os.Getenv("PORT")

	logger.Info(ctx, fmt.Sprintf("API's running at port %v", port))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		logger.Error(ctx, err.Error())
	}
}
