.PHONY: all build run start

LOCAL_PORT:=5436
DOCKER_PORT:=5432
PASSWORD:=qwerty
USERNAME:=postgres
HOST:=localhost
SSLMODE:=disable
DATABASE:=postgres

all: build run start

build:
	docker build -t postgres_flood_control .

run:
	docker run --name bd_flood_control -p $(LOCAL_PORT):$(DOCKER_PORT) -e POSTGRES_PASSWORD=$(PASSWORD) -d postgres_flood_control

start:
	go run cmd/main.go
