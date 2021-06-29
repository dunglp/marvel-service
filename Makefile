GOOS ?= "linux"
GOARCH ?= "amd64"
IMAGE_NAME = "dunglp-xendit-assessment"

.PHONY : run
run : deps lint unit-test build-docker run-docker

.PHONY: deps
deps:
	go mod tidy -v
	go mod vendor -v

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./svc ./cmd

.PHONY: gen
gen:
	mockery --all --dir ./pkg/service --output ./pkg/service/mocks
	mockery --all --dir ./pkg/domain --output ./pkg/domain/mocks
	mockery --all --dir ./pkg/marvel --output ./pkg/marvel/mocks

.PHONY: lint
lint:
	test -s ./golangci-lint || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b "$(PWD)" v1.31.0
	./golangci-lint run -v --timeout 10m ./...

.PHONY: unit-test
unit-test:
	go test -v -race -mod=vendor -cover -coverprofile=cover.out $$(go list ./... | grep -v /test)

.PHONY: view-unit-test-results
view-unit-test-results:
	go tool cover -html=cover.out

.PHONY: check-swagger
check-swagger: ## Check if go-swagger tool is in place
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

.PHONY: generate-swagger
generate-swagger:
	 swagger generate spec -o ./swagger.yaml --scan-models

.PHONY: serve-swagger ## Serve the generated swagger specs
serve-swagger: check-swagger
	swagger serve -F=swagger --port 3003 swagger.yaml

.PHONY: build-docker
build-docker:
	docker build -t $(IMAGE_NAME) .

.PHONY: run-docker
run-docker:
	docker run --rm --name $(IMAGE_NAME) -p 8080:8080 $(IMAGE_NAME)
