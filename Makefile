TOOLDIR := .bin
PROTOC_GEN_GO_VERSION := v1.27.1
PROTOC_GEN_GO_GRPC_VERSION := v1.1.0
export GOBIN=$(abspath $(TOOLDIR))
export PATH := $(GOBIN):$(PATH)

prototool:
	mkdir -p $(TOOLDIR)
	cd $(TOOLDIR); go get google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	cd $(TOOLDIR); go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION)
	cd $(TOOLDIR); curl -L -o p.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.18.1/protoc-3.18.1-linux-x86_64.zip
	cd $(TOOLDIR); unzip -o p.zip

.PHONY: generate
generate: prototool
	mkdir -p pkg/gen/go
	.bin/bin/protoc --go_out=pkg/gen/go --go-grpc_out=pkg/gen/go -I proto/ proto/*

