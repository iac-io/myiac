#!/bin/bash

export TF_VAR_cluster_state_bucket="qumu-dev-tf-state-dev"
export TF_VAR_state_bucket_prefix="terraform/state/cluster"
export TF_VAR_credentials_file="/Users/dfernandez/qumu-dev_account.json"

# https://cloud.google.com/docs/authentication/production
export GOOGLE_APPLICATION_CREDENTIALS=/Users/dfernandez/qumu-dev_account.json
gsutil mb -p qumu-dev gs://$TF_VAR_cluster_state_bucket
terraform init -backend-config "bucket=$TF_VAR_cluster_state_bucket" \
-backend-config "prefix=$TF_VAR_state_bucket_prefix" -backend-config="credentials=$TF_VAR_credentials_file"
