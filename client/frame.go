package client

import (
	"ConfigurationTools/model"
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// InitWin 初始化界面
func InitWin() {

	tmw := new(TabMainWindow)

	if err := (MainWindow{
		Title:      "农机前装配置管理工具",
		AssignTo:   &tmw.MainWindow,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		MenuItems: []MenuItem{
			Menu{
				Text: "&工具",
				Items: []MenuItem{
					Action{
						Text:        "导出CAN历史",
						OnTriggered: tmw.newCanExportPanel,
					},
				},
			},
			Menu{
				Text: "&帮助",
				Items: []MenuItem{
					Action{
						Text:        "关于",
						OnTriggered: tmw.about,
					},
				},
			},
		},
		Size:   Size{1100, 680},
		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			TabWidget{
				AssignTo: &tmw.TabWidget,
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	//	brush, err := walk.NewSystemColorBrush(walk.SysColorWindow)
	//	if err == nil {
	//		tmw.ToolBar().SetBackground(brush)
	//		tmw.MainWindow.SetBackground(brush)
	//	}

	InitVehicleTypePage(tmw)

	tmw.Run()
}

type TabMainWindow struct {
	*walk.MainWindow
	TabWidget *walk.TabWidget
}

func InitVehicleTypePage(tmw *TabMainWindow) {

	tp, err := walk.NewTabPage()
	if err != nil {
		log.Fatal(err)
	}

	if (tp.SetTitle("首页")); err != nil {
		log.Fatal(err)
	}

	//	tb := TransparentBrush{}
	//	b, err := tb.Create()
	//	if err == nil {
	//		tp.SetBackground(b)
	//	}

	tp.SetLayout(walk.NewHBoxLayout())

	_, err1 := NewVehicleTypeList(tp, tmw)
	if err1 != nil {
		log.Fatal(err1)
	}

	if err := tmw.TabWidget.Pages().Add(tp); err != nil {
		log.Fatal(err)
	}

	if err := tmw.TabWidget.SetCurrentIndex(tmw.TabWidget.Pages().Len() - 1); err != nil {
		log.Fatal(err)
	}
}

// AddGroupCan 为车辆类型删除、更新部分配置信息
func (mw *TabMainWindow) EditGroupCan(vt *model.VehicleType) {

	tp, err := walk.NewTabPage()
	if err != nil {
		log.Fatal(err)
	}

	if (tp.SetTitle(vt.TypeName)); err != nil {
		log.Fatal(err)
	}

	menu, err := walk.NewMenu()
	closeAction := walk.NewAction()
	closeAction.SetText("关闭")
	closeAction.Triggered().Attach(func() {
		mw.TabWidget.Pages().Remove(tp)
	})
	menu.Actions().Add(closeAction)

	tp.SetContextMenu(menu)
	tp.SetLayout(walk.NewHBoxLayout())

	_, err = NewEditPanel(tp, vt, mw)
	if err != nil {
		log.Fatal(err)
	}

	if err := mw.TabWidget.Pages().Add(tp); err != nil {
		log.Fatal(err)
	}

	if err := mw.TabWidget.SetCurrentIndex(mw.TabWidget.Pages().Len() - 1); err != nil {
		log.Fatal(err)
	}
}

// AddGroupCan 为车辆类型增加配置信息
func (mw *TabMainWindow) AddGroupCan(vt *model.VehicleType) {

	tp, err := walk.NewTabPage()
	if err != nil {
		log.Fatal(err)
	}

	if (tp.SetTitle(vt.TypeName)); err != nil {
		log.Fatal(err)
	}

	menu, err := walk.NewMenu()
	closeAction := walk.NewAction()
	closeAction.SetText("关闭")
	closeAction.Triggered().Attach(func() {
		mw.TabWidget.Pages().Remove(tp)
	})
	menu.Actions().Add(closeAction)

	tp.SetContextMenu(menu)
	tp.SetLayout(walk.NewHBoxLayout())

	_, err = NewAddPanel(tp, vt, mw)
	if err != nil {
		log.Fatal(err)
	}

	if err := mw.TabWidget.Pages().Add(tp); err != nil {
		log.Fatal(err)
	}

	if err := mw.TabWidget.SetCurrentIndex(mw.TabWidget.Pages().Len() - 1); err != nil {
		log.Fatal(err)
	}
}

func (mw *TabMainWindow) about() {
	walk.MsgBox(mw, "关于", "前装 CAN 配置工具\r\n作者：杜建平\r\n邮箱：dujianping@uml-tech.com", walk.MsgBoxIconInformation)
}

func (mw *TabMainWindow) newCanExportPanel() {
	tp, err := walk.NewTabPage()
	if err != nil {
		log.Fatal(err)
	}

	if (tp.SetTitle("导出CAN")); err != nil {
		log.Fatal(err)
	}

	menu, err := walk.NewMenu()
	closeAction := walk.NewAction()
	closeAction.SetText("关闭")
	closeAction.Triggered().Attach(func() {
		mw.TabWidget.Pages().Remove(tp)
	})
	menu.Actions().Add(closeAction)

	tp.SetContextMenu(menu)
	tp.SetLayout(walk.NewHBoxLayout())

	err = NewExportPanel(tp, mw)
	if err != nil {
		log.Fatal(err)
	}

	if err := mw.TabWidget.Pages().Add(tp); err != nil {
		log.Fatal(err)
	}

	if err := mw.TabWidget.SetCurrentIndex(mw.TabWidget.Pages().Len() - 1); err != nil {
		log.Fatal(err)
	}
}
