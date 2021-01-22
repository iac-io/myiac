#!/bin/bash

REGISTRY_WITH_REPO=eu.gcr.io/moneycol
VERSION=app-latest

docker build -t $REGISTRY_WITH_REPO/myiac:$VERSION -f Dockerfiles/Dockerfile . \
--build-arg CURRENT_HELM_VERSION=2.16.1 \
--build-arg EXTRA_WORKDIR_ORIG=/workdir/charts \
--build-arg EXTRA_WORKDIR_DEST=/home/app/charts