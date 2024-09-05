postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=4713a4cd628778cd1c37a95518f3eaf3 -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root Numeris_DB

dropdb:
	docker exec -it postgres dropdb Numeris_DB

create-migration:
	@if [ -n "$(name)" ]; then \
		migrate create -ext sql -dir migrations -seq $(name); \
	else \
		echo "Error: Missing 'name' variable" >&2; \
		echo "Usage: make create-migration name=\"<NAME_OF_MIGRATION_FILE>\"" >&2; \
		exit 1; \
	fi

migrate-up:
	migrate -path migrations -database "postgresql://root:4713a4cd628778cd1c37a95518f3eaf3@localhost:5432/Numeris_DB?sslmode=disable" -verbose up

migrate-down:
	migrate -path migrations -database "postgresql://root:4713a4cd628778cd1c37a95518f3eaf3@localhost:5432/Numeris_DB?sslmode=disable" -verbose down

mock-user-repo:
	mockgen -package mocked -destination internal/mock/user_repo.go  github.com/zde37/Numeris-Task/internal/repository UserRepository

mock-invoice-repo:
	mockgen -package mocked -destination internal/mock/invoice_repo.go  github.com/zde37/Numeris-Task/internal/repository InvoiceRepository

mock-user-service:
	mockgen -package mocked -destination internal/mock/user_service.go  github.com/zde37/Numeris-Task/internal/service UserService

mock-invoice-service:
	mockgen -package mocked -destination internal/mock/invoice_service.go  github.com/zde37/Numeris-Task/internal/service InvoiceService

test:
	go test -v -cover -short -count=1 ./...
	 
stress:
	go test -v -cover -count=1 ./internal/controller -tags=stress

run:
	go run cmd/main.go

build-run:
	go build -o numeris-task cmd/main.go && ./numeris-task

.PHONY: postgres createdb dropdb createmigration migrateup migratedown mock-user-repo mock-invoice-repo mock-user-service mock-invoice-service test stress server build-run