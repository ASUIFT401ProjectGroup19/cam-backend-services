TOOLDIR := .bin

export GOBIN=$(abspath $(TOOLDIR))
export PATH := $(GOBIN):$(PATH)

tools:
	mkdir -p $(TOOLDIR)
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

build:
	CGO_ENABLED=0 GOOS=linux go build -o .bin/authserver ./cmd/authentication/

image:
	docker build -t backend/authentication .

debug:
	docker compose up debug-auth -d

run:
	docker compose up backend-auth -d --build

stop:
	docker compose down -v
