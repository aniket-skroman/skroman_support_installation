run:
	go run application.go 

liveserver:
	nodemon --exec go run application.go --signal SIGTERM

migratecreate:
	migrate create -ext sql -dir db/migrations/ -seq alter_client_id_complaint

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" --verbose down

migratefix:
	migrate -path db/migrations/ -database postgres://postgres:root@localhost:5432/postgres?sslmode=disable force 2

sqlc:
	sqlc generate

PHONY:
	run, liveserver, migratecreate, migrateup, migratedown
