#!/bin/bash

# protoc 3
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
unzip protoc-3.7.1-linux-x86_64.zip -d protoc3
sudo mv protoc3/bin/* /usr/local/bin/
sudo mv protoc3/include/* /usr/local/include/

# pip3
curl -sS https://bootstrap.pypa.io/get-pip.py | sudo python3

# grpcio tools
pip3 install --user grpcio-tools

# go dependencies
go get -u google.golang.org/grpc \
        golang.org/x/tools/go/packages \
        github.com/golang/protobuf/protoc-gen-go \
        github.com/satori/go.uuid \
        github.com/pebbe/zmq4 \
        github.com/golang/mock/gomock \
        github.com/hyperledger/sawtooth-sdk-go \
        github.com/jessevdk/go-flags \
        github.com/stretchr/testify/mock \
        github.com/btcsuite/btcd/btcec \
        gopkg.in/yaml.v2 \
        github.com/labstack/echo \
	github.com/labstack/echo/middleware

go install github.com/golang/mock/mockgen

cd $GOPATH/src/github.com/hyperledger/sawtooth-sdk-go && \
        go generate