package iperf3

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/astaxie/beego/logs"
)

type Config struct {
	ServerListen      string
	ServerPort        int
	ServerCount       int
	ServerInterval    int
	ServerLog         string
	ServerAutoStartup bool
	ServerAutoHide    bool
	ServerJsonFormat  bool

	ClientListen            string
	ClientAddress           string
	ClientPort              int
	ClientRunTime           int
	ClientOmitSec           int
	ClientProtocol          string
	ClientPayload           int
	ClientJsonFormat        bool
	ClientDontFragment      bool
	ClientZeroCopy          bool
	ClientNoDelay           bool
	ClientReverseMode       bool
	ClientBidirectionalMode bool
	ClientMaxmumSegment     bool
	ClientOnlyIPv4          bool
	ClientOnlyIPv6          bool
	ClientStreams           int
	ClientBandwidth         int
	ClientBandwidthUnit     string // KB,MB,GB
	ClientWindows           int
	ClientWindowsUnit       string
	ClientDscpValue         int
	ClientTypeService       int
	ClientRepeatCount       int
	ClientRepeatInterval    int
	ClientLog               string
}

var configCache = Config{
	ServerListen:      "0.0.0.0",
	ServerPort:        5201,
	ServerCount:       1,
	ServerInterval:    1,
	ServerLog:         "",
	ServerAutoStartup: false,
	ServerAutoHide:    false,
	ServerJsonFormat:  true,

	ClientListen:            "0.0.0.0",
	ClientAddress:           "127.0.0.1",
	ClientPort:              5201,
	ClientRunTime:           10,
	ClientOmitSec:           0,
	ClientProtocol:          "tcp",
	ClientPayload:           1024,
	ClientJsonFormat:        true,
	ClientDontFragment:      false,
	ClientZeroCopy:          false,
	ClientNoDelay:           false,
	ClientReverseMode:       false,
	ClientBidirectionalMode: false,
	ClientMaxmumSegment:     false,
	ClientOnlyIPv4:          false,
	ClientOnlyIPv6:          false,
	ClientStreams:           1,
	ClientBandwidth:         0,
	ClientBandwidthUnit:     "MB",
	ClientWindows:           0,
	ClientWindowsUnit:       "MB",
	ClientDscpValue:         0,
	ClientTypeService:       0,
	ClientRepeatCount:       1,
	ClientRepeatInterval:    0,
	ClientLog:               "",
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

	configFilePath = filepath.Join(ConfigDirGet(), name+".json")
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
