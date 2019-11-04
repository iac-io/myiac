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

