#!/bin/bash

export CHARTS_PATH=$1
REGISTRY_WITH_REPO=eu.gcr.io/moneycol
VERSION=0.5.1-app

# Before building copy charts from ../charts

docker build -t $REGISTRY_WITH_REPO/myiac:$VERSION -f Dockerfiles/Dockerfile . \
--build-arg CURRENT_HELM_VERSION=2.16.1 \
--build-arg EXTRA_WORKDIR_ORIG=/workdir/charts \
--build-arg EXTRA_WORKDIR_DEST=/home/app/charts