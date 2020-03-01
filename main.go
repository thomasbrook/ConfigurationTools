package main

import (
	"ConfigurationTools/client"
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	client.InitWin()
	// mw := new(MyMainWindow)
	// m := MainWindowInit(mw)
	// if _, err := m.Run(); err != nil {
	// 	fmt.Println(err)
	// }
}

type MyMainWindow struct {
	*walk.MainWindow
	sv *walk.ScrollView
}

func MainWindowInit(mw *MyMainWindow) *MainWindow {

	m := &MainWindow{
		AssignTo:   &mw.MainWindow,
		Title:      "MainWindow",
		Background: SolidColorBrush{Color: walk.RGB(255, 0, 0)},
		Layout:     VBox{},
		Size:       Size{550, 340},
		Children: []Widget{
			PushButton{
				Text:       "create",
				Background: TransparentBrush{},
				OnClicked: func() {
					childlist := mw.sv.Children()
					if childlist != nil {
						childlist.Clear()
						//						for i := childlist.Len() - 1; i >= 0; i-- {
						//							ch := childlist.At(i)
						//							ch.SetParent(nil)
						//							ch.Dispose()
						//						}
						fmt.Println("clear ok!")
					}
					loadChild(mw.sv)
					return
				},
			},
			ScrollView{
				AssignTo: &mw.sv,
				Layout: Flow{
					MarginsZero: true,
					Alignment:   AlignHNearVNear,
				},
				Background: TransparentBrush{},
				Children: []Widget{
					Composite{
						Layout:     VBox{},
						Background: TransparentBrush{},
						Children: []Widget{
							Label{
								Alignment: AlignHNearVCenter,
								Font:      Font{PointSize: 11},
								Text:      "hello walk ",
								Visible:   true,
								TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
							},
						},
					},
					Composite{
						Layout:     VBox{},
						Background: TransparentBrush{},
						Children: []Widget{
							Label{
								Alignment: AlignHNearVCenter,
								Font:      Font{PointSize: 11},
								Text:      "I am fine. ",
								Visible:   true,
							},
						},
					},
					Composite{
						Layout:     VBox{},
						Background: TransparentBrush{},
						Children: []Widget{
							Label{
								Alignment: AlignHNearVCenter,
								Font:      Font{PointSize: 11},
								Text:      "I am fine. ",
								Visible:   true,
								TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
							},
						},
					},
					Composite{
						Layout:     VBox{},
						Background: TransparentBrush{},
						Children: []Widget{
							Label{
								Alignment: AlignHNearVCenter,
								Font:      Font{PointSize: 11},
								Text:      "I am fine. ",
								Visible:   true,
								TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
							},
						},
					},
					Composite{
						Layout:     VBox{},
						Background: TransparentBrush{},
						Children: []Widget{
							Label{
								Alignment: AlignHNearVCenter,
								Font:      Font{PointSize: 11},
								Text:      "I am fine. ",
								Visible:   true,
								TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
							},
						},
					},
					Composite{
						Layout:     VBox{},
						Background: TransparentBrush{},
						Children: []Widget{
							Label{
								Alignment: AlignHNearVCenter,
								Font:      Font{PointSize: 11},
								Text:      "I am fine. ",
								Visible:   true,
								TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
							},
						},
					},
				},
			},
		},
	}

	return m
}

func loadChild(parent walk.Container) error {
	for i := 0; i < 8; i++ {
		w := Composite{
			Layout:     VBox{},
			Background: SolidColorBrush{Color: walk.RGB(0x00, 0x8B, 0x8B)},
			Children: []Widget{
				Composite{
					Layout: VBox{},
					Children: []Widget{

						Label{
							Alignment: AlignHNearVCenter,
							Font:      Font{PointSize: 11},
							Text:      "smile",
							Visible:   true,
							TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
						},
					},
				},
			},
		}

		err := w.Create(NewBuilder(parent))
		if err != nil {
			panic(err)
		}
	}

	return nil
}
