[![Build Status](https://dev.azure.com/oneeyedelf1/powershell.native/_apis/build/status/KnicKnic.go-windows?branchName=master)](https://dev.azure.com/oneeyedelf1/powershell.native/_build/latest?definitionId=5&branchName=master)

# Goal

To create bindings to allow you a go application to call various Microsoft Windows Server Apis.

## Contents (pkg/...)

Parts of the following cluster api sets
* [Cluster](pkg/cluster/Readme.md)
    * Microsoft Windows Failover Cluster bindings
* [kernel32](pkg/kernel32)
    * LocalAlloc & LocalFree 
* [ntdll](pkg/ntdll)
    * memcpy
