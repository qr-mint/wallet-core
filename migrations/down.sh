set -a
source .env
set +a

migrate -database postgres://${PG_USER}:${PG_PASSWORD}@${PG_HOST}:${PG_PORT}/${PG_DB_NAME}?sslmode=${PG_SSL_MODE} -path ./migrations down