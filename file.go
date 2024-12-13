package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
)

var DEFAULT_HOME string

func RunlogDirGet() string {
	dir := fmt.Sprintf("%s\\runlog", DEFAULT_HOME)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	}
	return dir
}

func ConfigDirGet() string {
	dir := fmt.Sprintf("%s\\config", DEFAULT_HOME)
	_, err := os.Stat(dir)
	if err != nil {
		os.MkdirAll(dir, 0644)
	}
	return dir
}

func ToolDirGet() string {
	dir := fmt.Sprintf("%s\\bin", DEFAULT_HOME)
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
		datadir = fmt.Sprintf("%s\\Iperf3Windows", datadir)
	}
	return datadir
}

func appDataDirInit() error {
	dir := appDataDir()
	_, err := os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0644)
		if err != nil {
			return err
		}
	}
	DEFAULT_HOME = dir
	return nil
}

func appInit(file string) error {
	body, err := BoxFile().Bytes(file)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	err = SaveToFile(filepath.Join(ToolDirGet(), file), body)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	return nil
}

func FileInit() error {
	err := appDataDirInit()
	if err != nil {
		return err
	}
	appInit("cygwin1.dll")
	// if err != nil {
	// 	return err
	// }
	appInit("iperf3.exe")
	// if err != nil {
	// 	return err
	// }
	return nil
}
