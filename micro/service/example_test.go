package service_test

import (
	"fmt"
	"github.com/wlibo666/mpaas-micro/micro/service"
)

func ExampleService_client() {
	svc, err := service.NewClientService("mymicroservice", "1.0.0")
	if err != nil {
		fmt.Printf("NewClientService failed,err:%s\n", err.Error())
		return
	}
	svc.Init()
}
