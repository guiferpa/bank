package postgres

import (
	"context"
	"testing"
	"time"

	"github/guiferpa/bank/domain/account"
	"github/guiferpa/bank/infra/infratest"
)

func TestRunSeed(t *testing.T) {
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

	suite := []struct {
		Describe string
		Spec     func(t *testing.T)
	}{
		{
			Describe: "Seed done successful",
			Spec: func(t *testing.T) {
				client, err := NewStorage(newStorageOptions)
				if err != nil {
					t.Error(err)
					return
				}

				dest := make([]OperationType, 0)
				if err := client.db.Select("*").Find(&dest).Error; err != nil {
					t.Error(err)
					return
				}

				if got, expected := len(dest), len(OperationTypeSeedData); got != expected {
					t.Errorf("unexpected value for find in operation_types table, got count: %v, expected count: %v", got, expected)
					return
				}
			},
		},
		{
			Describe: "Seed duplicated ignored sucessful",
			Spec: func(t *testing.T) {
				_, err := NewStorage(newStorageOptions)
				if err != nil {
					t.Error(err)
					return
				}

				client, err := NewStorage(newStorageOptions)
				if err != nil {
					t.Error(err)
					return
				}

				dest := make([]OperationType, 0)
				if err := client.db.Select("*").Find(&dest).Error; err != nil {
					t.Error(err)
					return
				}

				if got, expected := len(dest), len(OperationTypeSeedData); got != expected {
					t.Errorf("unexpected value for find in operation_types table, got count: %v, expected count: %v", got, expected)
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
				if _, err := client.CreateAccount(createAccountOptions); err == nil {
					t.Errorf("unexpected value for error, got: %v", err)
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

func TestGetAccountByID(t *testing.T) {
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
			Describe: "Got successful",
			Spec: func(t *testing.T) {
				documentNumber := "42"
				data := Account{DocumentNumber: documentNumber}
				if err := client.db.Create(&data).Error; err != nil {
					t.Error(err)
					return
				}

				acc, err := client.GetAccountByID(data.ID)
				if err != nil {
					t.Error(err)
					return
				}

				if got, expected := acc.DocumentNumber, documentNumber; got != expected {
					t.Errorf("unexpected DocumentNumber, got: %s, expected: %s", got, expected)
					return
				}
			},
		},
		{
			Describe: "Not found",
			Spec: func(t *testing.T) {
				_, err := client.GetAccountByID(2)
				if _, ok := err.(*account.StorageRepositoryGetAccountByIDError); !ok {
					t.Errorf("unexpected value for error, got: %v", err)
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

func TestCreateTransaction(t *testing.T) {
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
				acc := &Account{DocumentNumber: "42"}
				if err := client.db.Create(&acc).Error; err != nil {
					t.Error(err)
					return
				}

				transOptions := account.CreateTransactionOptions{
					AccountID:       acc.ID,
					OperationTypeID: 1, // COMPRA A VISTA
					Amount:          -10_00,
					EventDate:       time.Now(),
				}
				transID, err := client.CreateTransaction(transOptions)
				if err != nil {
					t.Error(err)
					return
				}

				var dest AccountTransaction
				if err := client.db.Select("*").Where("id = ?", transID).First(&dest).Error; err != nil {
					t.Error(err)
					return
				}

				if got, expected := dest.Amount, transOptions.Amount; got != expected {
					t.Errorf("unexpected transaction's Amount, got: %v, expected: %v", got, expected)
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
