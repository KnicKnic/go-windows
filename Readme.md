[![Build Status](https://dev.azure.com/oneeyedelf1/wasm-imagemagick/_apis/build/status/KnicKnic.go-failovercluster-api?branchName=master)](https://dev.azure.com/oneeyedelf1/wasm-imagemagick/_build/latest?definitionId=4&branchName=master)

# Goal

To create bindings to allow you a go application to call Microsoft Windows Server Failover Cluster Api.

Currently uses syscall to wrap the c clusapi.dll and resutils.dll code.

## Completed

Parts of the following cluster api sets
1. Cluster
1. Resource
1. Registry
1. Crypto

## TODO

* add comments for public functions
* Write more tests
* complete more wrappers for functions
