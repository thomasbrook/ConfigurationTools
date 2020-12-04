package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strconv"
	"strings"
)

func NewImportCanFromClipboardPanel(parent walk.Container, vt *model.VehicleTypeStats, mainWin *TabMainWindow) (*ImportCanFromClipboardPage, error) {
	cfc := &ImportCanFromClipboardPage{
		mainWin:           mainWin,
		targetVehicleType: vt,
	}

	_, err := (&dataSource.VehicleTypeEntity{TypeId: vt.TypeId}).GetVehicleType()
	if err != nil {
		return nil, err
	}

	noGroupIcon, _ := walk.Resources.Icon("/img/delete.ico")
	fieldExistIcon, _ := walk.Resources.Icon("/img/warn.ico")

	if err := (Composite{
		AssignTo: &cfc.Composite,
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{MarginsZero: true, Margins: Margins{Right: 3}},
				Children: []Widget{
					ImageView{
						Image:       "/img/paste.ico",
						ToolTipText: "从剪贴板黏贴",
						MinSize:     Size{Width: 22},
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							err := cfc.parserClipboardText()
							if err != nil {
								walk.MsgBox(cfc.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
								return
							}
						},
					},
					ImageView{
						Image:       "/img/clear.ico",
						ToolTipText: "清空",
						MinSize:     Size{Width: 22},
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							cfc.previewCanAdapter = nil
							cfc.previewCanTbl.SetModel(nil)
						},
					},
					HSpacer{},
					PushButton{
						Image: "/img/import.ico",
						Text:  "导入",
						OnClicked: func() {

							if len(cfc.previewCanAdapter) == 0 {
								walk.MsgBox(cfc.mainWin, "", "数据不存在", walk.MsgBoxIconInformation)
								return
							}

							vt, err := (&dataSource.VehicleTypeEntity{TypeId: cfc.targetVehicleType.TypeId}).GetVehicleType()
							if err != nil {
								walk.MsgBox(cfc.mainWin, "", err.Error(), walk.MsgBoxIconError)
								return
							}

							var existedGroupMap = make(map[string]string)
							var groupId = []string{}
							for _, v := range vt.CanGroup {
								existedGroupMap[strconv.Itoa(v.Code)] = v.Id
								groupId = append(groupId, v.Id)
							}

							// 过滤掉不存在的分组，以及判断是否需要跳过已存在的字段
							var isSkip bool
							var data []*model.CanDetailWithGroup
							for _, v := range cfc.previewCanAdapter {
								if _, isExist := existedGroupMap[v.Remark]; isExist {

									groupId, isExist := existedGroupMap[v.Remark]
									if isExist {
										v.GroupInfoId = groupId
									}

									data = append(data, v)

									if !isSkip {
										idx := strings.Index(v.Note, "字段已存在")
										if idx >= 0 {
											isSkip = true
										}
									}
								}
							}

							var mode = 2
							if !isSkip {
								err = (&dataSource.VehicleTypeEntity{}).SyncCanFromCsv(groupId, data, mode)
								if err != nil {
									walk.MsgBox(cfc.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
									return
								}
								mynotify.Info("导入成功")
								return
							}

							cmd, err := RunConfirmDialog(cfc.mainWin)
							if err != nil {
								mynotify.Error("确认窗口初始化失败：" + err.Error())
								return
							}

							if cmd == walk.DlgCmdCancel || cmd == walk.DlgCmdNone {
								return
							}

							if cmd == walk.DlgCmdOK {
								mode = 1
							} else if cmd == walk.DlgCmdIgnore {
								mode = 2
							}

							err = (&dataSource.VehicleTypeEntity{}).SyncCanFromCsv(groupId, data, mode)
							if err != nil {
								walk.MsgBox(cfc.mainWin, "执行失败", err.Error(), walk.MsgBoxIconError)
								return
							}
							mynotify.Info("导入成功")
						},
					},
				},
			},
			TableView{
				AssignTo:         &cfc.previewCanTbl,
				AlternatingRowBG: true,
				//AlternatingRowBGColor: walk.RGB(239, 239, 239),
				ColumnsOrderable: true,
				Columns: []TableViewColumn{
					{Name: "Index", Title: "#", Frozen: true, Width: 60, Alignment: AlignCenter},
					{Name: "GroupName", Title: "分组名", Width: 120, Alignment: AlignCenter},
					{Name: "Sort", Title: "分组排序", Alignment: AlignCenter, Width: 60},
					{Name: "Remark", Title: "分组标识", Alignment: AlignCenter, Width: 60},
					{Name: "Chinesename", Title: "别名", Alignment: AlignCenter, Width: 140},
					{Name: "FieldName", Title: "字段名", Alignment: AlignCenter, Width: 140},
					{Name: "OutfieldId", Title: "CAN编码", Width: 60, Alignment: AlignFar},
					{Name: "Unit", Title: "单位", Alignment: AlignCenter, Width: 60},
					{Name: "DataType", Title: "数据类型", Alignment: AlignCenter, Width: 100, FormatFunc: func(value interface{}) string {
						switch value {
						case "1":
							return "日期时间"
						case "2":
							return "数字枚举"
						case "3":
							return "数据"
						case "4":
							return "其他"
						case "5":
							return "文本枚举"
						case "6":
							return "文本多枚举值"
						case "7":
							return "多字段组合多枚举值"
						default:
							return ""
						}
					}},
					{Name: "Formula", Title: "转换公式", Alignment: AlignCenter, Width: 100},
					{Name: "Decimals", Title: "小数位", Alignment: AlignCenter, Width: 50},
					{Name: "DataMap", Title: "值域", Alignment: AlignCenter, Width: 160},
					{Name: "IsAlarm", Title: "是否软报警", Alignment: AlignCenter, Width: 75, FormatFunc: func(value interface{}) string {
						switch value {
						case "0":
							return ""
						case "1":
							return "√"
						default:
							return ""
						}
					}},
					{Name: "IsAnalysable", Title: "是否可分析", Alignment: AlignCenter, Width: 75, FormatFunc: func(value interface{}) string {
						switch value {
						case 0:
							return ""
						case 1:
							return "√"
						case 2:
							return "√√"
						default:
							return ""
						}
					}},
					{Name: "IsDelete", Title: "是否删除", Alignment: AlignCenter, Width: 75, FormatFunc: func(value interface{}) string {
						switch value {
						case 0:
							return ""
						case 1:
							return "√"
						default:
							return ""
						}
					}},
					{Name: "OutfieldSn", Title: "CAN排序", Alignment: AlignCenter, Width: 75},
				},
				StyleCell: func(style *walk.CellStyle) {
					m := cfc.previewCanTbl.Model().([]*model.CanDetailWithGroup)
					item := m[style.Row()]
					switch style.Col() {
					case 1:
						idx := strings.Index(item.Note, "分组不存在")
						if idx != -1 {
							//style.TextColor = walk.RGB(255, 0, 0)
							style.Image = noGroupIcon
						}
					case 6:
						idx := strings.Index(item.Note, "字段已存在")
						if idx != -1 {
							//style.TextColor = walk.RGB(255, 0, 0)
							style.Image = fieldExistIcon
						}
					}
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true, Margins: Margins{Top: 3}},
				Children: []Widget{
					Label{
						AssignTo: &cfc.statusLbl,
						Font:     Font{PointSize: 10},
					},
					HSpacer{},
					ImageView{Image: "/img/delete.ico"},
					Label{Text: "：当前车系下，无此分组，CAN字段将不会导入"},
					ImageView{Image: "/img/warn.ico"},
					Label{Text: "：当前车系下，已存在该字段，可替换/跳过"},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	return cfc, nil
}

type ImportCanFromClipboardPage struct {
	*walk.Composite

	mainWin           *TabMainWindow
	targetVehicleType *model.VehicleTypeStats

	previewCanTbl     *walk.TableView
	previewCanAdapter []*model.CanDetailWithGroup

	statusLbl *walk.Label
}

func (cfc *ImportCanFromClipboardPage) parserClipboardText() error {
	txt, err := walk.Clipboard().Text()
	if err != nil {
		return err
	}

	var clipboardCan []*model.CanDetailWithGroup
	err = json.Unmarshal([]byte(txt), &clipboardCan)
	if err != nil {
		return err
	}

	// 查询当前机型已有can配置，并哈希存储
	existedCanMap := make(map[string]bool)
	existedCan, err := (&dataSource.VehicleTypeEntity{}).ListAllCan(cfc.targetVehicleType.TypeId)
	if err != nil {
		return err
	}

	for i := 0; i < len(existedCan); i++ {
		key := fmt.Sprintf("%s_%s", existedCan[i].Remark, existedCan[i].OutfieldId)
		existedCanMap[key] = true
	}

	// 查询当前车系分组
	existedGroupMap := make(map[string]bool)
	vehicleType, err := (&dataSource.VehicleTypeEntity{TypeId: cfc.targetVehicleType.TypeId}).GetVehicleType()
	if err != nil {
		return err
	}

	for i := 0; i < len(vehicleType.CanGroup); i++ {
		existedGroupMap[strconv.Itoa(vehicleType.CanGroup[i].Code)] = true
	}

	cfc.previewCanAdapter = nil

	// 遍历剪贴板CAN内容，转换处理，并填充TableView
	for idx, item := range clipboardCan {

		item.Index = idx + 1

		item.Note = ""

		_, isExist := existedGroupMap[item.Remark]
		if !isExist {
			item.Note = "分组不存在"
		} else {
			_, isExist := existedCanMap[fmt.Sprintf("%s_%s", item.Remark, item.OutfieldId)]
			if isExist {
				item.Note = "字段已存在"
			}
		}
	}

	cfc.previewCanAdapter = clipboardCan

	cfc.previewCanTbl.SetModel(cfc.previewCanAdapter)
	cfc.statusLbl.SetText(fmt.Sprintf("共 %d 项", len(cfc.previewCanAdapter)))
	if len(cfc.previewCanAdapter) == 0 {
		walk.MsgBox(cfc.mainWin, "提示", "导入文件无数据", walk.MsgBoxIconWarning)
	}
	return nil
}
