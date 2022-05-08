#!/bin/sh

set -e

go build -ldflags=-w -o ./dist/service
docker-compose up --build app
