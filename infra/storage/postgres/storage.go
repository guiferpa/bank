package postgres

import (
	"errors"
	"fmt"
	"github/guiferpa/bank/domain/account"

	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (ps *PostgresStorage) GetAccountByID(accountID uint) (account.Account, error) {
	var dest Account
	if err := ps.db.Select("*").Where("id = ?", accountID).First(&dest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return account.Account{}, account.NewStorageRepositoryGetAccountByIDError(account.StorageAccountNotFoundErrorCode, "account not found")
		}

		return account.Account{}, err
	}

	return account.Account{
		ID:             dest.ID,
		DocumentNumber: dest.DocumentNumber,
	}, nil
}

func (ps *PostgresStorage) HasAccountByDocumentNumber(documentNumber string) (bool, error) {
	var dest int64
	if err := ps.db.Model(&Account{}).Select("*").Where("document_number = ?", documentNumber).Count(&dest).Error; err != nil {
		return false, err
	}

	return dest > 0, nil
}

func (ps *PostgresStorage) CreateTransaction(opts account.CreateTransactionOptions) (uint, error) {
	model := &AccountTransaction{
		AccountID:       opts.AccountID,
		OperationTypeID: opts.OperationTypeID,
		Amount:          opts.Amount,
		EventDate:       opts.EventDate,
	}
	if err := ps.db.Create(&model).Error; err != nil {
		return 0, err
	}

	return model.ID, nil
}

func (ps *PostgresStorage) runSeed() error {
	return ps.db.Debug().Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(OperationTypeSeedData, len(OperationTypeSeedData)).Error
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

	if err := db.AutoMigrate(&Account{}, &OperationType{}, &AccountTransaction{}); err != nil {
		return nil, err
	}

	ps := &PostgresStorage{db}

	// Reason which make me to do seed on my own hand: https://github.com/go-gorm/gorm/issues/5339
	if err := ps.runSeed(); err != nil {
		return nil, err
	}

	return ps, nil
}
