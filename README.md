# Hashing Service

This is a Service that can hash any type of file, from applications to documents and media. Hashing a file results in a unique and irreversible identifier. If this identifier is registered in an immutable network, as a blockchain ledger, the owner of the file can prove at anytime afterwards that the file was not modified.

This hashing service signs the hash to guarantee that the file was hashed by itselft, and sends it to https://github.com/lacchain/credential-server to write it down in the LACChain blockchain network.   

The Hashing Service and the [Credential Service](https://github.com/lacchain/credential-server) are being used for the notarizing tool named [LACChain Notarizer](http://notarizer.lacchain.net/), that you can use to register and verify the hash any file in the LACChain Blockchain Network for free. For any questions about the tool, you can also read the [LACChain Notarizer FAQ](https://medium.com/@lacchain.official/lacchain-notarizer-faq-6ae3dbb3441e).

## Prerequisites

* Go 1.12+ installation or later
* **GOPATH** environment variable is set correctly
* docker version 17.03 or later

## Package overview

1. **lib** contains most of auxiliar code.
2. **model** contains data models of requests and responses of APIs
3. **main.go** exposes endpoints to consume the hashing service
4. **html** contains the user interface  

## Install

```
$ git clone https://github.com/lacchain/hashing-service

$ export GO111MODULE=on

$ cd hashing-service
$ go build
```

## Run

* Run Hashing Service

```
./hashing-service start
```

* After that, you can go to http://localhost:9000 and start to generate document hash credentials
