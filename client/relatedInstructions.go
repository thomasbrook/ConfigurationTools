package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/mynotify"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

func NewRelatedInstructionsPanel(parent walk.Container, mainWin *TabMainWindow) (*RelatedInstructionsPage, error) {
	cmdCfg, err := loadCommandConfig()
	if err != nil {
		walk.MsgBox(mainWin, "", err.Error(), walk.MsgBoxIconError)
		cmdCfg = new(CommandConfigTable)
	}

	rip := &RelatedInstructionsPage{
		mainWin:         mainWin,
		commandTblModel: cmdCfg,
	}

	terType, err := new(dataSource.TerminalTypeEntity).ListTerminalType()
	if err != nil {
		walk.MsgBox(mainWin, "", err.Error(), walk.MsgBoxIconError)
		terType = []*dataSource.TerminalTypeEntity{}
	}

	if err := (Composite{
		AssignTo: &rip.Composite,
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					ComboBox{
						AssignTo:      &rip.terTypeCb,
						BindingMember: "Id",
						DisplayMember: "TypeName",
						ToolTipText:   "请选择需要关联指令的终端类型",
						MinSize:       Size{Width: 160},
						Model:         terType,
						OnCurrentIndexChanged: func() {
							currIdx := rip.terTypeCb.CurrentIndex()
							if currIdx == -1 {
								return
							}

							m := rip.terTypeCb.Model().([]*dataSource.TerminalTypeEntity)
							relId, err := m[currIdx].ListCmdConfigRelId()
							if err != nil {
								walk.MsgBox(rip.mainWin, "", err.Error(), walk.MsgBoxIconError)
								return
							}

							relIdMap := make(map[int]bool)
							for _, val := range relId {
								relIdMap[val] = true
							}

							for _, item := range rip.commandTblModel.items {
								_, isExist := relIdMap[item.Id]
								if isExist {
									item.IsChecked = true
								} else {
									item.IsChecked = false
								}
							}
							rip.commandTblModel.PublishRowsReset()
						},
					},
					PushButton{
						Text: "关联",
						OnClicked: func() {
							currIdx := rip.terTypeCb.CurrentIndex()
							if currIdx == -1 {
								walk.MsgBox(rip.mainWin, "提示", "请选择终端类型", walk.MsgBoxIconInformation)
								return
							}
							m := rip.terTypeCb.Model().([]*dataSource.TerminalTypeEntity)

							var relId []int
							for _, val := range rip.commandTblModel.items {
								if val.IsChecked {
									relId = append(relId, val.Id)
								}
							}

							if len(relId) == 0 {
								walk.MsgBox(rip.mainWin, "提示", "没有选定项", walk.MsgBoxIconInformation)
								return
							}

							err = m[currIdx].RelateCommandConfig(relId)
							if err != nil {
								walk.MsgBox(rip.mainWin, "错误", err.Error(), walk.MsgBoxIconInformation)
								return
							}
							mynotify.Message("关联成功")
						},
					},
				},
			},
			TableView{
				AssignTo:         &rip.commandTbl,
				AlternatingRowBG: true,
				//AlternatingRowBGColor: walk.RGB(239, 239, 239),
				ColumnsOrderable: false,
				CheckBoxes:       true,
				MultiSelection:   true,
				Columns: []TableViewColumn{
					{Name: "Id", Title: "#", Frozen: true, Width: 35},
					{Name: "CmdName", Title: "指令名称", Width: 120},
					{Name: "Desc", Title: "描述", Width: 180},
					{Name: "Param", Title: "参数模板", Width: 260},
					{Name: "DataRange", Title: "数值范围", Width: 220},
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
				Model: rip.commandTblModel,
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	return rip, nil
}

type RelatedInstructionsPage struct {
	*walk.Composite
	mainWin *TabMainWindow

	terTypeCb       *walk.ComboBox
	commandTbl      *walk.TableView
	commandTblModel *CommandConfigTable
}

type CommandConfigTable struct {
	walk.TableModelBase
	items []*commandConfigModel
}

type commandConfigModel struct {
	Index          int
	Id             int
	CmdName        string
	Desc           string
	RenderTemplate string
	ParamTemplate  string
	GroupType      int
	ControlType    string
	CreateTime     time.Time
	IsChecked      bool
}

func loadCommandConfig() (*CommandConfigTable, error) {
	m := new(CommandConfigTable)

	data, err := new(dataSource.CommandConfigEntity).ListCommand()
	if err != nil {
		return nil, err
	}
	for idx, val := range data {
		m.items = append(m.items, &commandConfigModel{
			Index:          idx + 1,
			Id:             val.Id,
			CmdName:        val.CmdName,
			Desc:           val.Desc,
			RenderTemplate: val.RenderTemplate,
			ParamTemplate:  val.ParamTemplate,
			GroupType:      val.GroupType,
			CreateTime:     val.CreateTime,
			IsChecked:      false,
		})
	}
	m.PublishRowsReset()
	return m, nil
}

func (cct *CommandConfigTable) RowCount() int {
	return len(cct.items)
}

func (cct *CommandConfigTable) Value(row, col int) interface{} {
	item := cct.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.CmdName
	case 2:
		return item.Desc
	case 3:
		return item.RenderTemplate
	case 4:
		return item.ParamTemplate
	case 5:
		return item.GroupType
	case 6:
		return item.CreateTime
	}
	panic("unexpected col")
}

func (cct *CommandConfigTable) Checked(row int) bool {
	return cct.items[row].IsChecked
}

func (cct *CommandConfigTable) SetChecked(row int, checked bool) error {
	cct.items[row].IsChecked = checked
	return nil
}
