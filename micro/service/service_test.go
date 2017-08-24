package service

import (
	"github.com/go-ini/ini"
	"testing"
)

func TestIniConf(t *testing.T) {
	conf, err := ini.InsensitiveLoad("./default.ini")
	if err != nil {
		t.Fatalf("new config failed,err:%s\n", err.Error())
	}
	sec, err := conf.GetSection("custom")
	if err != nil {
		t.Fatalf("get section failed,err:%s", err.Error())
	}
	key, err := sec.GetKey("micro_type")
	if err != nil {
		t.Fatalf("get key failed,err:%s", err.Error())
	}
	t.Logf("micro_type:%s", key.String())
}

func TestNewService(t *testing.T) {
	svc, err := NewServerService("./default.ini")
	if err != nil {
		t.Fatalf("NewService failed,err:%s", err.Error())
	}
	return
	t.Log("will run svc...")
	err = svc.Run()
	if err != nil {
		t.Fatalf("run failed,err:%s", err.Error())
	}
}
