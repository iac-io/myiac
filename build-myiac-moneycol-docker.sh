#!/bin/bash

export CHARTS_PATH=$1
REGISTRY_WITH_REPO=eu.gcr.io/moneycol
VERSION=0.6.0-app1

# Before building copy charts from ../charts
docker build -t $REGISTRY_WITH_REPO/myiac:$VERSION -f Dockerfiles/Dockerfile . \
--build-arg CURRENT_HELM_VERSION=3.7.1 \
--build-arg EXTRA_WORKDIR_ORIG=/workdir/charts-dns \
--build-arg EXTRA_WORKDIR_DEST=/home/app/charts