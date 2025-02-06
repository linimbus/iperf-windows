# iperf-windows

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