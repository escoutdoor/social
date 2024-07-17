include .env

up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DSN) goose -dir="./migrations" up

reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DSN) goose -dir="./migrations" reset

db_status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DSN) goose -dir="./migrations" status
