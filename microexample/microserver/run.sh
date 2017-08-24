#!/bin/bash

export SERVICE_CONFIG_FILE=/letv/app/config.ini
export SERVER_LOG_FILE=/letv/logs/server.log

while [ 1 ]
do
    pid=`pidof microserver`
    if [ -z "$pid" ] ; then
        /letv/app/microserver
    fi
    sleep 5
done
