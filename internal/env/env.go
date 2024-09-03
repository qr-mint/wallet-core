package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	App          App
	Http         Http
	Jwt          Jwt
	PgSql        PgSql
	Redis        Redis
	Telegram     Telegram
	Sync         Sync
	Integrations Integrations
	Sentry       Sentry
	Elastic      Elastic
}

type App struct {
	Mode      string
	SecretKey string
	RateLimit int
}

type Http struct {
	Host            string
	Port            string
	HeaderOriginURL string
}

type Jwt struct {
	AccessTokenMaxLifeInMinutes  int
	RefreshTokenMaxLifeInMinutes int
}

type PgSql struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     int
	SslMode  string
}

type Redis struct {
	Host     string
	Password string
	Port     int
	Database int
}

type Sentry struct {
	Dsn string
}

type Telegram struct {
	BotToken               string
	BotID                  int
	QueryLifetimeInSeconds int
}

type Sync struct {
	PricesFrequencyInHours      int
	TransactionsIntervalMinutes int
	BalanceIntervalMinutes      int
	NftIntervalMinutes          int
}

type Integrations struct {
	Tonconsole struct {
		Host string
	}
	Toncenter struct {
		Host string
		Key  string
	}
	Tronscanapi struct {
		Host string
		Key  string
	}
	Trongrid struct {
		Host string
		Key  string
	}
	Tonhub struct {
		ConnectHost string
		MainnetHost string
	}
	Coingecko struct {
		Host string
	}
	SimpleSwap struct {
		Host   string
		ApiKey string
	}
	FinchPay struct {
		WidgetHost string
		Host       string
		PartnerId  string
		SecretKey  string
	}
	ChangeHero struct {
		Host   string
		Apikey string
	}
}

type Elastic struct {
	Host     string
	Port     string
	Password string
}

const appMode = "APP_MODE"
const appSecretKey = "APP_SECRET_KEY"
const appRateLimit = "APP_RATE_LIMIT"
const httpHost = "HTTP_HOST"
const httpPort = "HTTP_PORT"
const httpHeaderOriginURL = "HTTP_HEADER_ORIGIN_URL"

const jwtAccessTokenMaxLifetimeInMinutes = "JWT_ACCESS_TOKEN_MAX_LIFETIME_IN_MINUTES"
const jwtRefreshTokenMaxLifetimeInMinutes = "JWT_REFRESH_TOKEN_MAX_LIFETIME_IN_MINUTES"

const pgSqlHost = "PG_HOST"
const pgSqlUser = "PG_USER"
const pgSqlPassword = "PG_PASSWORD"
const pgSqlDbName = "PG_DB_NAME"
const pgSqlPort = "PG_PORT"
const pgSqlSslMode = "PG_SSL_MODE"

const redisHost = "REDIS_HOST"
const redisPassword = "REDIS_PASSWORD"
const redisPort = "REDIS_PORT"
const redisDatabase = "REDIS_DATABASE"

const telegramBotToken = "TELEGRAM_BOT_TOKEN"
const telegramBotId = "TELEGRAM_BOT_ID"
const telegramQueryLifetimeInSeconds = "TELEGRAM_QUERY_LIFETIME_IN_SECONDS"

const syncPricesFrequencyInHours = "SYNC_PRICES_FREQUENCY_IN_HOURS"
const syncTransactionsIntervalMinutes = "SYNC_TRANSACTIONS_INTERVAL_MINUTES"
const syncBalanceIntervalMinutes = "SYNC_BALANCE_INTERVAL_MINUTES"
const syncNftIntervalMinutes = "SYNC_NFT_INTERVAL_MINUTES"

const integrationsTonconsoleHost = "INTEGRATIONS_TONCONSOLE_HOST"
const integrationsToncenterHost = "INTEGRATIONS_TONCENTER_HOST"
const integrationsToncenterKey = "INTEGRATIONS_TONCENTER_KEY"
const integrationsTonhubConnectHost = "INTEGRATIONS_TONHUB_CONNECT_HOST"
const integrationsTonhubMainnetHost = "INTEGRATIONS_TONHUB_MAINNET_HOST"

const integrationsTronscanapiHost = "INTEGRATIONS_TRONSCANAPI_HOST"
const integrationsTronscanapiKey = "INTEGRATIONS_TRONSCANAPI_KEY"
const integrationsTrongridHost = "INTEGRATIONS_TRONGRID_HOST"
const integrationsTrongridKey = "INTEGRATIONS_TRONGRID_KEY"

