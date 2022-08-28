#!/bin/bash

## cleaning socket volume
docker container rm dapr-pluggable-component
docker container rm daprd-pluggable-component
docker volume rm examples_socket


COMPONENT=${1:-memory} docker-compose build
COMPONENT=${1:-memory} docker-compose up