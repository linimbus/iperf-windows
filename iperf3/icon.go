package iperf3

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
)

func IconLoadFromBox(filename string, size walk.Size) *walk.Icon {
	body, err := Asset(filename)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	dir := DEFAULT_HOME + "\\icon\\"
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0644)
		if err != nil {
			logs.Error(err.Error())
			return walk.IconApplication()
		}
	}
	filepath := dir + filename
	err = SaveToFile(filepath, body)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	icon, err := walk.NewIconFromFileWithSize(filepath, size)
	if err != nil {
		logs.Error(err.Error())
		return walk.IconApplication()
	}
	return icon
}

var ICON_Main *walk.Icon
var ICON_Status *walk.Icon
var ICON_Start *walk.Icon
var ICON_Stop *walk.Icon
var ICON_Flow *walk.Icon

func IconInit() {
	ICON_Main = IconLoadFromBox("main.ico", walk.Size{Width: 128, Height: 128})
	ICON_Start = IconLoadFromBox("start.ico", walk.Size{Width: 64, Height: 64})
	ICON_Stop = IconLoadFromBox("stop.ico", walk.Size{Width: 64, Height: 64})
	ICON_Status = IconLoadFromBox("status.ico", walk.Size{Width: 16, Height: 16})
	ICON_Flow = IconLoadFromBox("flow.ico", walk.Size{Width: 16, Height: 16})
}
