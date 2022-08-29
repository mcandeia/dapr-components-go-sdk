# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-buster AS build

ARG COMPONENT=memory

WORKDIR /sdk

## Copying sdk files
COPY internal internal
COPY *.go ./
COPY state state
COPY go.mod ./
COPY go.sum ./

WORKDIR /app

COPY examples/${COMPONENT}/go.mod ./
COPY examples/${COMPONENT}/go.sum ./
RUN go mod edit -replace github.com/mcandeia/dapr-components-go-sdk=/sdk/
RUN go mod download

COPY examples/${COMPONENT}/*.go ./
COPY examples/umask_entrypoint_wrapper.sh ./

RUN CGO_ENABLED=0 go build -o component

CMD ["./umask_entrypoint_wrapper.sh", "./component"]

## Deploy
#FROM gcr.io/distroless/static-debian11

#WORKDIR /

#COPY --from=build /component /component

#ENTRYPOINT ["/component"]