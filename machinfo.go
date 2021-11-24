// 获取机器信息
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// 获取系统信息
func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	strDisk := GetDiskInfoma()

	strDisk = "磁盘使用:\n" + strDisk + "\n"
	strMem := fmt.Sprintf("内存使用: %f\n", GetMemoryInfo())
	strCPU := fmt.Sprintf("CPU使用: %f\n", GetCPUInfoma())

	fmt.Fprintln(w, strDisk+strMem+strCPU)
}

// cpu信息
func GetCPUInfo() {
	fmt.Fprintln(w, GetCpuPercent())
}

// CPU占用比例
func GetCPUPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetCPUDetail() {

}

// 内存信息
func GetMemInfo() {
	fmt.Fprintln(w, GetMemPercent())
}

// 内存占用比例
func GetMemPercent() {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

// 内存详细信息
func GetMemDetail() {

}

// 磁盘信息 已经占用比例，空闲空间
func GetDiskInfo() string {
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Print("get Partitions failed, err:%v\n", err)
		return "disk not find"
	}

	strOut := ""
	for _, part := range parts {
		diskInfo, _ := disk.Usage(part.Mountpoint)
		strtmp := fmt.Sprintf("%v disk info:已经使用占比:%.2f%% 空闲空间:%.2fG\n", diskInfo.Path, diskInfo.UsedPercent, (float64)(diskInfo.Free)/1024/1024/1024)
		strOut += strtmp
	}

	return strOut
}

// 磁盘占用比例
func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

func GetDiskDetail() {

}
