package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/winservices"
)

type conf struct {
	Pdb_dir     string
	File_dir    string
	File_server string
}

var cf conf

// const cf conf
func load_default_config() error {
	// var cf conf
	var path string = "./conf.toml"
	if _, err := toml.DecodeFile(path, &cf); err != nil {
		return err
	}

	return nil
}

func load_config(path string) error {
	// var cf conf
	// var path string = "./conf.toml"
	if _, err := toml.DecodeFile(path, &cf); err != nil {
		return err
	}

	return nil
}

func Testdir(w http.ResponseWriter, request *http.Request) {
	str := "pdb目录" + cf.Pdb_dir + "\n" + "上传目录" + cf.File_dir + "\n" + "文件服务根目录" + cf.File_server
	fmt.Fprintln(w, str)
}

func Reloadconf(w http.ResponseWriter, request *http.Request) {
	if err := load_default_config(); err != nil {
		fmt.Fprintln(w, "reload error! %s", err)
	} else {
		fmt.Fprintln(w, "reload success")
	}
}

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
	filepath := cf.File_dir + filename

	w.Header().Set("Pragma", "No-cache")
	w.Header().Set("Cache-Control", "No-cache")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/text/plain")
	http.ServeFile(w, r, filepath)
}

func Uploadpdb(w http.ResponseWriter, request *http.Request) {
	//文件上传只允许POST方法
	if request.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Uploadfile error Method not allowed,Please use [Post] Method")
		return
	}

	//从表单中读取文件
	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, "Uploadfile error file = %s Read file error", file)
		return
	}

	//defer 结束时关闭文件
	defer file.Close()

	//创建文件
	filePath := cf.Pdb_dir + "/" + fileHeader.Filename
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintln(w, "Uploadfile error Create file error. path is ", filePath)
		return
	}

	//defer 结束时关闭文件
	defer newFile.Close()

	//将文件写到本地
	_, err = io.Copy(newFile, file)
	if err != nil {
		_, _ = io.WriteString(w, "Write file error")
		return
	}
	fmt.Fprintln(w, "Uploadfile success filePath = ", filePath)
}

func Uploadfile(w http.ResponseWriter, request *http.Request) {
	//文件上传只允许POST方法
	if request.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Uploadfile error Method not allowed,Please use [Post] Method")
		return
	}

	//从表单中读取文件
	file, fileHeader, err := request.FormFile("file")
	if err != nil {
		fmt.Fprintln(w, "Uploadfile error file = %s Read file error", file)
		return
	}

	//defer 结束时关闭文件
	defer file.Close()

	//创建文件
	filePath := cf.File_dir + "/" + fileHeader.Filename
	newFile, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintln(w, "Uploadfile error Create file error. path is ", filePath)
		return
	}

	//defer 结束时关闭文件
	defer newFile.Close()

	//将文件写到本地
	_, err = io.Copy(newFile, file)
	if err != nil {
		_, _ = io.WriteString(w, "Write file error")
		return
	}
	fmt.Fprintln(w, "Uploadfile success filePath = ", filePath)
}

// 清理目录中一个月前的的临时文件  文件格式 日期_创建时间
func Cleanfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "start clean file ... please wait !!")
	return
	curDay := getCurDay()

	// 读取上传目录下文件名
	_dir, err := ioutil.ReadDir(cf.File_dir)
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
				removeDir := cf.File_dir + dirname
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

func help(w http.ResponseWriter, r *http.Request) {
	ret := `
	/ 访问文件服务
	/Cleanfile 清理文件
	` +
		`/Uploadfile 上传文件 post方法 上传地址：` + cf.File_dir + "\n" +
		`	/Uploadpdb  上传文件 post方法 上传地址：` + cf.Pdb_dir + `
	/Downfile 下载文件 无法直接调用
	/GetSystemInfo 获取系统信息(内存 磁盘 cpu)
	/TestDir  打印文件路径（注意无法重新设置文件服务器根目录，需要重新设置根目录需要重启服务）
	/Reloadconf 重新加载路径
	`

	fmt.Fprintln(w, ret)
}

func main() {

	// fmt.Println(os.Args)
	if len(os.Args) < 1 {
		if err := load_default_config(); err != nil {
			log.Println("init load failed!!! %s", err)
		} else {
			log.Println("init load success", cf.Pdb_dir)
		}
	} else {
		fmt.Println(os.Args[1])
		fmt.Println(reflect.TypeOf(os.Args[1]))
		if err := load_config(os.Args[1]); err != nil {
			log.Println("init load failed!!! %s", err)
		} else {
			log.Println("init load success", cf.Pdb_dir)
		}
	}

	// for i, v := range os.Args {
	// 	fmt.Printf("args[%v]=%v\n", i, v)
	// }

	// var cf conf

	mux := http.NewServeMux()

	//
	mux.HandleFunc("/help", help)

	// 其他接口 清理文件，上传下载文件
	mux.HandleFunc("/Testdir", Testdir)
	mux.HandleFunc("/Reloadconf", Reloadconf)
	mux.HandleFunc("/Cleanfile", Cleanfile)
	mux.HandleFunc("/Uploadpdb", Uploadpdb)
	mux.HandleFunc("/Uploadfile", Uploadfile)
	mux.HandleFunc("/Downfile", Downfile)

	// 获取系统信息
	mux.HandleFunc("/GetSystemInfo", GetSystemInfo)
	mux.HandleFunc("/GetServiceInfo", GetServiceInfo)

	// 其他操作
	// svn更新
	// mux.HandleFunc("UpdateSvn", UpdateSvn)
	// // kill进程
	// mux.HandleFunc("KillPython", KillPython)

	// 文件服务器
	mux.Handle("/", http.FileServer(http.Dir(cf.File_server)))

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
