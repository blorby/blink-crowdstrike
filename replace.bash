#!/usr/bin/env bash

VAR="$1"

for file in "build/Dockerfile" "build/Jenkinsfile" "cmd/main/main.go" "go.mod"; do sed -i "" "s/EXAMPLE/${VAR}/g" $file; done