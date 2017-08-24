#!/bin/bash

export HTTPSERVER_LOG_FILE=/letv/app/httpserver.log
export HTTPSERVER_LISTEN_ADDR=:8080
export SERVICE1_NAME=microserver1

while [ 1 ]
do
    pid=`pidof httpserver`
    if [ -z "$pid" ] ; then
        /letv/app/httpserver
    fi
    sleep 5
done
