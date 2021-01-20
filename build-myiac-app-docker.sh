#!/bin/bash

REGISTRY_WITH_REPO=eu.gcr.io/moneycol
VERSION=app-latest

docker build -t $REGISTRY_WITH_REPO/myiac:$VERSION -f Dockerfiles/myiac-app/Dockerfile .