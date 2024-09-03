FROM golang:1.22.2

WORKDIR /usr/src/wallet-core

COPY . .

RUN go mod vendor
RUN go mod download
RUN go mod verify

RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
RUN apt-get update && apt-get install -y migrate
