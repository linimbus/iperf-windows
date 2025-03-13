package iperf3

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Config struct {
	ServerPort        int
	ServerInterval    int
	ServerLog         string
	ServerAutoStartup bool
	ServerAutoHide    bool
	ServerJsonFormat  bool

	ClientAddress       string
	ClientPort          int
	ClientRunTime       int
	ClientOmitSec       int
	ClientProtocol      string
	ClientPayload       int
	ClientJsonFormat    bool
	ClientDontFragment  bool
	ClientZeroCopy      bool
	ClientNoDelay       bool
	ClientStreams       int
	ClientBandwidth     int
	ClientBandwidthUnit string // KB,MB,GB
	ClientDscp          int
	ClientLog           string
}

var configCache = Config{
	ServerPort:        5201,
	ServerInterval:    1,
	ServerLog:         "",
	ServerAutoStartup: false,
	ServerAutoHide:    false,
	ServerJsonFormat:  true,

	ClientAddress:       "127.0.0.1",
	ClientPort:          5201,
	ClientRunTime:       10,
	ClientOmitSec:       0,
	ClientProtocol:      "tcp",
	ClientPayload:       1024,
	ClientJsonFormat:    true,
	ClientDontFragment:  false,
	ClientZeroCopy:      false,
	ClientNoDelay:       false,
	ClientStreams:       1,
	ClientBandwidth:     0,
	ClientBandwidthUnit: "MB",
	ClientDscp:          0,
	ClientLog:           "",
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

func ConfigInit(name string) {

	defer func() {
		if configCache.ClientLog == "" {
			configCache.ClientLog = DataDirGet()
		}

		if configCache.ServerLog == "" {
			configCache.ServerLog = DataDirGet()
		}

		configSyncToFile()
	}()

	configFilePath = fmt.Sprintf("%s%c%s", ConfigDirGet(), os.PathSeparator, name+".json")
	_, err := os.Stat(configFilePath)
	if err != nil {
		err = configSyncToFile()
		if err != nil {
			logs.Error("config sync to file fail, %s", err.Error())
			return
		}
	}
	value, err := os.ReadFile(configFilePath)
	if err != nil {
		logs.Error("read config file from app data dir fail, %s", err.Error())
		configSyncToFile()
		return
	}
	err = json.Unmarshal(value, &configCache)
	if err != nil {
		logs.Error("json unmarshal config fail, %s", err.Error())
		configSyncToFile()
		return
	}
}
