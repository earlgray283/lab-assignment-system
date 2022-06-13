FROM golang:1.18 AS build

WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN \
    GOOS=linux \
    GOARCH=amd64 \
    go build -mod=readonly -v -o app

FROM --platform=linux/amd64 ubuntu:22.04

RUN apt-get update && apt-get -y --no-install-recommends \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=build /work/app .

ENTRYPOINT [ "/app" ]
