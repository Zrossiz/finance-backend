#!/bin/sh
docker compose --env-file .env -f deployment/docker-compose.yml up --build -d
go run cmd/migrator/migrator.go -h localhost