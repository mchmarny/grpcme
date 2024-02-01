VERSION    :=$(shell cat .version)
YAML_FILES :=$(shell find . ! -path "./vendor/*" -type f -regex ".*y*ml" -print)
REG_URI    ?=ghcr.io/mchmarny
REPO_NAME  :=$(shell basename $(PWD))

all: help

.PHONY: api
api: ## Generates the API code and documentation
	tools/gen-api

.PHONY: version
version: ## Prints the current version
	@echo $(VERSION)

.PHONY: tidy
tidy: ## Updates the go modules and vendors all dependancies 
	go mod tidy
	go mod vendor

.PHONY: upgrade
upgrade: ## Upgrades all dependancies 
	go get -d -u ./...
	go mod tidy
	go mod vendor

.PHONY: test
test: tidy ## Runs unit tests
	go test -count=1 -race -covermode=atomic -coverprofile=cover.out ./...

.PHONY: lint
lint: lint-go lint-yaml ## Lints the entire project 
	@echo "Completed Go and YAML lints"

.PHONY: lint
lint-go: ## Lints the entire project using go 
	golangci-lint -c .golangci.yaml run

.PHONY: lint-yaml
lint-yaml: ## Runs yamllint on all yaml files (brew install yamllint)
	yamllint -c .yamllint $(YAML_FILES)

.PHONY: vulncheck
vulncheck: ## Checks for soource vulnerabilities
	govulncheck -test ./...

.PHONY: server
server: ## Runs uncpiled version of the server
	go run cmd/server/main.go

.PHONY: ui
ui: ## Runs uncpiled version of the ui client
	go run cmd/ui/main.go

.PHONY: headless
headless: ## Runs uncpiled version of the headless client
	go run cmd/headless/main.go

.PHONY: image
image: ## Builds the server images
	@echo "Building server image..."
	KO_DOCKER_REPO=$(REG_URI)/$(REPO_NAME)-server \
    GOFLAGS="-ldflags=-X=main.version=$(VERSION)" \
    ko build cmd/server/main.go --image-refs .digest --bare --tags $(VERSION),latest
	@echo "Building ui image..."
	KO_DOCKER_REPO=$(REG_URI)/$(REPO_NAME)-ui \
    GOFLAGS="-ldflags=-X=main.version=$(VERSION)" \
    ko build cmd/ui/main.go --image-refs .digest --bare --tags $(VERSION),latest
	@echo "Building headless image..."
	KO_DOCKER_REPO=$(REG_URI)/$(REPO_NAME)-headless \
    GOFLAGS="-ldflags=-X=main.version=$(VERSION)" \
    ko build cmd/headless/main.go --image-refs .digest --bare --tags $(VERSION),latest

.PHONY: tag
tag: ## Creates release tag 
	git tag -s -m "version bump to $(VERSION)" $(VERSION)
	git push origin $(VERSION)

.PHONY: tagless
tagless: ## Delete the current release tag 
	git tag -d $(VERSION)
	git push --delete origin $(VERSION)

.PHONY: clean
clean: ## Cleans bin and temp directories
	go clean
	rm -fr ./vendor
	rm -fr ./bin

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
