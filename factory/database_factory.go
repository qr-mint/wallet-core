package factory

import (
	"database/sql"
	"fmt"
	pgsql "gitlab.com/golib4/psql-connection/connection"
	"nexus-wallet/pkg/repository"
	"nexus-wallet/pkg/transaction"
	"time"
)

func (f *ServiceFactory) createConnection() (*pgsql.Connection, error) {
	connection, err := pgsql.NewConnection(pgsql.Params{
		Host:     f.env.PgSql.Host,
		User:     f.env.PgSql.User,
		Password: f.env.PgSql.Password,
		DbName:   f.env.PgSql.DbName,
		Port:     int32(f.env.PgSql.Port),
		SslMode:  f.env.PgSql.SslMode,

		MaxOpenConns:    75,
		MaxIdleConns:    20,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("pgsql connection failed: %w", err)
	}

	return connection, nil
}

func (f *ServiceFactory) createBaseRepository(db *sql.DB) (*repository.BaseRepository, error) {
	newLogger, err := f.createLogger("sql query")
	if err != nil {
		return nil, fmt.Errorf("can not create logger while creating baseRepository: %s", err)
	}

	baseRepository := repository.NewBaseRepository(db, newLogger, repository.Params{LogInfo: true})
	return &baseRepository, nil
}

func (f *ServiceFactory) createTransactionManager(db *sql.DB) transaction.Manager {
	return transaction.NewSQLManager(db)
}
