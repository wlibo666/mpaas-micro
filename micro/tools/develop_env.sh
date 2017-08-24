#!/bin/bash
curdir=`pwd`

function log() {
    echo "`date` $1 $2 $3 $4"
}

function init() {
    export PATH=$PATH:$HOME/bin:$GOROOT/bin:$GOPATH/bin
}

function install_micro() {
    pkgs=(
        github.com/micro/go-micro
        github.com/micro/go-plugins
        github.com/micro/cli
        github.com/micro/examples
        github.com/micro/go-api
        github.com/micro/go-bot
        github.com/micro/go-grpc
        github.com/micro/go-log
        github.com/micro/go-os
        github.com/micro/go-run
        github.com/micro/hipchat
        github.com/micro/mdns
        github.com/micro/micro
        github.com/micro/misc
        github.com/micro/protobuf
        github.com/go-ini/ini
        github.com/go-kit/kit/log
        github.com/go-kit/kit/log/level
        )
    for pkg in ${pkgs[@]}
    do
        log "go get $pkg"
        go get $pkg
    done    

    mv $GOPATH/src/github.com/micro/go-micro $GOPATH/src/github.com/micro/go-micro.bak
    mv $GOPATH/src/github.com/micro/go-plugins $GOPATH/src/github.com/micro/go-plugins.bak

    cd $GOPATH/src/github.com/micro/
    git clone https://github.com/wlibo666/go-micro.git
    git clone https://github.com/wlibo666/go-plugins.git
}

function add_vendor() {
    #log "download vendor: http://static.scloud.letv.cn/wcy/micro_service_vendor.tar.gz"
    curl "http://static.scloud.letv.cn/wcy/micro_service_vendor.tar.gz" -o $curdir/micro_service_vendor.tar.gz
    cd $curdir
    rm -rf vendor/
    tar xf micro_service_vendor.tar.gz
    rm -f micro_service_vendor.tar.gz

    newroot=${curdir##$GOPATH/src/}
    #log "exec sed -i s%github.com/wlibo666/mpaas-micro/microexample%$newroot%g vendor/vendor.json"
    sed -i "s%github.com/wlibo666/mpaas-micro/microexample%$newroot%g" vendor/vendor.json
}

function install_proto3() {
    cd $GOPATH
    curl "http://static.scloud.letv.cn/wcy/protoc-3.3.0-linux-x86_64.zip" -o $GOPATH/protoc-3.3.0-linux-x86_64.zip
    unzip protoc-3.3.0-linux-x86_64.zip
    rm $GOPATH/protoc-3.3.0-linux-x86_64.zip
}

function install_protobuf() {
    # if not install protoc,install it
    if [ ! -f $GOPATH/bin/protoc ] ; then
        install_proto3
    else
        # if version is not 3.3.x,reinstall it
        version=`$GOPATH/bin/protoc --version | grep "3.3"`
        if [ -z "$version" ] ; then
            install_proto3
        fi
    fi

    filesize=`stat $GOPATH/bin/protoc-gen-go | grep Size| awk '{print $2}'`
    if [ "$filesize" != "4192585" ] ; then
        curl "http://static.scloud.letv.cn/wcy/protoc-gen-go" -o $GOPATH/bin/protoc-gen-go
    fi
}

init
#install_micro
add_vendor
install_protobuf
