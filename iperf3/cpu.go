package iperf3

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func init() {
	go func() {
		for {
			time.Sleep(time.Second * 2)

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

			info := fmt.Sprintf("CPU: %.2f%% MEM: %.2f%%", cpuPercent[0], memInfo.UsedPercent)

			ServerStatusUpdate(info)
			ClientStatusUpdate(info)

			logs.Info("system resource info %s", info)
		}
	}()

}
