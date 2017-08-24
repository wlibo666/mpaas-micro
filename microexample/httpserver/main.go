package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/wlibo666/mpaas-micro/micro/log"
	"github.com/wlibo666/mpaas-micro/micro/service"
	calculator "github.com/wlibo666/mpaas-micro/microexample/proto"
	"golang.org/x/net/context"
)

var (
	logFile     string
	listenAddr  string
	serviceName string
)

var (
	client calculator.CalculatorClient
)

type response struct {
	Num1 int64  `json:"num1"`
	Num2 int64  `json:"num2"`
	Res  int64  `json:"res"`
	Err  string `json:"err"`
}

// 初始化微服务客户端
func initServiceClient(name string) error {
	svc, err := service.NewClientService("microserver1", "")
	if err != nil {
		return err
	}
	svc.Init()
	//根据服务名称创建微服务客户端
	client = calculator.NewCalculatorClient("microserver1", svc.Client())
	return nil
}

// 从URL里获取参数
func getParameters(r *http.Request) (int64, int64, error) {
	params := r.URL.Query()
	tmpstr := params.Get("num1")
	if len(tmpstr) == 0 {
		return 0, 0, fmt.Errorf("lost num1")
	}
	num1, err := strconv.Atoi(tmpstr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid num1:%s", tmpstr)
	}
	tmpstr = params.Get("num2")
	if len(tmpstr) == 0 {
		return 0, 0, fmt.Errorf("lost num2")
	}
	num2, err := strconv.Atoi(tmpstr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid num2:%s", tmpstr)
	}
	return int64(num1), int64(num2), nil
}

func writeResData(num1, num2, res int64, err string, w http.ResponseWriter) error {
	d := &response{
		Num1: num1,
		Num2: num2,
		Res:  res,
		Err:  err,
	}
	data, e := json.Marshal(d)
	if e != nil {
		log.Error("MarshalErr", e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Marshal failed,err:%s", e.Error())))
		return e
	}
	w.Write(data)
	return nil
}

// 加法处理
func addHandler(w http.ResponseWriter, r *http.Request) {
	num1, num2, err := getParameters(r)
	if err != nil {
		writeResData(0, 0, 0, err.Error(), w)
		return
	}
	// 调用加法接口
	addRes, err := client.Add(context.Background(), &calculator.AddReq{
		Num1: num1,
		Num2: num2,
	})
	if err != nil {
		writeResData(0, 0, 0, err.Error(), w)
		return
	}
	log.Debug("cmd", "add", "num1", num1, "num2", num2, "res", addRes.Res)
	writeResData(num1, num2, addRes.Res, "", w)
}

// 乘法处理
func mulHandler(w http.ResponseWriter, r *http.Request) {
	num1, num2, err := getParameters(r)
	if err != nil {
		writeResData(0, 0, 0, err.Error(), w)
		return
	}
	// 调用乘法接口
	mulRes, err := client.Mul(context.Background(), &calculator.MulReq{
		Num1: num1,
		Num2: num2,
	})
	if err != nil {
		writeResData(0, 0, 0, err.Error(), w)
		return
	}
	log.Debug("cmd", "mul", "num1", num1, "num2", num2, "res", mulRes.Res)
	writeResData(num1, num2, mulRes.Res, "", w)
}

func GetEnvVar(key, defV string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		return defV
	}
	return v
}

func main() {
	logFile = GetEnvVar("HTTPSERVER_LOG_FILE", "./httpserver.log")
	listenAddr = GetEnvVar("HTTPSERVER_LISTEN_ADDR", ":8080")
	serviceName = GetEnvVar("SERVICE1_NAME", "microserver1")

	// 创建日志文件
	err := log.NewDefaultLogger(logFile, log.DEBUG)
	if err != nil {
		fmt.Printf("NewDefaultLogger:%s failed,err:%s", logFile, err.Error())
		return
	}
	log.Info("program", os.Args[0], "status", "start")
	log.Info("logfile", logFile, "listenAddr", listenAddr, "serviceName", serviceName)

	// 初始化微服务客户端
	err = initServiceClient(serviceName)
	if err != nil {
		log.Error("initServiceClient", err.Error())
		return
	}

	// 设置http接口处理函数
	http.HandleFunc("/v1/add", addHandler)
	http.HandleFunc("/v1/mul", mulHandler)

	// 启动http服务
	err = http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Error("ListenAndServeErr", err.Error())
	}
}
