#
FROM golang:1.16.4-alpine3.13 AS build

COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download

COPY . /src
WORKDIR /src
RUN go build -o dispatcher-worker cmd/dispatcher/main.go

#
FROM alpine:3.13

WORKDIR /app
COPY ./configs/ /app/configs/
COPY --from=build /src/dispatcher /app/

ENV PATH="/app:${PATH}"

ENTRYPOINT ["dispatcher-worker"]