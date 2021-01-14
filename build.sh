#!/bin/bash
## Convert to makefile?
#CHARTS_LOCATION=$1
#cp -R helperScripts $GOPATH/bin
#
## -av --progress
## cp -R $CHARTS_LOCATION $GOPATH/bin
#rsync -a $CHARTS_LOCATION $GOPATH/bin/charts --exclude .git
#
#cp -R terraform $GOPATH/bin
#
#go build -o $GOPATH/bin/myiac github.com/iac-io/myiac/app


env GOOS=darwin GOARCH=amd64 go build cmd/myiac/myiac.go
mv myiac myiac-darwin-amd64
env GOOS=linux GOARCH=amd64 go build cmd/myiac/myiac.go
mv myiac myiac-linux-amd64