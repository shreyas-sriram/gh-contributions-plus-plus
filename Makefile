SRCS := $(shell find . -name '*.go')
LINTERS := \
	golang.org/x/lint/golint \
	honnef.co/go/tools/cmd/staticcheck
APP_NAME := gh-contributions-aggregator


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
golint: 
	for file in $(SRCS); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: lint
lint: golint vet staticcheck

.PHONY: test
test:
	ENV=test go test -v -race -coverprofile=coverage.out ./...

.PHONY: test.coverage
test.coverage:
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	go clean -i ./...

.PHONY: build.linux
build.linux:
	env GOOS=linux go build -ldflags="-s -w" -o bin/$(APP_NAME) cmd/api/main.go

.PHONY: build.mac
build.mac:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(APP_NAME) cmd/api/main.go

.PHONY: deploy.sls
deploy.sls:
	sls deploy -v

.PHONY: sls
sls: build.linux deploy.sls

.PHONY: start
start:
	DEPLOY=server ./bin/$(APP_NAME)

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: docker.build
docker.build:
	docker build . -t $(APP_NAME) 

.PHONY: docker.run
docker.run:
	docker run -p 3000:3000 $(APP_NAME) 
