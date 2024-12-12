package main

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func init() {
	go func() {
		for {
			time.Sleep(time.Second)

			cpuPercent, err := cpu.Percent(0, false)
			if err != nil {
				logs.Warning("Get CPU Usage failed, %s", err.Error())
				continue
			}

			memInfo, err := mem.VirtualMemory()
			if err != nil {
				logs.Warning("Get Memory Usage failed, %s", err.Error())
				continue
			}

			logs.Info("CPU 使用率: %.2f%%\n", cpuPercent[0])
			logs.Info("内存使用率: %.2f%%\n", memInfo.UsedPercent)
		}
	}()

}
