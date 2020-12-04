package client

import (
	"ConfigurationTools/configurationManager"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"ConfigurationTools/utils"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"os"
	"strconv"
	"strings"
)

// InitWin 初始化界面
func InitWin() {

	tmw := new(TabMainWindow)

	if err := (MainWindow{
		Title:      "CAN 数据配置工具",
		AssignTo:   &tmw.MainWindow,
		Background: SolidColorBrush{Color: walk.RGB(0xED, 0xED, 0xED)},
		//Background: GradientBrush{
		//	Vertexes: []walk.GradientVertex{
		//		{X: 0, Y: 0, Color: walk.RGB(255, 255, 127)},
		//		{X: 1, Y: 0, Color: walk.RGB(127, 191, 255)},
		//		{X: 0.5, Y: 0.5, Color: walk.RGB(255, 255, 255)},
		//		{X: 1, Y: 1, Color: walk.RGB(127, 255, 127)},
		//		{X: 0, Y: 1, Color: walk.RGB(255, 127, 127)},
		//	},
		//	Triangles: []walk.GradientTriangle{
		//		{0, 1, 2},
		//		{1, 3, 2},
		//		{3, 4, 2},
		//		{4, 0, 2},
		//	},
		//},
		MenuItems: []MenuItem{
			Menu{
				Text: "&工具箱",
				Items: []MenuItem{
					Action{
						Text:        "&指令管理",
						OnTriggered: tmw.openCommandManagePanel,
					},
					Action{
						Text:        "&关联指令",
						OnTriggered: tmw.openRelatedInstructionsPanel,
					},
					Action{
						Text:        "&发送指令",
						OnTriggered: tmw.OpenSendInstructionsPanel,
					},
					Action{
						Text:        "&导出CAN",
						OnTriggered: tmw.exportHistoryCanPanel,
						Visible:     true,
					},
				},
			},
			Menu{
				Text:     "&数据源配置",
				AssignTo: &tmw.configMenu,
				Items: []MenuItem{
					Menu{
						AssignTo: &tmw.targetPlatform,
						Text:     "&目标平台",
					},
				},
			},
			Menu{
				AssignTo: &tmw.helpMenu,
				Text:     "&帮助",
				Items: []MenuItem{
					Action{
						Text: "&帮助文档",
						OnTriggered: func() {
							dir, err := os.Getwd()
							if err == nil {
								dir = strings.ReplaceAll(dir, "\\", "/")
								RunBuiltinWebView(fmt.Sprintf("%s/help.html", dir))
							}
						},
					},
					Action{
						Text:        "&关于",
						OnTriggered: tmw.about,
					},
				},
			},
		},
		Size:   Size{1200, 742},
		Layout: VBox{MarginsZero: true, SpacingZero: true},
		Children: []Widget{
			TabWidget{
				//Background: TransparentBrush{},
				AssignTo: &tmw.TabWidget,
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	tmw.SetX((utils.GetSystemMetrics(0) - tmw.WidthPixels()) / 2)
	tmw.SetY((utils.GetSystemMetrics(1) - tmw.HeightPixels()) / 2)

	icon, err := walk.Resources.Icon("/favicon.ico")
	if err == nil {
		tmw.MainWindow.SetIcon(icon)
	}

	dbKey, err := configurationManager.GetIni(configurationManager.DATABASE)
	if err != nil {
		dbKey = "-1"
	}

	canCfgUrl := utils.KnownConfigUrl()
	for _, item := range canCfgUrl {
		a := walk.NewAction()
		a.SetText(item.Key)
		a.Triggered().Attach(tmw.openCanConfig(item.Value))

		tmw.helpMenu.Actions().Insert(0, a)
	}

	for i := 0; i < len(configurationManager.DBConfig); i++ {

		// 如果初始化值异常，默认初始化第一个选项
		if dbKey == "-1" {
			err := configurationManager.SetIni(configurationManager.DATABASE, strconv.Itoa(configurationManager.DBConfig[i].Code))
			if err != nil {
				mynotify.Error("初始化失败," + err.Error())
				return
			}

			err = configurationManager.SetIni(configurationManager.Model, configurationManager.DBConfig[i].Mode)
			if err != nil {
				mynotify.Error("初始化失败," + err.Error())
				return
			}

			err = configurationManager.SetIni(configurationManager.DATABASE_CONNSTR, configurationManager.DBConfig[i].ConnStr)
			if err != nil {
				mynotify.Error("初始化失败," + err.Error())
				return
			}

			dbKey = strconv.Itoa(configurationManager.DBConfig[i].Code)
		}

		action := walk.NewAction()
		action.SetText(fmt.Sprintf("%s-%s", configurationManager.DBConfig[i].Name, strings.Title(configurationManager.DBConfig[i].Mode)))
		action.SetToolTip(fmt.Sprintf("%s", configurationManager.DBConfig[i].Code))
		action.Triggered().Attach(tmw.setOption(configurationManager.DBConfig[i].Code, configurationManager.DBConfig[i].Mode, configurationManager.DBConfig[i].ConnStr))
		if dbKey == strconv.Itoa(configurationManager.DBConfig[i].Code) {
			tmw.MainWindow.SetTitle(fmt.Sprintf("CAN 数据配置管理 %s", action.Text()))
			action.SetChecked(true)
		}

		tmw.targetPlatform.Actions().Add(action)
	}

	InitVehicleTypePage(tmw)

	tmw.Run()
}

type TabMainWindow struct {
	*walk.MainWindow
	TabWidget      *walk.TabWidget
	targetPlatform *walk.Menu
	configMenu     *walk.Menu

	helpMenu *walk.Menu
}

func InitVehicleTypePage(tmw *TabMainWindow) {
	tp, err := walk.NewTabPage()
	if err != nil {
		log.Fatal(err)
	}
	tp.SetWidth(50)
	if (tp.SetTitle(" 主页 ")); err != nil {
		log.Fatal(err)
	}

	tp.SetLayout(walk.NewHBoxLayout())

	if err := tmw.TabWidget.Pages().Add(tp); err != nil {
		log.Fatal(err)
	}

	if err := tmw.TabWidget.SetCurrentIndex(tmw.TabWidget.Pages().Len() - 1); err != nil {
		log.Fatal(err)
	}

	_, err = NewVehicleTypeList(tp, tmw)
	if err != nil {
		mynotify.Error("窗口初始化失败，" + err.Error())
	}
}

func (mw *TabMainWindow) CanManage(vt *model.VehicleTypeStats) {
	tp := mw.newTab(fmt.Sprintf("【%s】CAN管理", vt.TypeName))
	_, err := NewCanManagePanel(tp, vt, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

// EditGroupCan 为车辆类型删除、更新部分配置信息
func (mw *TabMainWindow) EditGroupCan(vt *model.VehicleTypeStats) {
	tp := mw.newTab(fmt.Sprintf("【%s】批量编辑", vt.TypeName))
	_, err := NewEditCanPanel(tp, vt, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

// AddGroupCan 为车辆类型增加配置信息
func (mw *TabMainWindow) ImportCanFromBigData(vt *model.VehicleTypeStats) {
	tp := mw.newTab(fmt.Sprintf("【%s】从大数据导入", vt.TypeName))
	_, err := ImportCanFromBigdataPanel(tp, vt, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

// EditVehicleType 编辑车型
func (mw *TabMainWindow) EditVehicleType(vt *model.VehicleTypeStats) {
	tp := mw.newTab(fmt.Sprintf("【%s】编辑车型", vt.TypeName))
	_, err := NewVehicleTypeEditPage(tp, vt, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

// ImportCanFromExisting 从已有车型导入CAN
func (mw *TabMainWindow) ImportCanFromExisting(vt *model.VehicleTypeStats) {
	tp := mw.newTab(fmt.Sprintf("【%s】从其他车型导入", vt.TypeName))
	_, err := ImportCanFromExistingPanel(tp, vt, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

func (mw *TabMainWindow) NewImportCanFromCsvFile(vt *model.VehicleTypeStats) {
	tp := mw.newTab(fmt.Sprintf("【%s】从CSV导入", vt.TypeName))
	_, err := NewImportCanFromCsvFilePanel(tp, vt, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

func (mw *TabMainWindow) NewImportCanFromClipboard(vt *model.VehicleTypeStats) {
	tp := mw.newTab(fmt.Sprintf("【%s】从剪贴板导入", vt.TypeName))
	_, err := NewImportCanFromClipboardPanel(tp, vt, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

func (mw *TabMainWindow) AddVehicleType() {
	tp := mw.newTab("添加车型")
	_, err := NewVehicleTypeAddPage(tp, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

func (mw *TabMainWindow) about() {
	walk.MsgBox(mw, "", "前装CAN配置工具\r\n\r\n软件研发部\r\n北京博创联动科技有限公司", walk.MsgBoxIconInformation)
}

// 指令管理
func (mw *TabMainWindow) openCommandManagePanel() {
	tp := mw.newTab("指令管理")

	_, err := NewCommandManagePanel(tp, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

// 关联指令
func (mw *TabMainWindow) openRelatedInstructionsPanel() {
	tp := mw.newTab("关联指令")
	_, err := NewRelatedInstructionsPanel(tp, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

// 发送指令
func (mw *TabMainWindow) OpenSendInstructionsPanel() {
	tp := mw.newTab("发送指令")
	_, err := NewSendInstructionsPanel(tp, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

// 到处CAN历史
func (mw *TabMainWindow) exportHistoryCanPanel() {
	tp := mw.newTab("导出 CAN")

	err := NewExportHistoryCanPanel(tp, mw)
	if err != nil {
		mynotify.Error("窗口初始化失败," + err.Error())
	}
}

func (mw *TabMainWindow) setOption(code int, mode string, url string) func() {

	return func() {
		actions := mw.targetPlatform.Actions()
		for i := 0; i < actions.Len(); i++ {
			if actions.At(i).ToolTip() == strconv.Itoa(code) {
				actions.At(i).SetChecked(true)
				continue
			}
			mw.targetPlatform.Actions().At(i).SetChecked(false)
		}

		err := configurationManager.SetIni(configurationManager.DATABASE, strconv.Itoa(code))
		if err != nil {
			mynotify.Info("设置失败")
			return
		}

		err = configurationManager.SetIni(configurationManager.Model, mode)
		if err != nil {
			mynotify.Info("设置失败")
			return
		}

		err = configurationManager.SetIni(configurationManager.DATABASE_CONNSTR, url)
		if err != nil {
			mynotify.Info("设置失败")
			return
		}

		mynotify.Info("设置成功，重启后生效")
	}
}

func (mw *TabMainWindow) openCanConfig(url string) func() {
	return func() {
		RunBuiltinWebView(url)
	}
}

// 创建新页面
func (mw *TabMainWindow) newTab(tabTitle string) *walk.TabPage {
	tp, err := walk.NewTabPage()
	if err != nil {
		log.Fatal(err)
	}

	if (tp.SetTitle(tabTitle)); err != nil {
		log.Fatal(err)
	}

	menu, err := walk.NewMenu()
	closeAction := walk.NewAction()
	closeAction.SetText("关闭")
	closeAction.Triggered().Attach(func() {
		mw.TabWidget.Pages().Remove(tp)
		tp.Dispose()
	})
	menu.Actions().Add(closeAction)

	tp.SetContextMenu(menu)
	tp.SetLayout(walk.NewHBoxLayout())

	if err := mw.TabWidget.Pages().Add(tp); err != nil {
		log.Fatal(err)
	}

	if err := mw.TabWidget.SetCurrentIndex(mw.TabWidget.Pages().Len() - 1); err != nil {
		log.Fatal(err)
	}

	return tp
}
