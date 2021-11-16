#!/usr/bin/env bash
if [ ! -f ./go.mod ]; then
  go mod init
fi
go mod vendor
swag init
env GOOS=linux GOARCH=amd64 go build -v main.go
docker-compose -f docker-compose.yml build