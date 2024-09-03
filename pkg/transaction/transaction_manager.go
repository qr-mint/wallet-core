package transaction

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type Manager interface {
	BeginTransaction() (driver.Tx, error)
	RollbackTransaction(tx driver.Tx) error
	CommitTransaction(tx driver.Tx) error
	WithTransaction(callback func(tx driver.Tx) error) error
}

type manager struct {
	db *sql.DB
}

func NewSQLManager(db *sql.DB) Manager {
	return &manager{db: db}
}

func (m manager) BeginTransaction() (driver.Tx, error) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %s", err)
	}

	return tx, nil
}

func (m manager) RollbackTransaction(tx driver.Tx) error {
	err := tx.Rollback()
	if err != nil {
		return fmt.Errorf("transaction rollback failed: %v", err)
	}

	return nil
}

func (m manager) CommitTransaction(tx driver.Tx) error {
	err := tx.Commit()
	if err != nil {
		return fmt.Errorf("transaction commit failed: %v", err)
	}

	return nil
}

func (m manager) WithTransaction(callback func(tx driver.Tx) error) error {
	tx, err := m.BeginTransaction()
	if err != nil {
		return fmt.Errorf("begin transaction error: %s", err)
	}
	callbackErr := callback(tx)
	if callbackErr != nil {
		err = tx.Rollback()
		if err != nil {
			return fmt.Errorf("transaction rollback failed: %s. while error occurs in transaction execution: %s", err, callbackErr)
		}
		return fmt.Errorf("transaction failed: %s", callbackErr)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("transaction commit error: %s", err)
	}
	return nil
}
