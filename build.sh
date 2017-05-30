#!/bin/sh
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o authadm utilities/authadm.go
git add .
git commit -am "Binary build"
git push
