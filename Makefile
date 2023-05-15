nats:
	nats-streaming-server

publish:
	go run internal/services/nats-streaming/publish/publisher.go

start:
	go run cmd/app/main.go