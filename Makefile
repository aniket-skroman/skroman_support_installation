run:
	go run application.go 

liveserver:
	nodemon --exec go run application.go --signal SIGTERM

migratecreate:
	migrate create -ext sql -dir db/migrations/ -seq change_device_id_type

migrateup:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" --verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://postgres:root@localhost:5432/postgres?sslmode=disable" --verbose down

migratefix:
	migrate -path db/migrations/ -database postgres://postgres:root@localhost:5432/postgres?sslmode=disable force 1

migrateup_maindb:
	migrate -path db/migrations/ -database  "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_client_complaints" --verbose up

migratedown_maindb:
	migrate -path db/migrations/ -database  "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_client_complaints" down

migratefix_maindb:
	migrate -path db/migrations/ -database  "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_client_complaints" force 3

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

PHONY:
	run, liveserver, migratecreate, migrateup, migratedown
