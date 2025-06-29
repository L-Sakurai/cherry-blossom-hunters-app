#!/bin/bash
cd /workspace/cherry-blossom-hunters-app/ && go mod tidy
go run cmd/server/main.go