// Package service provide micro service function to create micro service with
// grpc and http.
package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
)

const (
	mpaasDeployFlag = "MESOS_CONTAINER_NAME"
)

var (
	default_consul     string = "10.135.28.154:8500,10.185.30.76:8500,10.112.34.55:8500"
	consul_from_remote string = ""
)

func updateConsul(from string) {
	if len(from) == 0 {
		return
	}
	resp, err := http.Get(from)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || resp.ContentLength < 10 {
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	consul_from_remote = string(data)
}

func isMpaasDeploy() bool {
	if strings.Contains(os.Getenv(mpaasDeployFlag), "mesos") {
		return true
	}
	return false
}

func getMpaasDeoloyAddr() string {
	return os.Getenv("HOST") + ":" + os.Getenv("PORT")
}

func setenv(key, value string) {
	if len(strings.Trim(value, " ")) == 0 {
		unsetenv(key)
		return
	}
	os.Setenv(key, value)
}

func unsetenv(key string) {
	os.Unsetenv(key)
}

func genVersion() string {
	now := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()/10)
}

func printConf(secs map[string]*ini.Section) {
	sec, ok := secs["custom"]
	if ok {
		fmt.Printf("%s:%s\n", "micro_type", getKeyString(sec, "micro_type", ""))
		fmt.Printf("%s:%s\n", "server_name", getKeyString(sec, "server_name", ""))
		fmt.Printf("%s:%s\n", "server_version", getKeyString(sec, "server_version", ""))
		fmt.Printf("%s:%s\n", "server_address", getKeyString(sec, "server_address", ""))
	}
	sec, ok = secs["commom"]
	if ok {
		fmt.Printf("%s:%s\n", "server_id", getKeyString(sec, "server_id", ""))
		fmt.Printf("%s:%s\n", "server_advertise", getKeyString(sec, "server_advertise", ""))
		fmt.Printf("%s:%s\n", "client_request_timeout", getKeyString(sec, "client_request_timeout", "3s"))
		fmt.Printf("%s:%s\n", "client_pool_size", getKeyString(sec, "client_pool_size", "2"))
		fmt.Printf("%s:%s\n", "client_pool_ttl", getKeyString(sec, "client_pool_ttl", "5m"))
		fmt.Printf("%s:%s\n", "client_retries", getKeyString(sec, "client_retries", "2"))
	}
	sec, ok = secs["registry"]
	if ok {
		fmt.Printf("%s:%s\n", "registry", getKeyString(sec, "registry", "consul"))
		fmt.Printf("%s:%s\n", "registry_addr_from", getKeyString(sec, "registry_addr_from", ""))
		fmt.Printf("%s:%s\n", "registry_address", getKeyString(sec, "registry_address", default_consul))
		fmt.Printf("%s:%s\n", "registry_ttl", getKeyString(sec, "registry_ttl", ""))
		fmt.Printf("%s:%s\n", "registry_interval", getKeyString(sec, "registry_interval", ""))
	}
}

func initEnviroment(secs map[string]*ini.Section) {
	sec, ok := secs["custom"]
	if ok {
		setenv("MICRO_SERVER_NAME", getKeyString(sec, "server_name", "mpaas_micro_service"))
		setenv("MICRO_SERVER_VERSION", getKeyString(sec, "server_version", genVersion()))
		setenv("MICRO_SERVER_ADDRESS", getKeyString(sec, "server_address", ""))
	}
	sec, ok = secs["commom"]
	if ok {
		setenv("MICRO_SERVER_ID", getKeyString(sec, "server_id", ""))
		if isMpaasDeploy() {
			setenv("MICRO_SERVER_ADVERTISE", getMpaasDeoloyAddr())
		} else {
			setenv("MICRO_SERVER_ADVERTISE", getKeyString(sec, "server_advertise", ""))
		}
		setenv("MICRO_CLIENT_REQUEST_TIMEOUT", getKeyString(sec, "client_request_timeout", "3s"))
		setenv("MICRO_CLIENT_POOL_SIZE", getKeyString(sec, "client_pool_size", "2"))
		setenv("MICRO_CLIENT_POOL_TTL", getKeyString(sec, "client_pool_ttl", "5m"))
		setenv("MICRO_CLIENT_RETRIES", getKeyString(sec, "client_retries", "2"))
	}
	sec, ok = secs["registry"]
	if ok {
		setenv("MICRO_REGISTRY", getKeyString(sec, "registry", "consul"))
		updateConsul(getKeyString(sec, "registry_addr_from", ""))
		if len(consul_from_remote) == 0 {
			setenv("MICRO_REGISTRY_ADDRESS", getKeyString(sec, "registry_address", default_consul))
		} else {
			default_consul = consul_from_remote
			setenv("MICRO_REGISTRY_ADDRESS", default_consul)
		}
	}
}

