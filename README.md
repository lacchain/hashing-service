# Hahing Service

This is a Service that hash any type of file, from applications to documents and media to verify the integrity. 

The Service sign the hash to guarantee that document was hashed by this Service and send it to https://github.com/lacchain/credential-server to written down on blockchain   

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