package iperf3

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

var _name string
var _home string

var APPLICATION_NAME = "Iperf3Windows"

func RunlogDirGet() string {
	dir := fmt.Sprintf("%s\\runlog\\%s", _home, _name)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	}
	return dir
}

func ConfigDirGet() string {
	dir := fmt.Sprintf("%s\\config\\%s", _home, _name)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	}
	return dir
}

func ToolDirGet() string {
	dir := fmt.Sprintf("%s\\bin", _home)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	}
	return dir
}

func IconDirGet() string {
	dir := fmt.Sprintf("%s\\icon", _home)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	}
	return dir
}

func DataDirGet() string {
	dir := fmt.Sprintf("%s\\data\\%s", _home, _name)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	}
	return dir
}

func appDataDir() string {
	datadir := os.Getenv("APPDATA")
	if datadir == "" {
		datadir = os.Getenv("CD")
	}
	if datadir == "" {
		datadir = ".\\"
	} else {
		datadir = filepath.Join(datadir, APPLICATION_NAME)
	}
	return datadir
}

func appDataDirInit(name string) error {
	dir := appDataDir()
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	}
	_home = dir
	_name = name
	return nil
}

func appInit(file string) {
	body, err := Asset(file)
	if err != nil {
		logs.Error("app init failed, %s", err.Error())
		return
	}
	err = SaveToFile(filepath.Join(ToolDirGet(), file), body)
	if err != nil {
		logs.Error("save to file failed, %s", err.Error())
		return
	}
}

func FileInit(name string) {
	appDataDirInit(name)
	appInit("cygwin1.dll")
	appInit("iperf3.exe")
}
