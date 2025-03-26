# syntax=docker/dockerfile:1

### buid stage
FROM golang:1.23-alpine AS build

WORKDIR /usr/src/bot

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading 
# them in subsequent builds if they change
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/bot

### deploy stage
FROM alpine:3.21

WORKDIR /usr/local/bin

COPY --from=build /usr/local/bin/bot .

CMD ["bot"]
