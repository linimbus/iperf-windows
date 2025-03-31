# iperf-windows

[![Download IPerf3-Windows](https://img.shields.io/sourceforge/dm/iperf3-windows.svg)](https://sourceforge.net/projects/iperf3-windows/files/latest/download)

[![Download IPerf3-Windows](https://a.fsdn.com/con/app/sf-download-button)](https://sourceforge.net/projects/iperf3-windows/files/latest/download)

### 准备环境

- Windows10/11 64位
- [Golang SDK](https://studygolang.com/dl/golang/go1.23.3.windows-amd64.msi)

### 环境变量

- GOPROXY=<https://goproxy.cn,direct>
- GOPATH=D:\workspace\golang

### 准备工具

```
go install github.com/akavel/rsrc@latest
go install github.com/jteeuwen/go-bindata/...@latest
```

### 编译

```
.\build.bat
```

### 测试
