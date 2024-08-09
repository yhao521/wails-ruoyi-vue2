package monitor

import (
	"mySparkler/backend/model/monitor"
)

var (
	SysOperLog = serviceSysOperLog{}
)

type serviceSysOperLog struct{}

func (s serviceSysOperLog) AddSysOperLog(data string) {
	var operLog = monitor.SysOperLog{
		Title:      "错误日志",
		JsonResult: string(data),
		Status:     "1",
	}
	operLog.OperationLogAdd()
}
