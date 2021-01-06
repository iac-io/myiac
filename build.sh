#!/bin/bash
# Convert to makefile?
CHARTS_LOCATION=$1
cp -R helperScripts $GOPATH/bin

# -av --progress
# cp -R $CHARTS_LOCATION $GOPATH/bin
rsync -a $CHARTS_LOCATION $GOPATH/bin/charts --exclude .git

cp -R terraform $GOPATH/bin

go build -o $GOPATH/bin/myiac github.com/iac-io/myiac/app