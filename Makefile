include .env

run:
	@docker-compose up

test:
	@go test ./... -v

migrations_up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(POSTGRES_URL_LOCALHOST) goose -dir="./migrations" up

migrations_reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(POSTGRES_URL_LOCALHOST) goose -dir="./migrations" reset

db_status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(POSTGRES_URL_LOCALHOST) goose -dir="./migrations" status