func newDefaultConfig(serviceName, serviceVersion string, conf map[string]string) (map[string]*ini.Section, error) {
	tmpConf := conf
	if tmpConf == nil {
		tmpConf = make(map[string]string)
	}
	tmpConf["server_name"] = serviceName
	tmpConf["server_version"] = serviceVersion

	file := ini.Empty()
	err := file.NewSections("custom", "commom", "registry")
	if err != nil {
		return nil, err
	}
	newConf := make(map[string]*ini.Section)
	newConf["custom"] = file.Section("custom")
	newConf["commom"] = file.Section("commom")
	newConf["registry"] = file.Section("registry")

	// custom section
	v, ok := tmpConf["micro_type"]
	if ok && len(v) > 0 {
		file.Section("custom").NewKey("micro_type", v)
	} else {
		file.Section("custom").NewKey("micro_type", "grpc")
	}
	v, ok = tmpConf["server_name"]
	if ok && len(v) > 0 {
		file.Section("custom").NewKey("server_name", v)
	} else {
		file.Section("custom").NewKey("server_name", "mpaas_micro_service")
	}
	v, ok = tmpConf["server_version"]
	if ok && len(v) > 0 {
		file.Section("custom").NewKey("server_version", v)
	} else {
		file.Section("custom").NewKey("server_version", "")
	}
	v, ok = tmpConf["server_address"]
	if ok && len(v) > 0 {
		file.Section("custom").NewKey("server_address", v)
	} else {
		file.Section("custom").NewKey("server_address", "")
	}
	// common section
	v, ok = tmpConf["server_id"]
	if ok && len(v) > 0 {
		file.Section("commom").NewKey("server_id", v)
	} else {
		file.Section("commom").NewKey("server_id", "")
	}
	v, ok = tmpConf["server_advertise"]
	if ok && len(v) > 0 {
		file.Section("commom").NewKey("server_advertise", v)
	} else {
		file.Section("commom").NewKey("server_advertise", "")
	}
	v, ok = tmpConf["client_request_timeout"]
	if ok && len(v) > 0 {
		file.Section("commom").NewKey("client_request_timeout", v)
	} else {
		file.Section("commom").NewKey("client_request_timeout", "3s")
	}
	v, ok = tmpConf["client_pool_size"]
	if ok && len(v) > 0 {
		file.Section("commom").NewKey("client_pool_size", v)
	} else {
		file.Section("commom").NewKey("client_pool_size", "2")
	}
	v, ok = tmpConf["client_pool_ttl"]
	if ok && len(v) > 0 {
		file.Section("commom").NewKey("client_pool_ttl", v)
	} else {
		file.Section("commom").NewKey("client_pool_ttl", "5m")
	}
	v, ok = tmpConf["client_retries"]
	if ok && len(v) > 0 {
		file.Section("commom").NewKey("client_retries", v)
	} else {
		file.Section("commom").NewKey("client_retries", "2")
	}
	// registry section
	v, ok = tmpConf["registry"]
	if ok && len(v) > 0 {
		file.Section("registry").NewKey("registry", v)
	} else {
		file.Section("registry").NewKey("registry", "consul")
	}
	v, ok = tmpConf["registry_address"]
	if ok && len(v) > 0 {
		file.Section("registry").NewKey("registry_address", v)
	} else {
		file.Section("registry").NewKey("registry_address", default_consul)
	}
	v, ok = tmpConf["registry_ttl"]
	if ok && len(v) > 0 {
		file.Section("registry").NewKey("registry_ttl", v)
	} else {
		file.Section("registry").NewKey("registry_ttl", "0")
	}
	v, ok = tmpConf["registry_interval"]
	if ok && len(v) > 0 {
		file.Section("registry").NewKey("registry_interval", v)
	} else {
		file.Section("registry").NewKey("registry_interval", "300")
	}
	v, ok = tmpConf["check_interval"]
	if ok && len(v) > 0 {
		file.Section("registry").NewKey("check_interval", v)
	} else {
		file.Section("registry").NewKey("check_interval", "15")
	}

	return newConf, nil
}