const integrationsCoingeckoHost = "INTEGRATIONS_COINGECKO_HOST"

const integrationsSimpleSwapHost = "INTEGRATIONS_SIMPLESWAP_HOST"
const integrationsSimpleSwapApiKey = "INTEGRATIONS_SIMPLESWAP_API_KEY"

const integrationsFinchPayWidgetHost = "INTEGRATIONS_FINCHPAY_WIDGET_HOST"
const integrationsFinchPayHost = "INTEGRATIONS_FINCHPAY_HOST"
const integrationsFinchPayPartnerId = "INTEGRATIONS_FINCHPAY_PARTNER_ID"
const integrationsFinchPaySecretKey = "INTEGRATIONS_FINCHPAY_SECRET_KEY"

const integrationsChangeHeroHost = "INTEGRATIONS_CHANGEHERO_HOST"
const integrationsChangeHeroApiKey = "INTEGRATIONS_CHANGEHERO_API_KEY"

const sentryDsn = "SENTRY_DSN"

const elasticsearchHost = "ELASTICSEARCH_HOST"
const elasticsearchPort = "ELASTICSEARCH_PORT"
const elasticsearchPassword = "ELASTICSEARCH_PASSWORD"

func (e *Env) Load() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("env load error: %w", err)
	}

	appRateLimit := os.Getenv(appRateLimit)
	appRateLimitInt, err := strconv.Atoi(appRateLimit)
	if err != nil {
		return fmt.Errorf("env load appRateLimit error: %w", err)
	}

	e.App = App{
		Mode:      os.Getenv(appMode),
		SecretKey: os.Getenv(appSecretKey),
		RateLimit: appRateLimitInt,
	}

	e.Http = Http{
		Host:            os.Getenv(httpHost),
		Port:            os.Getenv(httpPort),
		HeaderOriginURL: os.Getenv(httpHeaderOriginURL),
	}

	accessTokenMaxLifetime := os.Getenv(jwtAccessTokenMaxLifetimeInMinutes)
	accessTokenMaxLifetimeInt, err := strconv.Atoi(accessTokenMaxLifetime)
	if err != nil {
		return fmt.Errorf("env load jwtAccessTokenMaxLifetimeInMinutes error: %w", err)
	}

	refreshTokenMaxLifetime := os.Getenv(jwtRefreshTokenMaxLifetimeInMinutes)
	refreshTokenMaxLifetimeInt, err := strconv.Atoi(refreshTokenMaxLifetime)
	if err != nil {
		return fmt.Errorf("env load jwtRefreshTokenMaxLifetimeInMinutes error: %w", err)
	}

	e.Jwt = Jwt{
		AccessTokenMaxLifeInMinutes:  accessTokenMaxLifetimeInt,
		RefreshTokenMaxLifeInMinutes: refreshTokenMaxLifetimeInt,
	}

	port, err := strconv.Atoi(os.Getenv(pgSqlPort))
	if err != nil {
		return fmt.Errorf("env load pgSqlPort error: %w", err)
	}
	e.PgSql = PgSql{
		Host:     os.Getenv(pgSqlHost),
		User:     os.Getenv(pgSqlUser),
		Password: os.Getenv(pgSqlPassword),
		DbName:   os.Getenv(pgSqlDbName),
		Port:     port,
		SslMode:  os.Getenv(pgSqlSslMode),
	}

	redisPort, err := strconv.Atoi(os.Getenv(redisPort))
	if err != nil {
		return fmt.Errorf("env load redisPort error: %w", err)
	}
	redisDatabase, err := strconv.Atoi(os.Getenv(redisDatabase))
	if err != nil {
		return fmt.Errorf("env load redisDatabase error: %w", err)
	}
	e.Redis = Redis{
		Host:     os.Getenv(redisHost),
		Password: os.Getenv(redisPassword),
		Port:     redisPort,
		Database: redisDatabase,
	}

	telegramQueryLifetimeInSeconds := os.Getenv(telegramQueryLifetimeInSeconds)
	telegramQueryLifetimeInSecondsInt, err := strconv.Atoi(telegramQueryLifetimeInSeconds)
	if err != nil {
		return fmt.Errorf("env load telegramQueryLifetimeInSeconds error: %w", err)
	}
	telegramBotId := os.Getenv(telegramBotId)
	telegramBotIdInt, err := strconv.Atoi(telegramBotId)
	if err != nil {
		return fmt.Errorf("env load telegramBotId error: %w", err)
	}

	e.Telegram = Telegram{
		BotToken:               os.Getenv(telegramBotToken),
		BotID:                  telegramBotIdInt,
		QueryLifetimeInSeconds: telegramQueryLifetimeInSecondsInt,
	}

	syncPricesFrequencyInHours := os.Getenv(syncPricesFrequencyInHours)
	syncPricesFrequencyInHoursInt, err := strconv.Atoi(syncPricesFrequencyInHours)
	if err != nil {
		return fmt.Errorf("env load syncPricesFrequencyInHours error: %w", err)
	}
	syncTransactionsIntervalMinutes := os.Getenv(syncTransactionsIntervalMinutes)
	syncTransactionsIntervalMinutesInt, err := strconv.Atoi(syncTransactionsIntervalMinutes)
	if err != nil {
		return fmt.Errorf("env load syncTransactionsIntervalMinutes error: %w", err)
	}
	syncBalanceIntervalMinutes := os.Getenv(syncBalanceIntervalMinutes)
	syncBalanceIntervalMinutesInt, err := strconv.Atoi(syncBalanceIntervalMinutes)
	if err != nil {
		return fmt.Errorf("env load syncBalanceIntervalMinutes error: %w", err)
	}
	syncNftIntervalMinutes := os.Getenv(syncNftIntervalMinutes)
	syncNftIntervalMinutesInt, err := strconv.Atoi(syncNftIntervalMinutes)
	if err != nil {
		return fmt.Errorf("env load syncNftIntervalMinutes error: %w", err)
	}

	e.Sync = Sync{
		PricesFrequencyInHours:      syncPricesFrequencyInHoursInt,
		TransactionsIntervalMinutes: syncTransactionsIntervalMinutesInt,
		BalanceIntervalMinutes:      syncBalanceIntervalMinutesInt,
		NftIntervalMinutes:          syncNftIntervalMinutesInt,
	}

	e.Integrations = Integrations{
		Tonconsole: struct{ Host string }{Host: os.Getenv(integrationsTonconsoleHost)},
		Toncenter: struct {
			Host string
			Key  string
		}{Host: os.Getenv(integrationsToncenterHost), Key: os.Getenv(integrationsToncenterKey)},

		Tonhub: struct {
			ConnectHost string
			MainnetHost string
		}{ConnectHost: os.Getenv(integrationsTonhubConnectHost), MainnetHost: os.Getenv(integrationsTonhubMainnetHost)},

		Trongrid: struct {
			Host string
			Key  string
		}{Host: os.Getenv(integrationsTrongridHost), Key: os.Getenv(integrationsTrongridKey)},
		Tronscanapi: struct {
			Host string
			Key  string
		}{Host: os.Getenv(integrationsTronscanapiHost), Key: os.Getenv(integrationsTronscanapiKey)},

		Coingecko: struct{ Host string }{Host: os.Getenv(integrationsCoingeckoHost)},

		SimpleSwap: struct {
			Host   string
			ApiKey string
		}{Host: os.Getenv(integrationsSimpleSwapHost), ApiKey: os.Getenv(integrationsSimpleSwapApiKey)},

		FinchPay: struct {
			WidgetHost string
			Host       string
			PartnerId  string
			SecretKey  string
		}{WidgetHost: os.Getenv(integrationsFinchPayWidgetHost), Host: os.Getenv(integrationsFinchPayHost), PartnerId: os.Getenv(integrationsFinchPayPartnerId), SecretKey: os.Getenv(integrationsFinchPaySecretKey)},

		ChangeHero: struct {
			Host   string
			Apikey string
		}{Host: os.Getenv(integrationsChangeHeroHost), Apikey: os.Getenv(integrationsChangeHeroApiKey)},
	}

	e.Sentry = Sentry{
		Dsn: os.Getenv(sentryDsn),
	}

	e.Elastic = Elastic{
		Host:     os.Getenv(elasticsearchHost),
		Port:     os.Getenv(elasticsearchPort),
		Password: os.Getenv(elasticsearchPassword),
	}

	return nil
}
