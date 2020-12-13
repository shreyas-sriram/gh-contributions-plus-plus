SRCS := $(shell find . -name '*.go')
LINTERS := \
	golang.org/x/lint/golint \
	honnef.co/go/tools/cmd/staticcheck
APP_NAME := gh-contributions-aggregator

BANNER:=\
    "\n"\
		"/**\n"\
    " * @project       $(APP_NAME)\n"\
    " */\n"\
    "\n"

## deps				: Download dependencies
.PHONY: deps
deps:
	go get -d -v ./...

## testdeps			: Download test dependencies
.PHONY: testdeps
testdeps:
	go get -d -v -t ./...
	go get -v $(LINTERS)

## install			: Install dependencies
.PHONY: install
install: deps
	go install ./...

## golint				: Run golint on all *.go files
.PHONY: golint
golint: 
	for file in $(SRCS); do \
		golint $${file}; \
		if [ -n "$$(golint $${file})" ]; then \
			exit 1; \
		fi; \
	done

## vet				: Run vet
.PHONY: vet
vet:
	go vet ./...

## staticcheck			: Run staticcheck
.PHONY: staticcheck
staticcheck:
	staticcheck ./...

## lint				: Run golint, vet and staticcheck
.PHONY: lint
lint: golint vet staticcheck

## test				: Run tests
.PHONY: test
test:
	ENV=test go test -v -race -coverprofile=coverage.out ./...

## test.coverage			: Show test coverage report
.PHONY: test.coverage
test.coverage:
	go tool cover -html=coverage.out

## clean				: Clean application objects
.PHONY: clean
clean:
	go clean -i ./...

## build.linux			: Build application for Linux runtime
.PHONY: build.linux
build.linux:
	env GOOS=linux go build -ldflags="-s -w" -o bin/$(APP_NAME) cmd/api/main.go

## build.mac			: Build application for Mac runtime
.PHONY: build.mac
build.mac:
	env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(APP_NAME) cmd/api/main.go

## sls.build			: Build application for Linux runtime
.PHONY: sls.build
sls.build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/$(APP_NAME) cmd/api/main.go

## sls.deploy			: Deploy application as serverless
.PHONY: sls.deploy
sls.deploy:
	sls deploy -v

## sls				: Build and deploy application as serverless
.PHONY: sls
sls: build.linux deploy.sls

## server.start			: Run application as server from runtime binary
.PHONY: server.start
server.start:
	DEPLOY=server ./bin/$(APP_NAME)

## run				: Run application from main.go
.PHONY: run
run:
	go run cmd/api/main.go

## docker.build			: Build application as docker
.PHONY: docker.build
docker.build:
	docker build . -t $(APP_NAME) 

## docker.run			: Run application in docker
.PHONY: docker.run
docker.run:
	docker run -p 3000:3000 $(APP_NAME)

## help				: Show all available make targets
.PHONY : help
help : 
	@echo $(BANNER)
	@echo \\tMake targets
	@echo -----------------------------
	@sed -n 's/^##//p' Makefile | sort
