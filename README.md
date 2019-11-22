# MyIAC

Infrastructure as code. GCP for now.

## Setup

* Create a GCP service account from here: https://cloud.google.com/iam/docs/creating-managing-service-accounts#iam-service-accounts-create-console
* Download the json file and store in your user dir.
* MyIAC will take control from there

## Build with custom executable name

```
go build -o $GOPATH/bin/myiac github.com/dfernandezm/myiac/app
myiac
```

if not, the main package would be used as executable:

```
go install github.com/dfernandezm/myiac/app
app
```

## Golang tutorials

* Structs as classes: https://golangbot.com/structs-instead-of-classes/
* Go packages: https://www.callicoder.com/golang-packages/
* Run commands: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
* Code style: https://golang.org/doc/code.html
* Constructors and initializing structs: https://stackoverflow.com/questions/37135193/how-to-set-default-values-in-go-structs
