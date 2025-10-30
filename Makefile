run: build
	@./bin/didlydoodash.exe

dev:
	go run cmd/api/main.go

build:
	@go build -o bin/didlydoodash.exe cmd/api/main.go

database:
	@go run cmd/migrate/main.go

drop:
	@go run cmd/drop/main.go