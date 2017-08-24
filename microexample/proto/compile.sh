#!/bin/bash
# go get github.com/micro/protobuf/proto
# go get github.com/micro/protobuf/protoc-gen-go
# cd $GOPATH/src/github.com/micro/protobuf/protoc-gen-go ; go install ; cd - ;
curdir=`pwd`
export PATH=$PATH:$GOPATH/bin
$GOPATH/bin/protoc -I$GOPATH/src --go_out=plugins=micro:$GOPATH/src \
    $curdir/calculator.proto