.DEFAULT_GOAL := help

# HOST is only used for API specs generation
HOST ?= localhost:8191

depends: ## Install & build dependencies
	go get ./...
	go build ./...
	go mod tidy

run:
	go run cmd/api/main.go

mod:
	go mod tidy && go mod vendor

specs: ## Generate swagger specs
	HOST=$(HOST) sh scripts/specs-gen.sh

specs-linux: ## Generate swagger specs linux
	HOST=$(HOST) sh scripts/specs-gen-linux.sh

build-image: specs
	sh ./build_image.sh