// 获取机器信息
package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

// 获取系统信息
func GetSystemInfo() string {
	strDisk := GetDiskDetail()

	strDisk = "磁盘使用:\n" + strDisk + "\n"
	// strMem := fmt.Sprintf("内存使用: %f\n", GetMemoryInfo())
	strMem := "内存使用:\n" + GetMemDetail() + "\n"
	strCPU := fmt.Sprintf("CPU使用: %f\n", GetCPUInfoma())

	fmt.Println(strDisk + strMem + strCPU)
	return strDisk + strMem + strCPU
}

// cpu信息
func GetCPUInfo() {
	fmt.Fprintln(w, GetCPUPercent())
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
	fmt.Print(GetMemPercent())
}

// 内存占用比例
func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

// 内存详细信息
func GetMemDetail() string {
	memInfo, _ := mem.VirtualMemory()
	strRet := fmt.Sprintf("使用比例: %f%%, 总共：%.2fG, 已经使用：%.2fG, 未使用：%.2fG\n", memInfo.UsedPercent, (float64)(memInfo.Total)/1024/1024/1024, (float64)(memInfo.Used)/1024/1024/1024, (float64)(memInfo.Available)/1024/1024/1024)
	return strRet
}

// 磁盘信息 已经占用比例，空闲空间
func GetDiskInfo() {
	fmt.Print(GetDiskDetail())
}

// 磁盘占用比例
func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

func GetDiskDetail() string {
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
