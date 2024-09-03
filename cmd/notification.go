package main

import (
	"log"
	"nexus-wallet/internal/env"
	"time"

	pgsql "gitlab.com/golib4/psql-connection/connection"
)

// DATABASE_HOST=wallet-postgres
// DATABASE_PASSWORD=HivGYs2Jol
// DATABASE_NAME=bot_users
// DATABASE_USER=nexus
// DATABASE_PORT=5433

func main() {
	envData := &env.Env{}
	err := envData.Load()
	connection, err := pgsql.NewConnection(pgsql.Params{
		Host:            envData.PgSql.Host,
		User:            envData.PgSql.User,
		Password:        envData.PgSql.Password,
		DbName:          envData.PgSql.DbName,
		Port:            int32(envData.PgSql.Port),
		SslMode:         envData.PgSql.SslMode,
		MaxOpenConns:    75,
		MaxIdleConns:    20,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil {
		log.Fatalf("pgsql connection failed: %w", err)
	}

	smt, err := connection.GetConnection().Prepare(`SELECT telegram_id FROM bot_users`)
}
