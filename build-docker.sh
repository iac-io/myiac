#!/bin/bash
# Convert to makefile?
CHARTS_LOCATION=$1
TAG=$2
cp -R $CHARTS_LOCATION .
docker build . -t $TAG
rm -rf charts