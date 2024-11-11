package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func AboutAction(mw *walk.MainWindow) {
	var ok *walk.PushButton
	var about *walk.Dialog
	var err error

	_, err = Dialog{
		AssignTo:      &about,
		Title:         "About",
		Icon:          walk.IconInformation(),
		MinSize:       Size{Width: 300, Height: 200},
		DefaultButton: &ok,
		Layout:        VBox{},
		Children: []Widget{
			TextLabel{
				Text:    "helloworld",
				MinSize: Size{Width: 250, Height: 200},
				MaxSize: Size{Width: 290, Height: 400},
			},
			Label{
				Text:          "Version: " + VersionGet(),
				TextAlignment: AlignCenter,
			},
			VSpacer{
				MinSize: Size{Height: 10},
			},
			PushButton{
				Text:      "OK",
				OnClicked: func() { about.Cancel() },
			},
		},
	}.Run(mw)

	if err != nil {
		log.Println(err.Error())
	}
}
