package main

import (
	"fmt"
	"github/guiferpa/bank/domain/account"
	"github/guiferpa/bank/handler/http/rest"
	"github/guiferpa/bank/infra/storage/postgres"
	"net/http"
	"os"
)

func main() {
	storage, err := postgres.NewStorage(postgres.NewStorageOptions{
		Host:         os.Getenv("DATABASE_HOST"),
		User:         os.Getenv("DATABASE_USER"),
		Password:     os.Getenv("DATABASE_PASSWORD"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		Port:         os.Getenv("DATABASE_PORT"),
	})
	if err != nil {
		panic(err)
	}
	service := account.NewUseCaseService(storage)
	handler := rest.NewHTTPHandler(service)

	port := os.Getenv("PORT")
	http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}
