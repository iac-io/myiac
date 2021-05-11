#!/bin/bash

if [ -z "${2}" ]; then
  TERRAFORM_PATH='internal/terraform/cluster'
else
  TERRAFORM_PATH=${2}
fi

if [ -z "${3}" ]; then
  HELM_PATH=""
else
  HELM_PATH="-v ${3}:/helm"
fi

# In order to use this image need to build it first
# With this command created image will use current user id for app user inside image:
# docker build -t myiac:latest --build-arg UID=$(id -u) --build-arg GID=$(id -g) Dockerfiles/DevEnv/
# Other available args:
# TERRAFORM_VERSION
# HELM_VERSION
# KUBE_PROMPT_VERSION


docker run -ti --rm \
  --name myiac \
  -v ${PWD}/:/workdir \
  -w /workdir \
  -v ${1}:/account.json \
  -v ${TERRAFORM_PATH}/:/terraform \
  ${HELM_PATH} \
  myiac:latest zsh