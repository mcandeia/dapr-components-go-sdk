# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

ARG COMPONENT=memory

WORKDIR /app

COPY ${COMPONENT}/go.mod ./
COPY ${COMPONENT}/go.sum ./
RUN go mod download

COPY ${COMPONENT}/*.go ./

RUN CGO_ENABLED=0 go build -o /run

## Deploy
FROM gcr.io/distroless/static-debian11

WORKDIR /

COPY --from=build /run /run

CMD ["/run"]