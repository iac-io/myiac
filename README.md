# MyIAC

Infrastructure as code. GCP for now.

## Setup

* Create a GCP service account with admin privileges from here: https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
* Download the json file and store in your user directory (`$HOME/account.json`)
* Install Go [here](https://golang.org/dl/)
* Check `$GOPATH` is pointing to `~/go` (Mac OS X), edit your `bashrc` or `zshrc` to ensure it's the case. 
* Ensure the following structure exists (Go Workspace)
  - `mkdir $GOPATH/src`
  - `mkdir $GOPATH/pkg`
  - `mkdir $GOPATH/bin`
  - `mkdir $GOPATH/src/github.com/dfernandezm/myiac`
* Add also `$GOPATH/go/bin` to your path
* Clone the project inside the `$GOPATH/src/github.com/dfernandezm/myiac`
* If deployment is required, copy a `charts` folder into `$GOPATH/bin`. This folder should follow the structure `charts/{appName}`. Inside there a typical structure for a Helm chart should be present (templates, values.yaml...)

## Go build with custom executable name

```
# inside myiac folder
$ go get ./...
$ go build -o $GOPATH/bin/myiac github.com/dfernandezm/myiac/app
$ myiac
```

if not, the main package would be used as executable:

```
go install github.com/dfernandezm/myiac/app
app
```

## Get usage help

```
myiac help
myiac [subcommand] help
```

## Setting up SSL in traefik

https://github.com/dfernandezm/myiac/blob/e0cbdde19ed9c4b8da750481e175e936c66d113c/kubernetes/cluster/README.md


## Golang tutorials

* Structs as classes: https://golangbot.com/structs-instead-of-classes/
* Go packages: https://www.callicoder.com/golang-packages/
* Run commands: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
* Code style: https://golang.org/doc/code.html
* Constructors and initializing structs: https://stackoverflow.com/questions/37135193/how-to-set-default-values-in-go-structs
