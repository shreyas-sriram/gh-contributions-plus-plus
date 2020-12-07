SRCS := $(shell find . -name '*.go')
LINTERS := \
	golang.org/x/lint/golint \
	github.com/kisielk/errcheck \
	honnef.co/go/tools/cmd/staticcheck
APP := gh-contributions-aggregator


.PHONY: deps
deps:
	go get -d -v ./...

.PHONY: updatedeps
updatedeps:
	go get -d -v -u -f ./...

.PHONY: testdeps
testdeps:
	go get -d -v -t ./...
	go get -v $(LINTERS)

.PHONY: updatetestdeps
updatetestdeps:
	go get -d -v -t -u -f ./...
	go get -u -v $(LINTERS)

.PHONY: install
install: deps
	go install ./...

.PHONY: golint
golint: testdeps
	for file in $(SRCS); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

.PHONY: vet
vet: testdeps
	go vet ./...

.PHONY: errcheck
errcheck: testdeps
	errcheck ./...

.PHONY: staticcheck
staticcheck: testdeps
	staticcheck ./...

.PHONY: lint
lint: golint vet errcheck staticcheck

.PHONY: test
test: testdeps lint
	go test -v -race ./test/...

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: build
build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: docker.build
docker.build: build
	docker build . -t $(APP_NAME) 

.PHONY: docker.run
docker.run: build-docker
	docker run -p 3000:3000 $(APP_NAME) 
