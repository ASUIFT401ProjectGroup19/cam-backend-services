TOOLDIR := .bin
PROTOC_GEN_GO_VERSION := v1.27.1
PROTOC_GEN_GO_GRPC_VERSION := v1.1.0
BUF_VERSION := v1.0.0-rc5
export GOBIN=$(abspath $(TOOLDIR))
export PATH := $(GOBIN):$(PATH)

tools:
	mkdir -p $(TOOLDIR)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION) github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@$(BUF_VERSION) github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@$(BUF_VERSION)

.PHONY: lint
lint: tools
	buf lint proto/

.PHONY: generate
generate: lint
	buf generate proto/

build: generate
	docker build -t backend/authentication .

run:
	docker compose up -d

stop:
	docker compose down -v
