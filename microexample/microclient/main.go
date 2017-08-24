package main

import (
	"fmt"
	"golang.org/x/net/context"

	"github.com/wlibo666/mpaas-micro/micro/service"
	calculator "github.com/wlibo666/mpaas-micro/microexample/proto"
)

func main() {
	// 根据服务名称和版本号创建微服务结构(版本可以为空)
	svc, err := service.NewClientService("microserver1", "")
	if err != nil {
		fmt.Printf("NewService failed,err:%s\n", err.Error())
		return
	}
	svc.Init()
	//根据服务名称创建微服务客户端
	cli := calculator.NewCalculatorClient("microserver1", svc.Client())
	// 调用加法接口
	addRes, err := cli.Add(context.Background(), &calculator.AddReq{
		Num1: 66,
		Num2: 99,
	})
	if err != nil {
		fmt.Printf("add failed,err:%s\n", err.Error())
		return
	}
	fmt.Printf("66 + 99 = %d\n", addRes.Res)
	// 调用乘法接口
	mulRes, err := cli.Mul(context.Background(), &calculator.MulReq{
		Num1: 66,
		Num2: 99,
	})
	if err != nil {
		fmt.Printf("mul failed,err:%s\n", err.Error())
		return
	}
	fmt.Printf("66 * 99 = %d\n", mulRes.Res)
}
