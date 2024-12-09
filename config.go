package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Config struct {
	ConnectivityURL string
	FilterInterface string
	OutputDirectory string
	RestfulHeader   string
	RestfulMethod   string
	RestfulURL      string
	Interval        int
}

var configCache = Config{
	ConnectivityURL: "https://test.ipw.cn/",
	FilterInterface: "",
	OutputDirectory: "",
	RestfulHeader:   "",
	RestfulMethod:   "POST",
	RestfulURL:      "",
	Interval:        60,
}

var configFilePath string
var configLock sync.Mutex

func configSyncToFile() error {
	configLock.Lock()
	defer configLock.Unlock()

	value, err := json.MarshalIndent(configCache, "\t", " ")
	if err != nil {
		logs.Error("json marshal config fail, %s", err.Error())
		return err
	}
	return os.WriteFile(configFilePath, value, 0664)
}

func ConfigGet() *Config {
	return &configCache
}

func ConnectivityURLSave(url string) error {
	configCache.ConnectivityURL = url
	return configSyncToFile()
}

func FilterInterfaceSave(filter string) error {
	configCache.FilterInterface = filter
	return configSyncToFile()
}

func OutputDirectorySave(dir string) error {
	configCache.OutputDirectory = dir
	return configSyncToFile()
}

func RestfulHeaderSave(key, value string) error {
	configCache.RestfulHeader = fmt.Sprintf("%s:%s", key, value)
	return configSyncToFile()
}

func RestfulMethodSave(method string) error {
	configCache.RestfulMethod = method
	return configSyncToFile()
}

func RestfulURLSave(url string) error {
	configCache.RestfulURL = url
	return configSyncToFile()
}

func IntervalSave(value int) error {
	configCache.Interval = value
	return configSyncToFile()
}

func ConfigInit() error {
	configFilePath = fmt.Sprintf("%s%c%s", ConfigDirGet(), os.PathSeparator, "config.json")

	_, err := os.Stat(configFilePath)
	if err != nil {
		err = configSyncToFile()
		if err != nil {
			logs.Error("config sync to file fail, %s", err.Error())
			return err
		}
	}

	value, err := os.ReadFile(configFilePath)
	if err != nil {
		logs.Error("read config file from app data dir fail, %s", err.Error())
		configSyncToFile()

		return err
	}

	err = json.Unmarshal(value, &configCache)
	if err != nil {
		logs.Error("json unmarshal config fail, %s", err.Error())
		configSyncToFile()

		return err
	}

	return nil
}
