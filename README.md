# Purpose
This repository is meant to stimulate ideas on how to store public data records on Hyperledger Sawtooth. 
It's a slightly modified copy of xo_go transaction processor, part of github.com/hyperledger/sawtooth-sdk-go

# RFC
Please read and review the RFC located [here](docs/RFC.md).

# System Requirements
1. OS Packages
    ```
    sudo apt-get update
    sudo apt-get -y upgrade
    sudo apt install -y zip curl python3 python3-pip pkg-config
    ```

2. Install protobuf compilers (make sure to get a 3.x.x version from [here](https://github.com/protocolbuffers/protobuf/releases))
    ```
    curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
    unzip protoc-3.7.1-linux-x86_64.zip -d protoc3
    sudo mv protoc3/bin/* /usr/local/bin/
    sudo mv protoc3/include/* /usr/local/include/
    ```

3. Install python's grpcio-tools library
    ```
    sudo su - 
    python3 -m pip install grpcio-tools
    ```

# Installation

Please see [Packaging As A Service](docs/PackageAsService.md)

# Client Usage
Please see [Usage](docs/Usage.md)
---

# Contributing: Development Requirements

1. Install golang version 1.12
    ```
    wget https://dl.google.com/go/go1.12.2.linux-amd64.tar.gz

    sudo tar -xvf go1.12.2.linux-amd64.tar.gz
    sudo mv go /usr/bin
    ```

    Add the following to your ~/.profile
    ```
    export GOROOT="/usr/bin/go"
    export GOPATH="$HOME/go"

    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
    ```

2. Install go dependencies
    ```
    go install github.com/golang/mock/mockgen && \
    go get -u google.golang.org/grpc \
        github.com/golang/protobuf/protoc-gen-go \
        github.com/satori/go.uuid \
        github.com/pebbe/zmq4 \
        github.com/golang/mock/gomock \
        github.com/hyperledger/sawtooth-sdk-go \
        github.com/jessevdk/go-flags \
        github.com/stretchr/testify/mock \
        github.com/btcsuite/btcd/btcec \
        gopkg.in/yaml.v2

    cd $GOPATH/src/github.com/hyperledger/sawtooth-sdk-go && \
        go generate
    ```

# Manual Testing Results
Read @ [here](docs/TestCases.md)

# Enhancements
Please review the Issues for enhancements to this project.