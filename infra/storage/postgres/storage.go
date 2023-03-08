package postgres

import (
	"fmt"
	"github/guiferpa/bank/domain/account"

	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func (ps *PostgresStorage) CreateAccount(opts account.CreateAccountOptions) (uint, error) {
	model := &Account{DocumentNumber: opts.DocumentNumber}
	if err := ps.db.Create(model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

type NewStorageOptions struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         string
}

func NewStorage(opts NewStorageOptions) (*PostgresStorage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%s", opts.User, opts.Password, opts.Host, opts.Port, opts.DatabaseName, "disable")
	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	exr, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := exr.Ping(); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&Account{}); err != nil {
		return nil, err
	}

	return &PostgresStorage{db}, nil
}
