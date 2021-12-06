import requests

file = {
    "sample_file": open("D:\\gitProject\\fileserver\\upload\\1.txt", "rb"),
    "Content-Type": "application/octet-stream",
    "Content-Disposition": "form-data",
    "filename" : "1.txt"
}

headers = {
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3",
    "Accept-Encoding": "gzip, deflate",
    "Accept-Language": "zh-CN,zh;q=0.9",
    "Cache-Control": "max-age=0",
    "Connection": "keep-alive",
    # "Content-Type": "multipart/form-data",
    # "Host": "10.222.222.7",
    # "Origin": "http://10.222.222.7",
    # "Referer": "http://10.222.222.7/src/html.php/html/system_samples",
    "Upgrade-Insecure-Requests": "1",
    "User-Agent": "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.87 Safari/537.36"
}

filepath="D:\\gitProject\\fileserver\\upload\\1.txt"
# requests.head
# data=filepath#{"file":filepath}
data = {
    "sample_name" : "1.txt",
    "owner_group" : "/data/atp/pcap/custom/test",
    "type" : "1",
    "sample_file_path" : "",
    "description_file_path" : "",
    # "description_file":""
}
# response = requests.post("http://localhost:3000/upload", headers=headers, files=file, data=data)
files = {'file': open(filepath, 'rb')}
response = requests.post("http://localhost:8081/Uploadfile", files=files)
print(response.status_code)