# Dockerfiles

## Root Dockerfile

This Dockerfile allows packaging `myiac` tool in Docker container so it can be deployed or run as standalone app.
This enables its usage within CI/CD systems than run builds/packers as Docker containers. It also makes `myiac` a 
regular Go app that can perform a number of Infrastructure related checks/actions within Kubernetes clusters or other
systems where Docker can run.

Moneycol note:
- This `Dockefile` should be run through the script `build-myiac-moneycol-docker.sh`

## DevEnv

TBA