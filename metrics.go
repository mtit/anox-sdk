package sdk

import (
	"log"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func getSystemMetrics() (cpuPercent float64, memTotal, memAvail int64) {
	percents, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("[Anox SDK] 读取CPU指标失败: %v", err)
	} else if len(percents) > 0 {
		cpuPercent = percents[0]
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("[Anox SDK] 读取RAM指标失败: %v", err)
	} else {
		memTotal = int64(vm.Total / 1024 / 1024)
		memAvail = int64(vm.Available / 1024 / 1024)
	}

	return
}
