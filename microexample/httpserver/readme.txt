应用名称: testmicroserver2
介绍: 一个HTTP服务，提供加法和乘法操作，在后端通过调用微服务 testmicroserver1获得结果

测试接口:
curl "http://calc.scloud.letv.cn/v1/add?num1=123&num2=456" -x 10.176.28.229:80 -v
curl "http://calc.scloud.letv.cn/v1/mul?num1=123&num2=456" -x 10.176.28.229:80 -v