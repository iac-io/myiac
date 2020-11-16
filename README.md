# MyIAC

Infrastructure as code. GCP for now.

## Setup

* Create a GCP service account with admin privileges from here: https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
* Download the json file and store in your user directory (`$HOME/account.json`)
* Install Go [here](https://golang.org/dl/), it should be Go `v1.13` at least
* Clone the project
* If deployments are going to be run, export the environment variable `CHARTS_PATH` pointing to a folder that contains Helm Charts. This folder should follow the structure `charts/{appName}`. Inside, the typical structure for a Helm chart should be present (templates, values.yaml...)
```
export CHARTS_PATH=/path/to/charts
```

## Build executable

```
$ go build cmd/myiac/myiac.go
$ ./myiac

# Move it to PATH folder to run command from anywhere
# $ mv ./myiac /usr/local/bin
```

## Get usage help

```
myiac help
myiac [subcommand] help
```

## Known issues & work in progress

Currently only the following subcommands are working as expected in the `master` branch (the project is undergoing a massive cleanup & refactor):

- `setupEnvironment`
- `deploy`
- `crypt`
- `createSecret`
- `dockerBuild`
- `createCert` (basic usage)

The cluster needs to be created beforehand as the `createCluster` and `destroyCluster` subcommands are currently not working.

## Golang tutorials

* Structs as classes: https://golangbot.com/structs-instead-of-classes/
* Go packages: https://www.callicoder.com/golang-packages/
* Run commands: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
* Code style: https://golang.org/doc/code.html
* Constructors and initializing structs: https://stackoverflow.com/questions/37135193/how-to-set-default-values-in-go-structs

## Setting up SSL in traefik

https://github.com/dfernandezm/myiac/blob/e0cbdde19ed9c4b8da750481e175e936c66d113c/kubernetes/cluster/README.md

