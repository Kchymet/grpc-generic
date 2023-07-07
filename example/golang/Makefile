
.PHONY: run
run: build
	go run ./cmd/...

.PHONY: build
buid: gen
	go build ./...

.PHONY: test
test: build
	go test ./...

.PHONY: gen
gen: gogen mockgen bufgen

.PHONY: gogen
gogen:
	go generate ./...

.PHONY: bufgen
bufgen:
	buf generate --template api/buf.gen.yaml --path api

.PHONY: mockgen
mockgen:
	go run go.uber.org/mock/mockgen@latest -- source=github.com/grpc/grpc-go/server.go --destination=mocks/mock_grpc_server.go --package=mocks
