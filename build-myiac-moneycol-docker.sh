#!/bin/bash

export CHARTS_PATH=$1
REGISTRY_WITH_REPO=eu.gcr.io/moneycol
VERSION=0.5.0-app

#if [[ -n "$CHARTS_PATH" ]]; then
#  rsync -rv --exclude=.git "$CHARTS_PATH" .
#fi

# Before building copy charts from ../charts

docker build -t $REGISTRY_WITH_REPO/myiac:$VERSION -f Dockerfiles/Dockerfile . \
--build-arg CURRENT_HELM_VERSION=2.16.1 \
--build-arg EXTRA_WORKDIR_ORIG=/workdir/charts \
--build-arg EXTRA_WORKDIR_DEST=/home/app/charts