func parseConf(confFile string) (map[string]*ini.Section, error) {
	conf, err := ini.InsensitiveLoad(confFile)
	if err != nil {
		return nil, err
	}
	secs := make(map[string]*ini.Section)
	for _, sec := range conf.Sections() {
		secs[sec.Name()] = sec
	}
	return secs, nil
}

func getKeyString(sec *ini.Section, keyName, defaultV string) string {
	if sec == nil {
		return ""
	}
	return sec.Key(keyName).MustString(defaultV)
}

func getKeyInt(sec *ini.Section, keyName string, defaultV int) int {
	if sec == nil {
		return 0
	}
	return sec.Key(keyName).MustInt(defaultV)
}

// NewServerService generate a micro service by config file,eg:
// github.com/wlibo666/mpaas-micro/micro/service/default.ini.
// parmeter configFile: config file name(format:ini)
func NewServerService(configFile string) (micro.Service, error) {
	secs, err := parseConf(configFile)
	if err != nil {
		return nil, err
	}
	//printConf(secs)
	initEnviroment(secs)

	var opts []micro.Option
	sec, ok := secs["registry"]
	if ok {
		if len(consul_from_remote) == 0 {
			r := consul.NewRegistry(registry.Addrs(strings.Split(getKeyString(sec, "registry_address", default_consul), ",")...))
			opts = append(opts, micro.Registry(r))
		} else {
			r := consul.NewRegistry(registry.Addrs(strings.Split(default_consul, ",")...))
			opts = append(opts, micro.Registry(r))
		}

		opts = append(opts, micro.RegisterInterval(time.Duration(getKeyInt(sec, "registry_interval", 60))*time.Second))
		opts = append(opts, micro.CheckInterval(time.Duration(getKeyInt(sec, "check_interval", 30))*time.Second))
	}
	srvType := "grpc"
	sec, ok = secs["custom"]
	if ok {
		srvType = getKeyString(sec, "micro_type", "grpc")
		if !isMpaasDeploy() {
			opts = append(opts, micro.CheckTCP(getKeyString(sec, "server_address", "")))
		} else {
			opts = append(opts, micro.CheckTCP(getMpaasDeoloyAddr()))
		}
	}

	var svr micro.Service
	switch srvType {
	case "http":
		svr = micro.NewService(opts...)
	case "grpc":
		svr = grpc.NewService(opts...)
	default:
		svr = grpc.NewService(opts...)
	}
	return svr, nil
}

// NewClientService generate a micro service by service name,service version.
// serviceName is service's name,serviceVersion is service's version,
// optConfs is optional config item,the same with server config file's item
// (github.com/wlibo666/mpaas-micro/micro/service/default.ini)
func NewClientService(serviceName, serviceVersion string, optConfs ...map[string]string) (micro.Service, error) {
	var secs map[string]*ini.Section
	var err error
	if len(optConfs) > 0 {
		secs, err = newDefaultConfig(serviceName, serviceVersion, optConfs[0])
	} else {
		secs, err = newDefaultConfig(serviceName, serviceVersion, nil)
	}
	if err != nil {
		return nil, err
	}
	//printConf(secs)
	initEnviroment(secs)

	var opts []micro.Option
	sec, ok := secs["registry"]
	if ok {
		r := consul.NewRegistry(registry.Addrs(strings.Split(getKeyString(sec, "registry_address", default_consul), ",")...))
		opts = append(opts, micro.Registry(r))
		opts = append(opts, micro.RegisterInterval(time.Duration(getKeyInt(sec, "registry_interval", 60))*time.Second))
		opts = append(opts, micro.CheckInterval(time.Duration(getKeyInt(sec, "check_interval", 30))*time.Second))
	}
	srvType := "grpc"
	sec, ok = secs["custom"]
	if ok {
		srvType = getKeyString(sec, "micro_type", "grpc")
		opts = append(opts, micro.CheckTCP(getKeyString(sec, "server_address", "")))
	}

	var svr micro.Service
	switch srvType {
	case "http":
		svr = micro.NewService(opts...)
	case "grpc":
		svr = grpc.NewService(opts...)
	default:
		svr = grpc.NewService(opts...)
	}
	return svr, nil
}
