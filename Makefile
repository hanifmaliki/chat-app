run-server:
	go run ./cmd/server/main.go

run-client:
	go run ./cmd/client/main.go

run-docker:
	docker-compose up --build