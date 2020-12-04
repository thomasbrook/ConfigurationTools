package client

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

type editForm struct {
	form      *walk.Composite
	ConfigUrl string
	SearchStr string
}

func (ef *editForm) runExportDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	icon, _ := walk.Resources.Icon("/img/search1.ico")

	return Dialog{
		AssignTo:      &dlg,
		Title:         "CAN关键词(字段名、编码)搜索，多个换行",
		Icon:          icon,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:   &db,
			DataSource: ef,
		},
		MinSize: Size{360, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					TextEdit{
						MinSize:     Size{100, 50},
						Text:        Bind("SearchStr"),
						VScroll:     true,
						ToolTipText: "输入搜索关键词（编码、字段名），多个换行显示",
					},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "查询",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "取消",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}
