package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/mynotify"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

func NewCommandManagePanel(parent walk.Container, mainWin *TabMainWindow) (*CommandManagePage, error) {
	cmp := &CommandManagePage{
		mainWin: mainWin,
	}

	if err := (Composite{
		AssignTo: &cmp.Composite,
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true, Margins: Margins{Right: 3}},
				Children: []Widget{
					HSpacer{},
					LineEdit{
						AssignTo:    &cmp.searchLe,
						ToolTipText: "请输入指令名称或描述关键词",
						MaxSize:     Size{Width: 180},
						OnKeyDown: func(key walk.Key) {
							if key == walk.KeyReturn {
								cmp.loadCommand(cmp.searchLe.Text())
							}
						},
					},
					PushButton{
						Text: "查询",
						OnClicked: func() {
							cmp.loadCommand(cmp.searchLe.Text())
						},
					},
					PushButton{
						Text: "添加",
						OnClicked: func() {
							var cmdcfg commandConfig
							if cmd, err := cmdcfg.runCommandDialog(cmp.mainWin); err != nil {
								walk.MsgBox(cmp.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
								return
							} else if cmd == walk.DlgCmdOK {
								cce := &dataSource.CommandConfigEntity{
									CmdName:        cmdcfg.CommandName,
									Url:            cmdcfg.Url,
									Desc:           cmdcfg.Description,
									RenderTemplate: cmdcfg.RenderTemplate,
									ParamTemplate:  cmdcfg.ParamTemplate,
									GroupType:      cmdcfg.GroupType,
								}
								_, err := cce.AddCommand()
								if err != nil {
									walk.MsgBox(cmp.mainWin, "添加失败", err.Error(), walk.MsgBoxIconError)
									return
								}
							}
						},
					},

					PushButton{
						Text: "删除选定项",
						OnClicked: func() {
							idx := cmp.commandTbl.CurrentIndex()
							if idx == -1 {
								walk.MsgBox(cmp.mainWin, "提示", "没有选定项", walk.MsgBoxIconWarning)
								return
							}

							var dlg *walk.Dialog
							var accepPB, cancelPB *walk.PushButton
							dlgcode, err := Dialog{
								AssignTo:      &dlg,
								Title:         "确认框",
								DefaultButton: &accepPB,
								CancelButton:  &cancelPB,

								MinSize: Size{Width: 300},
								Layout:  VBox{},
								Children: []Widget{
									Label{Text: "是否确认删除"},
									Composite{
										Layout: HBox{},
										Children: []Widget{
											HSpacer{},
											PushButton{
												AssignTo: &accepPB,
												Text:     "确认",
												OnClicked: func() {
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
							}.Run(cmp.mainWin)
							if err != nil {
								mynotify.Error("确认窗口初始化失败：" + err.Error())
								return
							}
							if dlgcode == walk.DlgCmdCancel || dlgcode == walk.DlgCmdNone {
								return
							}

							entities := cmp.commandTbl.Model().([]*dataSource.CommandConfigEntity)
							_, err = entities[idx].DeleteCommand()
							if err != nil {
								walk.MsgBox(cmp.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
								return
							}
							cmp.loadCommand(cmp.searchLe.Text())
						},
					},
				},
			},
			TableView{
				AssignTo:         &cmp.commandTbl,
				AlternatingRowBG: true,
				//AlternatingRowBGColor: walk.RGB(239, 239, 239),
				ColumnsOrderable: true,
				Columns: []TableViewColumn{
					{Name: "Id", Title: "#", Frozen: true, Width: 35},
					{Name: "CmdName", Title: "指令名称", Width: 120},
					{Name: "Url", Title: "URL", Width: 120},
					{Name: "Desc", Title: "描述", Width: 180},
					{Name: "RenderTemplate", Title: "渲染模板", Width: 260},
					{Name: "ParamTemplate", Title: "参数模板", Width: 220},
					{Name: "DefaultValue", Title: "默认值", Width: 260},
					{Name: "GroupType", Title: "指令分组", Width: 90, FormatFunc: func(value interface{}) string {
						switch value.(int) {
						case 0:
							return ""
						default:
							return ""
						}
					}},
					{Name: "CreateTime", Title: "创建时间", Width: 120, FormatFunc: func(value interface{}) string {
						if value != nil {
							return (value.(time.Time)).Format("2006/01/02 15:04:05")
						}
						return ""
					}},
				},
				OnItemActivated: func() {
					idx := cmp.commandTbl.CurrentIndex()
					entities := cmp.commandTbl.Model().([]*dataSource.CommandConfigEntity)

					cmdcfg := &commandConfig{
						Id:             entities[idx].Id,
						CommandName:    entities[idx].CmdName,
						Url:            entities[idx].Url,
						Description:    entities[idx].Desc,
						RenderTemplate: entities[idx].RenderTemplate,
						ParamTemplate:  entities[idx].ParamTemplate,
						DefaultValue:   entities[idx].DefaultValue,
						GroupType:      entities[idx].GroupType,
					}
					if cmd, err := cmdcfg.runCommandDialog(cmp.mainWin); err != nil {
						walk.MsgBox(cmp.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
						return
					} else if cmd == walk.DlgCmdOK {

						entities[idx].CmdName = cmdcfg.CommandName
						entities[idx].Url = cmdcfg.Url
						entities[idx].Desc = cmdcfg.Description
						entities[idx].RenderTemplate = cmdcfg.RenderTemplate
						entities[idx].ParamTemplate = cmdcfg.ParamTemplate
						entities[idx].DefaultValue = cmdcfg.DefaultValue
						entities[idx].GroupType = cmdcfg.GroupType

						_, err := entities[idx].UpdateCommand()
						if err != nil {
							walk.MsgBox(cmp.mainWin, "更新失败", err.Error(), walk.MsgBoxIconError)
							return
						}
					}
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true, Margins: Margins{Top: 3}},
				Children: []Widget{
					Label{
						AssignTo: &cmp.statusLbl,
						Font:     Font{PointSize: 10},
					},
					HSpacer{},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	cmp.loadCommand("")
	return cmp, nil
}

type CommandManagePage struct {
	*walk.Composite
	mainWin *TabMainWindow

	searchLe   *walk.LineEdit
	commandTbl *walk.TableView
	statusLbl  *walk.Label
}

func (cmp *CommandManagePage) loadCommand(searchStr string) {
	data, err := (&dataSource.CommandConfigEntity{CmdName: searchStr}).ListCommand()
	if err != nil {
		walk.MsgBox(cmp.mainWin, "加载失败", err.Error(), walk.MsgBoxIconError)
	} else {
		cmp.commandTbl.SetModel(data)
	}
}

type commandConfig struct {
	Id             int
	CommandName    string
	Url            string
	Description    string
	RenderTemplate string
	ParamTemplate  string
	DefaultValue   string
	GroupType      int
}

func (cmd *commandConfig) runCommandDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	var pro Property
	if cmd.Id == 0 {
		pro = Bind("'添加指令'")
	} else {
		pro = Bind("'指令详情' + (commandConfig.CommandName == '' ? '' : ' - '+commandConfig.CommandName)")
	}

	return Dialog{
		AssignTo:      &dlg,
		Title:         pro,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "commandConfig",
			DataSource:     cmd,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{500, 460},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "指令名称",
					},
					LineEdit{
						Text:      Bind("CommandName", Regexp{"^[\\s|\\S]{1,}$"}),
						MaxLength: 45,
					},
					Label{
						Text: "URL",
					},
					LineEdit{
						Text: Bind("Url"),
					},
					Label{
						Text: "指令描述",
					},
					LineEdit{
						Text:      Bind("Description"),
						MaxLength: 255,
					},
					Label{
						Text: "渲染模板",
					},
					TextEdit{
						Text:      Bind("RenderTemplate", Regexp{"^([^|,，]+\\|[^|,，]+\\|[^|,，]+[,，]?)+$"}),
						MaxSize:   Size{Height: 60},
						MinSize:   Size{Height: 60},
						MaxLength: 255,
					},
					Label{
						Text: "参数模板",
					},
					TextEdit{
						Text:      Bind("ParamTemplate", Regexp{"^[\\s|\\S]{1,}$"}),
						MaxSize:   Size{Height: 60},
						MinSize:   Size{Height: 60},
						MaxLength: 255,
					},
					Label{
						Text: "默认值",
					},
					TextEdit{
						Text:      Bind("DefaultValue"), //Regexp{"^([^=,，]+=[^=,，]+[,，]?)+$"}
						MaxSize:   Size{Height: 60},
						MinSize:   Size{Height: 60},
						MaxLength: 255,
					},
				},
			},
			VSpacer{},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "保存",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
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
