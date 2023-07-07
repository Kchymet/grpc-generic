.PHONY: build
build:
	go build ./...

.PHONY: test
test: build
	go test ./...

.PHONY: mockgen
mockgen:
	go run go.uber.org/mock/mockgen@latest --destination=mocks/mock_service_registrar.go --package=mocks google.golang.org/grpc ServiceRegistrar 
