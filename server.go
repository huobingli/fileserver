package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

const BaseUploadPath = "D:\\client_pack\\work_dir\\pdb_file"

func TimeParseYYYYMMDD(in string, sub string) (out time.Time, err error) {
	layout := "2006" + sub + "01" + sub + "02"
	out, err = time.ParseInLocation(layout, in, time.Local)
	if err != nil {
		return
	}
	return
}

func getCurDay() (date int) {
	curTime := time.Now()
	year := curTime.Year()
	month := int(curTime.Month())
	day := curTime.Day()
	return year*10000 + month*100 + day
}

func downfile(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("context")
	filepath := BaseUploadPath + filename

	w.Header().Set("Pragma", "No-cache")
	w.Header().Set("Cache-Control", "No-cache")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/text/plain")
	http.ServeFile(w, r, filepath)
}

// 清理目录中一个月前的的临时文件  文件格式 日期_创建时间
func cleanfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "start clean file ... please wait !!")

	curDay := getCurDay()

	// 读取上传目录下文件名
	_dir, err := ioutil.ReadDir(BaseUploadPath)
	if err != nil {
		return
	}

	for _, _file := range _dir {
		if _file.IsDir() {
			dirname := _file.Name()

			// 分割文件名，找到文件夹名中时间
			comm := strings.Index(dirname, "_")
			strDirDate := dirname[:comm]

			dirDate, err := strconv.Atoi(strDirDate)
			if err != nil {
				fmt.Print("strconv.Atoi, err:%v\n", err)
			}

			// 遍历删除一个月前的文件夹
			if dirDate < curDay-100 {
				fmt.Println(dirDate)
				removeDir := BaseUploadPath + dirname
				os.RemoveAll(removeDir)
			}
		}
	}

}

// 获取CPU使用情况
func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetCpuInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, GetCpuPercent())
}

func GetCPUInfoma() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	// fmt.Print(percent)
	// info, _ := cpu.Info()
	info, _ := cpu.Times(false)
	fmt.Print(info)
	return percent[0]
}

// func GetCPUDetail() string {
// 	info, _ := cpu.Info()
// 	fmt.Print(info)
// }

// 获取内存使用情况
func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

func GetMemoryInfo() float64 {
	memInfo, _ := mem.VirtualMemory()
	fmt.Print(memInfo)
	return memInfo.UsedPercent
}

func GetMemDetail() string {
	memInfo, _ := mem.VirtualMemory()
	strRet := fmt.Sprintf("使用比例: %f%%, 总共：%.2fG, 已经使用：%.2fG, 未使用：%.2fG\n", memInfo.UsedPercent, (float64)(memInfo.Total)/1024/1024/1024, (float64)(memInfo.Used)/1024/1024/1024, (float64)(memInfo.Available)/1024/1024/1024)
	return strRet
}

func GetMemInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, GetMemPercent())
}

// 获取磁盘使用情况
func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

// 后面想获取磁盘信息，预留一下口子与参考blog
// https://blog.csdn.net/whatday/article/details/109620192
// TODO需要调试一下输出
func getdiskinfo(w http.ResponseWriter, r *http.Request) {
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Print("get Partitions failed, err:%v\n", err)
		return
	}

	for _, part := range parts {
		fmt.Print("part:%v\n", part.String())
		diskInfo, _ := disk.Usage(part.Mountpoint)
		fmt.Print("disk info:used:%f free:%f\n", diskInfo.UsedPercent, diskInfo.Free)
	}

	ioStat, _ := disk.IOCounters()
	strOut := ""
	for k, v := range ioStat {
		strtmp := fmt.Sprintf("%v:%v\n", k, v)
		strOut += strtmp
	}

	fmt.Fprintln(w, strOut)
}

func GetDiskInfoma() string {
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Print("get Partitions failed, err:%v\n", err)
		return "disk not find"
	}

	strOut := ""
	for _, part := range parts {
		// fmt.Print("part:%v\n", part.String())
		diskInfo, _ := disk.Usage(part.Mountpoint)
		// fmt.Print(diskInfo)
		// fmt.Print("disk info:used:%f free:%f\n", diskInfo.UsedPercent, diskInfo.Free)
		strtmp := fmt.Sprintf("%v disk info:已经使用占比:%.2f%% 空闲空间:%.2fG\n", diskInfo.Path, diskInfo.UsedPercent, (float64)(diskInfo.Free)/1024/1024/1024)
		strOut += strtmp
	}

	// 磁盘IO
	// ioStat, _ := disk.IOCounters()
	// strOut := ""
	// for k, v := range ioStat {
	// 	strtmp := fmt.Sprintf("%v:%v\n", k, v)
	// 	strOut += strtmp
	// }

	return strOut
	// fmt.Fprintln(w, strOut)
}

// 获取系统信息 磁盘，内存，cpu使用信息
func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	strDisk := GetDiskInfoma()

	strDisk = "磁盘使用:\n" + strDisk + "\n"
	// strMem := fmt.Sprintf("内存使用: %f\n", GetMemoryInfo())
	strMem := "内存使用:\n" + GetMemDetail() + "\n"
	strCPU := fmt.Sprintf("CPU使用: %f\n", GetCPUInfoma())

	fmt.Fprintln(w, strDisk+strMem+strCPU)
}

func main() {
	mux := http.NewServeMux()

	// 其他接口
	mux.HandleFunc("/cleanfile", cleanfile)
	// mux.HandleFunc("/getdiskinfo", getdiskinfo)
	mux.HandleFunc("/downfile", downfile)
	mux.HandleFunc("/getmeminfo", GetMemInfo)
	mux.HandleFunc("/getcpuinfo", GetCpuInfo)
	mux.HandleFunc("/GetSystemInfo", GetSystemInfo)

	// 文件服务器
	mux.Handle("/", http.FileServer(http.Dir(BaseUploadPath)))

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
