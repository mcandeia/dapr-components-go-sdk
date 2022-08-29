#!/bin/bash

COMPONENT=${1:-memory} docker-compose build --no-cache
COMPONENT=${1:-memory} docker-compose up