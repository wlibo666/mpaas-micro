package main

import (
	"fmt"
	"os"

	"github.com/wlibo666/mpaas-micro/micro/log"
	"github.com/wlibo666/mpaas-micro/micro/service"
	calculator "github.com/wlibo666/mpaas-micro/microexample/proto"
	"golang.org/x/net/context"
)

var (
	configFile string
	logFile    string
)

/*
	加法/乘法 服务接口实现
*/
type calc struct{}

func (c *calc) Add(ctx context.Context, req *calculator.AddReq, res *calculator.AddRes) error {
	res.Res = req.Num1 + req.Num2
	log.Debug("addNum1", req.Num1, "addNum2", req.Num2, "addRes", res.Res)
	return nil
}

func (c *calc) Mul(ctx context.Context, req *calculator.MulReq, res *calculator.MulRes) error {
	res.Res = req.Num1 * req.Num2
	log.Debug("nulNum1", req.Num1, "mulNum2", req.Num2, "mulRes", res.Res)
	return nil
}

func GetEnvVar(key, defV string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return defV
	}
	return v
}

func main() {
	//  从环境变量获取必需参数 配置文件路径/日志文件路径
	configFile = GetEnvVar("SERVICE_CONFIG_FILE", "./config.ini")
	logFile = GetEnvVar("SERVER_LOG_FILE", "./server.log")

	// 打开日志文件
	err := log.NewDefaultLogger(logFile, log.DEBUG)
	if err != nil {
		fmt.Printf("log Init failed,err:%s\n", err.Error())
		return
	}

	log.Debug("logfile", logFile, "config", configFile)

	// 根据配置文件创建微服务
	svc, err := service.NewServerService(configFile)
	if err != nil {
		log.Error("NewServerService err", err.Error())
		return
	}
	// 服务初始化
	svc.Init()
	// 注册服务实体
	calculator.RegisterCalculatorHandler(svc.Server(), new(calc))
	// 启动服务
	err = svc.Run()
	if err != nil {
		log.Error("Run err", err.Error())
		return
	}
}
