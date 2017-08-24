consul node:10.135.28.154,10.112.34.55,10.185.30.76

注册服务
curl -X PUT "http://10.135.28.154:8500/v1/agent/service/register" -d '
{
    "ID":"micro-test1",
    "Name":"micro-test1",
    "Tags":[
        "t-789caa564a2acacf"
    ],
    "Port":8687,
    "Address":"127.0.0.1",
    "Check":{
        "TCP":"127.0.0.1:8687",
        "Interval": "20s",
        "DeregisterCriticalServiceAfter":"1m5s",
        "Status":"passing"
    },
    "Checks":null
}' -v


获取检测点
curl http://10.135.28.154:8500/v1/agent/checks|jq .

获取单个consul节点的健康检测信息
curl "http://10.135.28.154:8500/v1/health/node/10.135.28.154" |jq .

获取单个服务的健康信息
curl "http://10.135.28.154:8500/v1/health/checks/micro-test" |jq .

获取单个服务的健康节点
curl "http://10.135.28.154:8500/v1/health/service/micro-test" |jq .