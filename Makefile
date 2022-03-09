TOOLDIR := .bin

export GOBIN=$(abspath $(TOOLDIR))
export PATH := $(GOBIN):$(PATH)

tools:
	mkdir -p $(TOOLDIR)
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

sync-mods:
	go mod tidy
	go mod vendor

build: sync-mods
	CGO_ENABLED=0 GOOS=linux go build -o .bin/server ./cmd/

image: sync-mods
	docker build -t backend/server .

debug: sync-mods
	docker compose up debug -d

run: sync-mods
	docker compose up backend -d --build

stop:
	docker compose down -v
