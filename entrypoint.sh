#!/bin/bash

# Adapted from: 
# https://github.com/kubernetes-retired/contrib/blob/master/startup-script/manage-startup-script.sh

set -o errexit
set -o nounset
set -o pipefail

# Defined in the chart
# CHECK_INTERVAL_SECONDS="30"
# EXEC=(nsenter -t 1 -m -u -i -n -p --)

CHECKPOINT_PATH="/tmp/foo"

echo "====== Startup script Daemonset ====="

do_startup_script() {
  local err=0;
  bash -c  "${STARTUP_SCRIPT}" && err=0 || err=$?
  if [[ ${err} != 0 ]]; then
    echo "!!! startup-script failed! exit code '${err}'" 1>&2
    return 1
  fi

  touch "${CHECKPOINT_PATH}"
  echo "!!! startup-script succeeded!" 1>&2
  return 0
}

while :; do
  do_startup_script
  sleep "${CHECK_INTERVAL_SECONDS}"
done