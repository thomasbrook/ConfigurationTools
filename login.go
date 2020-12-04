package main

import (
	"ConfigurationTools/client"
	"ConfigurationTools/utils"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

func LoginWin() {

	lw := new(LoginWindow)

	if err := (MainWindow{
		AssignTo: &lw.MainWindow,
		Title:    "CAN 数据配置工具",
		Layout:   HBox{MarginsZero: true, SpacingZero: true},
		Size:     Size{360, 300},

		Background: GradientBrush{
			Vertexes: []walk.GradientVertex{
				{X: 0, Y: 0, Color: walk.RGB(255, 255, 127)},
				{X: 1, Y: 0, Color: walk.RGB(127, 191, 255)},
				{X: 0.5, Y: 0.5, Color: walk.RGB(255, 255, 255)},
				{X: 1, Y: 1, Color: walk.RGB(127, 255, 127)},
				{X: 0, Y: 1, Color: walk.RGB(255, 127, 127)},
			},
			Triangles: []walk.GradientTriangle{
				{0, 1, 2},
				{1, 3, 2},
				{3, 4, 2},
				{4, 0, 2},
			},
		},
		Children: []Widget{
			VSpacer{},
			Composite{
				Layout:    VBox{},
				Alignment: AlignHCenterVCenter,
				MaxSize:   Size{Width: 260},
				Children: []Widget{
					VSpacer{},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							Label{Text: "登录名"},
							LineEdit{AssignTo: &lw.userNameLe, Text: "Admin", Background: TransparentBrush{}},
							Label{Text: "密码"},
							LineEdit{AssignTo: &lw.pwdLe, Text: "123456", PasswordMode: true, Background: TransparentBrush{}},
						},
					},
					Composite{
						Layout: HBox{},
						Children: []Widget{
							Label{AssignTo: &lw.tipLbl, TextColor: walk.RGB(0xFF, 0x00, 0x00)},
							HSpacer{},
							PushButton{
								Text: "登录",
								ContextMenuItems: []MenuItem{
									Menu{Text: "测试"},
								},
								OnClicked: lw.Login,
							},
						},
					},
					VSpacer{},
				},
			},
			VSpacer{},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	lw.SetX((utils.GetSystemMetrics(0) - lw.WidthPixels()) / 2)
	lw.SetY((utils.GetSystemMetrics(1) - lw.HeightPixels()) / 2)

	icon, err := walk.Resources.Icon("/favicon.ico")
	if err == nil {
		lw.MainWindow.SetIcon(icon)
	}

	lw.Run()
}

type LoginWindow struct {
	*walk.MainWindow

	userNameLe *walk.LineEdit
	pwdLe      *walk.LineEdit
	tipLbl     *walk.Label
}

func (lw *LoginWindow) Login() {

	//if isConn := dataSource.Ping(); !isConn {
	//	mynotify.Error("数据库无法连接")
	//	return
	//}

	uname := lw.userNameLe.Text()
	pwd := lw.pwdLe.Text()

	fmt.Println(uname, pwd)
	if uname != "Admin" && pwd != "123456" {
		lw.tipLbl.SetText("用户名或密码错误")
		return
	}

	lw.tipLbl.SetText("")

	lw.Close()
	client.InitWin()
}
