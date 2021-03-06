package log

import (
	"fmt"
	"testing"
)

func TestDebug(t *testing.T) {
	err := NewDefaultLogger("default.log", DEBUG)
	if err != nil {
		fmt.Printf("new a.log failed,err:%s\n", err.Error())
		return
	}
	Debug("debug", "true")
	Info("info", "true")
	Warn("aaa", "true")
	Error("error", "true")

	err = NewLogger("newlog", "newlog.log")
	if err != nil {
		fmt.Printf("new b.log failed,err:%s\n", err.Error())
		return
	}
	DebugWithLogger("newlog", "debug", "yes")
	InfoWithLogger("newlog", "info", "yes")
	WarnWithLogger("newlog", "yes", "no")
	ErrorWithLogger("newlog", "err", "yes")
}
