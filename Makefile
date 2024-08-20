.PHONY: test build build-docker
test:
	go test -v ./...

build:
	go build cmd/api/main.go

build-docker:
	docker build --platform linux/amd64  -t willejs/port-service:latest .
