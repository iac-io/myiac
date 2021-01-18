#!/bin/bash

if [ -z "${2}" ]; then
  TERRAFORM_PATH='internal/terraform/cluster'
else
  TERRAFORM_PATH=${2}
fi

docker run -ti --rm \
  --name myiac \
  -v ${PWD}/:/workdir \
  -w /workdir \
  -v ${1}:/account.json \
  -v ${TERRAFORM_PATH}/:/terraform \
  myiac:latest zsh