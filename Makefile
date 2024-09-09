.PHONY: test build build-docker
test:
	go test -v ./...

build:
	go build cmd/api/main.go

run:
	PORT_FILE=data/ports.json go run cmd/api/main.go

run-jaeger:
	docker run -d --name jaeger \
		-e COLLECTOR_OTLP_ENABLED=true \
		-e COLLECTOR_OTLP_HTTP_PORT=4318 \
		-p 16686:16686 \
		-p 4318:4318 \
		-p 14250:14250 \
		-p 14268:14268 \
		jaegertracing/all-in-one:latest

build-docker:
	docker build --platform linux/amd64  -t willejs/port-service:latest .

deploy:
	helm install port-service helm/port-service