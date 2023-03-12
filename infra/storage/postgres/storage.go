package postgres

import (
	"errors"
	"fmt"
	"github/guiferpa/bank/domain/account"
	"github/guiferpa/bank/domain/log"

	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type PostgresStorage struct {
	db     *gorm.DB
	logger log.LoggerRepository
}

func (ps *PostgresStorage) CreateAccount(opts account.CreateAccountOptions) (uint, error) {
	model := &Account{DocumentNumber: opts.DocumentNumber}
	if err := ps.db.Create(model).Error; err != nil {
		return 0, account.NewInfraError(account.InfraUnknownError, err.Error())
	}

	return model.ID, nil
}

func (ps *PostgresStorage) GetAccountByID(accountID uint) (account.Account, error) {
	var dest Account
	if err := ps.db.Select("*").Where("id = ?", accountID).First(&dest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return account.Account{}, account.NewInfraError(account.InfraAccountNotFoundErrorCode, "account not found")
		}

		return account.Account{}, account.NewInfraError(account.InfraUnknownError, err.Error())
	}

	return account.Account{
		ID:             dest.ID,
		DocumentNumber: dest.DocumentNumber,
	}, nil
}

func (ps *PostgresStorage) HasAccountByDocumentNumber(documentNumber string) (bool, error) {
	var dest int64
	if err := ps.db.Model(&Account{}).Select("*").Where("document_number = ?", documentNumber).Count(&dest).Error; err != nil {
		return false, account.NewInfraError(account.InfraUnknownError, err.Error())
	}

	return dest > 0, nil
}

func (ps *PostgresStorage) HasOperationTypeByID(operationTypeID uint) (bool, error) {
	var dest int64
	if err := ps.db.Model(&OperationType{}).Select("*").Where("id = ?", operationTypeID).Count(&dest).Error; err != nil {
		return false, account.NewInfraError(account.InfraUnknownError, err.Error())
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
		return 0, account.NewInfraError(account.InfraUnknownError, err.Error())
	}

	return model.ID, nil
}

func (ps *PostgresStorage) RunSeed() error {
	return ps.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(OperationTypeSeedData, len(OperationTypeSeedData)).Error
}

type NewStorageOptions struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         string
	Logger       log.LoggerRepository
}

func NewStorage(opts NewStorageOptions) (*PostgresStorage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%v/%s?sslmode=%s", opts.User, opts.Password, opts.Host, opts.Port, opts.DatabaseName, "disable")
	db, err := gorm.Open(driver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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

	ps := &PostgresStorage{db, opts.Logger}

	return ps, nil
}
