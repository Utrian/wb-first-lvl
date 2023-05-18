nats:
	nats-streaming-server

publish:
	go run internal/services/nats-streaming/publish/publisher.go

run:
	go run cmd/app/main.go

build:
	go build cmd/app/main.go

main:
	./main