package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func VersionGet() string {
	return "v1.0.0"
}

func main() {
	var mainWin *walk.MainWindow

	var consoleMode *walk.ComboBox
	var consoleAddress *walk.ComboBox
	var consolePort *walk.NumberEdit
	var consoleProtocol *walk.ComboBox

	var buttonStart, buttonStop *walk.PushButton

	_, err := MainWindow{
		Title:    "iperf-windows " + VersionGet(),
		AssignTo: &mainWin,
		MinSize:  Size{Width: 300, Height: 200},
		Size:     Size{Width: 300, Height: 200},
		Layout: VBox{
			Alignment:   AlignHNearVNear,
			MarginsZero: true,
			Margins:     Margins{Left: 10, Top: 10, Right: 10, Bottom: 10},
		},
		MenuItems: []MenuItem{
			Action{
				Text: "Setting",
				OnTriggered: func() {

				},
			},
			Action{
				Text: "Mini Windows",
				OnTriggered: func() {

				},
			},
			Action{
				Text: "About",
				OnTriggered: func() {
					AboutAction(mainWin)
				},
			},
		},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "Mode: ",
					},
					ComboBox{
						AssignTo:     &consoleMode,
						CurrentIndex: 0,
						Model:        []string{"Server", "Client"},
						OnCurrentIndexChanged: func() {
						},
					},
					Label{
						Text: "Address:",
					},
					ComboBox{
						AssignTo:     &consoleAddress,
						CurrentIndex: 0,
						Model:        []string{""},
						OnCurrentIndexChanged: func() {
						},
					},
					Label{
						Text: "Port:",
					},
					NumberEdit{
						AssignTo:    &consolePort,
						Value:       float64(5201),
						ToolTipText: "1~65535",
						MaxValue:    65535,
						MinValue:    1,
						OnValueChanged: func() {
							// addLink.Port = int(consolePort.Value())
						},
					},

					Label{
						Text: "Protocol:",
					},
					ComboBox{
						AssignTo:     &consoleProtocol,
						CurrentIndex: 0,
						Model:        []string{"tcp", "udp"},
						OnCurrentIndexChanged: func() {
							// addLink.Protocol = consoleProtocol.Text()
						},
					},
					VSpacer{},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{

							PushButton{
								AssignTo: &buttonStart,
								Text:     "Start",
								OnClicked: func() {

								},
							},
							PushButton{
								AssignTo: &buttonStop,
								Text:     "Stop",
								OnClicked: func() {

								},
							},
						},
					},
				},
			},
		},
	}.Run()

	if err != nil {
		log.Fatalf(err.Error())
	} else {
		mainWin.Close()
	}
}
