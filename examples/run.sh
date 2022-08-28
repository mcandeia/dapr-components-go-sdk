#!/bin/bash

COMPONENT=${1:-memory} docker-compose build
COMPONENT=${1:-memory} docker-compose up -d