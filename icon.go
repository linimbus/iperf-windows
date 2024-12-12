package main

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
)

func IconLoadFromBox(filename string, size walk.Size) *walk.Icon {
	body, err := BoxFile().Bytes(filename)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	dir := DEFAULT_HOME + "\\icon\\"
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0644)
		if err != nil {
			logs.Error(err.Error())
			return nil
		}
	}
	filepath := dir + filename
	err = SaveToFile(filepath, body)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	icon, err := walk.NewIconFromFileWithSize(filepath, size)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	return icon
}

var ICON_Main *walk.Icon
var ICON_Status *walk.Icon
var ICON_Start *walk.Icon
var ICON_Stop *walk.Icon
var ICON_Downlink *walk.Icon
var ICON_Uplink *walk.Icon

var ICON_Max_Size = walk.Size{
	Width: 64, Height: 64,
}

var ICON_Min_Size = walk.Size{
	Width: 16, Height: 16,
}

func IconInit() error {
	ICON_Main = IconLoadFromBox("main.ico", ICON_Max_Size)
	ICON_Status = IconLoadFromBox("status.ico", ICON_Min_Size)
	ICON_Start = IconLoadFromBox("start.ico", ICON_Min_Size)
	ICON_Stop = IconLoadFromBox("stop.ico", ICON_Min_Size)
	ICON_Downlink = IconLoadFromBox("downlink.ico", ICON_Min_Size)
	ICON_Uplink = IconLoadFromBox("uplink.ico", ICON_Min_Size)
	return nil
}
