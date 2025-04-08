package iperf3

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
)

func VersionGet() string {
	return "v0.3.2"
}

func SaveToFile(name string, body []byte) error {
	return os.WriteFile(name, body, 0664)
}

func CapSignal(proc func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		proc()
		logs.Error("recv signcal %s, ready to exit", sig.String())
		os.Exit(-1)
	}()
}

func ByteView(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < (1024 * 1024) {
		return fmt.Sprintf("%.1fKB", float64(size)/float64(1024))
	} else if size < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fMB", float64(size)/float64(1024*1024))
	} else if size < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.1fGB", float64(size)/float64(1024*1024*1024))
	} else {
		return fmt.Sprintf("%.1fTB", float64(size)/float64(1024*1024*1024*1024))
	}
}

func InterfaceGet(iface *net.Interface) ([]net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, 0)
	for _, v := range addrs {
		ipone, _, err := net.ParseCIDR(v.String())
		if err != nil {
			continue
		}
		if len(ipone) > 0 {
			ips = append(ips, ipone)
		}
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("interface not any address")
	}
	return ips, nil
}

func InterfaceOptions() []string {
	output := []string{"0.0.0.0", "::"}
	ifaces, err := net.Interfaces()
	if err != nil {
		logs.Error("interface query failed, %s", err.Error())
		return output
	}
	for _, v := range ifaces {
		if v.Flags&net.FlagUp == 0 {
			continue
		}
		address, err := InterfaceGet(&v)
		if err != nil {
			logs.Warning("interface get failed, %s", err.Error())
			continue
		}
		for _, addr := range address {
			output = append(output, addr.String())
		}
	}
	return output
}

func InterfaceIndex(name string, interfaces []string) int {
	for i, v := range interfaces {
		if v == name {
			return i
		}
	}
	return 0
}

func CopyClipboard() (string, error) {
	text, err := walk.Clipboard().Text()
	if err != nil {
		logs.Error(err.Error())
		return "", fmt.Errorf("can not find the any clipboard")
	}
	return text, nil
}

func PasteClipboard(input string) error {
	err := walk.Clipboard().SetText(input)
	if err != nil {
		logs.Error(err.Error())
	}
	return err
}

func GetTimestamp() string {
	return time.Now().Format("2006-01-02T15-04-05")
}

func GetTimestampUS() int64 {
	return time.Now().UnixNano()
}

func FormatJSON(body []byte) ([]byte, error) {
	var data interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(data, "", "    ")
}
