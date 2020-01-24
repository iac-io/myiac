#!/bin/bash
# Convert to makefile?
CHARTS_LOCATION=$1
cp -R helperScripts $GOPATH/bin
cp -R $CHARTS_LOCATION $GOPATH/bin
cp -R terraform $GOPATH/bin

go build -o $GOPATH/bin/myiac github.com/dfernandezm/myiac/app