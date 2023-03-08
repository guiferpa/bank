package postgres

import (
	"context"
	"testing"
	"time"

	"github/guiferpa/bank/domain/account"
	"github/guiferpa/bank/infra/infratest"
)

func TestCreateAccount(t *testing.T) {
	env, err := infratest.NewEnvironment()
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	containerID, err := env.RunContainer(ctx, "postgres:14", "5432", []string{
		"POSTGRES_PASSWORD=Pa$$w0rd",
		"POSTGRES_DB=infra",
	})
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(2 * time.Second)

	newStorageOptions := NewStorageOptions{
		Host:         "localhost",
		User:         "postgres",
		Password:     "Pa$$w0rd",
		DatabaseName: "infra",
		Port:         "5432",
	}
	client, err := NewStorage(newStorageOptions)
	if err != nil {
		t.Error(err)
		return
	}

	suite := []struct {
		Describe string
		Spec     func(t *testing.T)
	}{
		{
			Describe: "Created successful",
			Spec: func(t *testing.T) {
				createAccountOptions := account.CreateAccountOptions{
					DocumentNumber: "42",
				}
				id, err := client.CreateAccount(createAccountOptions)
				if err != nil {
					t.Error(err)
					return
				}

				dest := Account{}
				if err := client.db.Select("*").Where("id = ?", id).Find(&dest).Error; err != nil {
					t.Error(err)
					return
				}

				if got, expected := dest.DocumentNumber, createAccountOptions.DocumentNumber; got != expected {
					t.Errorf("unexpected DocumentNumber, got: %s, expected: %s", got, expected)
					return
				}
			},
		},
		{
			Describe: "Duplicated account",
			Spec: func(t *testing.T) {
				createAccountOptions := account.CreateAccountOptions{
					DocumentNumber: "42",
				}
				_, err := client.CreateAccount(createAccountOptions)
				if err != nil {
					t.Error(err)
					return
				}
			},
		},
	}

	for _, s := range suite {
		t.Run(s.Describe, s.Spec)
	}

	defer func() {
		if err := env.KillContainer(ctx, containerID); err != nil {
			t.Error(err)
			return
		}
	}()
}
