//go:build integration

package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github/guiferpa/bank/pkg/docker"
)

func TestIntegrationForAPI(t *testing.T) {
	client, err := docker.NewEnvironment()
	if err != nil {
		t.Error(err)
		return
	}

	dbID, err := client.RunContainer(context.Background(), "postgres:14", "5431", "5432", []string{
		"POSTGRES_PASSWORD=Pa$$w0rd",
		"POSTGRES_DB=api",
	})
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		if err := client.KillContainer(context.Background(), dbID); err != nil {
			t.Error(err)
			return
		}
	}()

	t.Setenv("DATABASE_HOST", "localhost")
	t.Setenv("DATABASE_USER", "postgres")
	t.Setenv("DATABASE_PORT", "5431")
	t.Setenv("DATABASE_NAME", "api")
	t.Setenv("DATABASE_PASSWORD", "Pa$$w0rd")
	t.Setenv("PORT", "8080")

	time.Sleep(4 * time.Second)

	go main()

	time.Sleep(2 * time.Second)

	suite := []struct {
		Describe string
		Spec     func(t *testing.T)
	}{
		{
			Describe: "Created account successful",
			Spec: func(t *testing.T) {
				body := bytes.NewBufferString(`{"document_number": "10"}`)
				resp, err := http.Post("http://localhost:8080/api/v1/accounts", "application/json; chartset=utf-8", body)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusCreated; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"id\":1,\"document_number\":\"10\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Got duplicated account error when create account successful",
			Spec: func(t *testing.T) {
				body := bytes.NewBufferString(`{"document_number": "10"}`)
				resp, err := http.Post("http://localhost:8080/api/v1/accounts", "application/json; chartset=utf-8", body)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusConflict; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"code\":\"domain.1\",\"message\":\"account already exists\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Got EOF error when create account successful",
			Spec: func(t *testing.T) {
				body := bytes.NewBufferString("")
				resp, err := http.Post("http://localhost:8080/api/v1/accounts", "application/json; chartset=utf-8", body)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusBadRequest; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"code\":\"handler.3\",\"message\":\"missing request body\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Got invalid request body error when create account successful",
			Spec: func(t *testing.T) {
				body := bytes.NewBufferString(`{"document_number": ""}`)
				resp, err := http.Post("http://localhost:8080/api/v1/accounts", "application/json; chartset=utf-8", body)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusUnprocessableEntity; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"code\":\"handler.2\",\"message\":\"field document_number cannot be empty\",\"field\":\"document_number\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Got account already created",
			Spec: func(t *testing.T) {
				resp, err := http.Get("http://localhost:8080/api/v1/accounts/1")
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusOK; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"id\":1,\"document_number\":\"10\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Got account not found",
			Spec: func(t *testing.T) {
				resp, err := http.Get("http://localhost:8080/api/v1/accounts/1398")
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusNotFound; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"code\":\"infra.2\",\"message\":\"account not found\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Created account transaction successful",
			Spec: func(t *testing.T) {
				body := bytes.NewBufferString(`{"account_id": 1, "operation_type_id": 1, "amount": 15.45}`)
				resp, err := http.Post("http://localhost:8080/api/v1/accounts/transaction", "application/json; charset=utf-8", body)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusCreated; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"id\":1}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Got account not found when create account transaction",
			Spec: func(t *testing.T) {
				body := bytes.NewBufferString(`{"account_id": 8000, "operation_type_id": 1, "amount": 15.45}`)
				resp, err := http.Post("http://localhost:8080/api/v1/accounts/transaction", "application/json; charset=utf-8", body)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusNotFound; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"code\":\"infra.2\",\"message\":\"account not found\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Got operation type not found when create account transaction",
			Spec: func(t *testing.T) {
				body := bytes.NewBufferString(`{"account_id": 1, "operation_type_id": 10, "amount": 15.45}`)
				resp, err := http.Post("http://localhost:8080/api/v1/accounts/transaction", "application/json; charset=utf-8", body)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusNotFound; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"code\":\"domain.2\",\"message\":\"operation type doesn't exist\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Got error when create account transaction with amount equals zero",
			Spec: func(t *testing.T) {
				body := bytes.NewBufferString(`{"account_id": 1, "operation_type_id": 2, "amount": 0}`)
				resp, err := http.Post("http://localhost:8080/api/v1/accounts/transaction", "application/json; charset=utf-8", body)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := resp.StatusCode, http.StatusUnprocessableEntity; got != expected {
					t.Errorf("unexpected response status code, got: %v, expected: %v", got, expected)
					return
				}

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Error(err)
					return
				}
				defer resp.Body.Close()

				if got, expected := string(data), "{\"code\":\"handler.2\",\"message\":\"this value can't be zero\",\"field\":\"amount\"}\n"; got != expected {
					t.Errorf("unexpected response body, got: %v, expected: %v", got, expected)
					return
				}
			},
		},
	}

	for _, s := range suite {
		t.Run(s.Describe, s.Spec)
	}
}
