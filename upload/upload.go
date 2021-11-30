package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const BaseUploadPath = "D:\\gitProject\\fileserver\\upload\\test"

func handleUpload(w http.ResponseWriter, request *http.Request) {
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
	newFile, err := os.Create(BaseUploadPath + "/" + fileHeader.Filename)
	if err != nil {
		_, _ = io.WriteString(w, "Create file error")
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
	_, _ = io.WriteString(w, "Upload success")
}

func main() {
	http.HandleFunc("/upload", handleUpload)
	// http.HandleFunc("/download", handleDownload)

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal("Server run fail")
	}
}
