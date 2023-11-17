# syntax=docker/dockerfile:1

FROM golang:1.21 AS build-stage

WORKDIR /usr/bin/

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o opends cmd/opends/main.go

FROM debian:12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /usr/bin/opends /usr/bin/opends

EXPOSE 13000

ENTRYPOINT ["/usr/bin/opends", "-port", "13000"]

