FROM golang:1.15.3-alpine AS build_base

ENV CGO_ENABLED=1
ENV GO111MODULE=on
RUN apk add --no-cache git gcc g++

# Set the Current Working Directory inside the container
WORKDIR /src

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/app ./cmd/server/main.go

# Start fresh from a smaller image
FROM alpine:3.12
RUN apk add ca-certificates

WORKDIR /app

COPY --from=build_base /src/out/app /app/gh-contributions-plus-plus
COPY --from=build_base /src/data /app/data

RUN chmod +x gh-contributions-plus-plus

# This container exposes port 3000 to the outside world
EXPOSE 3000

# Run the binary program produced by `go install`
ENTRYPOINT ./gh-contributions-plus-plus
