package main

import (
	"fmt"
	"io"
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
	"github.com/shirou/gopsutil/winservices"
)

const BaseUploadPath = "D:\\client_pack\\work_dir\\pdb_file"

// func TimeParseYYYYMMDD(in string, sub string) (out time.Time, err error) {
// 	layout := "2006" + sub + "01" + sub + "02"
// 	out, err = time.ParseInLocation(layout, in, time.Local)
// 	if err != nil {
// 		return
// 	}
// 	return
// }

func getCurDay() (date int) {
	curTime := time.Now()
	year := curTime.Year()
	month := int(curTime.Month())
	day := curTime.Day()
	return year*10000 + month*100 + day
}

func Downfile(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("context")
	filepath := BaseUploadPath + filename

	w.Header().Set("Pragma", "No-cache")
	w.Header().Set("Cache-Control", "No-cache")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/text/plain")
	http.ServeFile(w, r, filepath)
}

const UploadPath = "D:\\gitProject\\fileserver\\upload\\test"

func Uploadfile(w http.ResponseWriter, request *http.Request) {
	fmt.Println("handle upload")
	//文件上传只允许POST方法
	if request.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}
	fmt.Println("handle upload1")
	//从表单中读取文件
	file, fileHeader, err := request.FormFile("file")
	fmt.Println(file)
	fmt.Println(fileHeader)
	if err != nil {
		_, _ = io.WriteString(w, "Read file error")
		return
	}
	fmt.Println("handle upload2")
	//defer 结束时关闭文件
	defer file.Close()
	fmt.Println("filename: " + fileHeader.Filename)

	//创建文件
	newFile, err := os.Create(UploadPath + "/" + fileHeader.Filename)
	if err != nil {
		_, _ = io.WriteString(w, "Create file error")
		return
	}
	fmt.Println("handle upload3")
	//defer 结束时关闭文件
	defer newFile.Close()

	//将文件写到本地
	_, err = io.Copy(newFile, file)
	if err != nil {
		_, _ = io.WriteString(w, "Write file error")
		return
	}
	fmt.Println("handle upload4")
	_, _ = io.WriteString(w, "Upload success")
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
func GetCPUInfo() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	info, _ := cpu.Times(false)
	fmt.Print(info)
	return percent[0]
}

// 获取内存使用情况
func GetMemInfo() string {
	memInfo, _ := mem.VirtualMemory()
	strRet := fmt.Sprintf("使用比例: %f%%, 总共：%.2fG, 已经使用：%.2fG, 未使用：%.2fG\n", memInfo.UsedPercent, (float64)(memInfo.Total)/1024/1024/1024, (float64)(memInfo.Used)/1024/1024/1024, (float64)(memInfo.Available)/1024/1024/1024)
	return strRet
}

// 获取磁盘使用情况
// 后面想获取磁盘信息，预留一下口子与参考blog
// https://blog.csdn.net/whatday/article/details/109620192
// TODO需要调试一下输出
func GetDiskInfo2(w http.ResponseWriter, r *http.Request) {
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

// 获取系统信息 磁盘，内存，cpu使用信息
func GetSystemInfo(w http.ResponseWriter, r *http.Request) {
	strDisk := GetDiskInfo()

	strDisk = "磁盘使用:\n" + strDisk + "\n"
	strMem := "内存使用:\n" + GetMemInfo() + "\n"
	strCPU := fmt.Sprintf("CPU使用: %f\n", GetCPUInfo())

	fmt.Fprintln(w, strDisk+strMem+strCPU)
}

func GetProcessInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "console")
}

func GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	services, _ := winservices.ListServices()

	fmt.Print(services)
	for _, service := range services {
		newservice, _ := winservices.NewService(service.Name)
		newservice.GetServiceDetail()
		fmt.Println("Name:", newservice.Name, "Binary Path:", newservice.Config.BinaryPathName, "State: ", newservice.Status.State)
	}

	fmt.Fprintln(w, "console")
}

func main() {
	mux := http.NewServeMux()

	// 其他接口 清理文件，上传下载文件
	mux.HandleFunc("/cleanfile", cleanfile)
	mux.HandleFunc("/Uploadfile", Uploadfile)
	mux.HandleFunc("/Downfile", Downfile)

	// 获取系统信息
	mux.HandleFunc("/GetSystemInfo", GetSystemInfo)
	mux.HandleFunc("/GetServiceInfo", GetServiceInfo)

	// 其他操作
	// svn更新
	mux.HandleFunc("UpdateSvn", UpdateSvn)
	// kill进程
	mux.HandleFunc("KillPython", KillPython)

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
