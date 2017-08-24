#!/bin/bash
jq_url="http://static.scloud.letv.cn/wcy/jq"
consul_addr=(
    http://10.135.28.154:8500
    http://10.185.30.76:8500
    http://10.112.34.55:8500
    )

# store parameter
serivce_name=""

function log() {
    echo "`date` $1 $2 $3 $4"
}

function usage() {
    echo "$0 usage:"
    echo "  listservices: list all service"
    echo "  serviceinfo: list service node info"
    echo "    eg: serviceinfo  servicename"
}

function check_dev() {
    flag=`which jq |grep "no jq"`
    if [ ! -z "$flag" ] ; then
        log "not found jq,now download it from:$jq_url"
        curl "$jq_url" -o /usr/bin/jq
        chmod +x /usr/bin/jq
    fi
}

function list_services() {
    for addr in ${consul_addr[@]}
    do
        res=`curl -s $addr/v1/agent/services | jq . | grep "Service" | awk -F'"' '{print $4}' | sort -u`
        if [ ! -z "$res" ] ; then
            echo "$res"
            break
        fi
    done
}

function list_service_info() {
    if [ -z "$serivce_name" ] ; then
        echo "lost service name when list service nodes"
        return
    fi
    for addr in ${consul_addr[@]}
    do
        res=`curl -s $addr/v1/health/service/$serivce_name | jq . | grep Success | awk '{print $4}'|awk -F: '{print $1 ":" $2}'|sort -u`
        if [ ! -z "$res" ] ; then
            echo "$res"
            break
        fi
    done
}

check_dev
case "$1" in
    "listservices")
        list_services
    ;;
    "serviceinfo")
        serivce_name=$2
        list_service_info
    ;;
    *)
        usage
    ;;
esac