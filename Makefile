.PHONY: test build proto

config:
	go mod tidy

test-local: start-env test stop-env

test:
	go test -timeout 30s ./...

test-race:
	go test -v -timeout 30s -race ./...

bench:
	go test -v -tags bench -benchmem -run '^$$' -bench Bench ./...

coverage:
	go test -tags integration -v -coverprofile=.coverage.out -timeout 30s ./...
	go tool cover -func=.coverage.out

start-env:
	docker-compose up -d psql sqs

stop-env:
	docker-compose down --rmi local