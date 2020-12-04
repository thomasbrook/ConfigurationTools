package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"ConfigurationTools/utils"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"strconv"
	"strings"
)

// ImportCanFromExistingPanel 从已有车型导入CAN
func ImportCanFromExistingPanel(parent walk.Container, vt *model.VehicleTypeStats, mainWin *TabMainWindow) (*CanCopyPage, error) {
	ccp := &CanCopyPage{
		refercanListTvModel: new(CanListTableModel),
		mainWin:             mainWin,
		targetVehicleType:   vt,
	}

	data, err := (&dataSource.VehicleTypeEntity{TypeId: vt.TypeId}).GetVehicleType()
	if err != nil {
		return ccp, err
	}

	fieldExistIcon, _ := walk.Resources.Icon("/img/warn.ico")

	if err := (Composite{
		AssignTo: &ccp.Composite,
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Composite{
						Layout:        VBox{MarginsZero: true, Margins: Margins{Right: 3}},
						StretchFactor: 2,
						Children: []Widget{
							Composite{
								Layout:  HBox{MarginsZero: true},
								MaxSize: Size{Height: 21},
								Children: []Widget{
									LineEdit{
										AssignTo: &ccp.searchStrLe,
										Text:     "",
										OnKeyDown: func(key walk.Key) {
											if key == walk.KeyReturn {
												m := ccp.NewVehicleTypeModel()
												ccp.vehicleTypeTv.SetModel(m)
											}
										},
									},
									PushButton{
										AssignTo: &ccp.queryPb,
										//Image:    "/img/search.ico",
										Text: "查询",
										OnClicked: func() {
											ccp.queryPb.SetEnabled(false)
											defer ccp.queryPb.SetEnabled(true)

											m := ccp.NewVehicleTypeModel()
											ccp.vehicleTypeTv.SetModel(m)
										},
									},
								},
							},
							TableView{
								AssignTo:         &ccp.vehicleTypeTv,
								AlternatingRowBG: true,
								//AlternatingRowBGColor: walk.RGB(239, 239, 239),
								ColumnsOrderable: true,
								//MultiSelection:        true,
								Columns: []TableViewColumn{
									{Name: "TypeName", Title: "车型"},
									{Name: "OrgName", Title: "厂商"},
								},
								OnItemActivated: func() {

									if !ccp.referCanGroupCb.Enabled() {
										return
									}

									ccp.vehicleTypeTv.SetEnabled(false)
									defer ccp.vehicleTypeTv.SetEnabled(true)
									currentIdx := ccp.vehicleTypeTv.CurrentIndex()

									vtm := ccp.vehicleTypeTv.Model().(*VehicleTypeModel)
									err := ccp.referVehicleTypeLbl.SetText(vtm.items[currentIdx].TypeName)
									if err != nil {
										return
									}

									canGroup := []model.GroupStats{}
									for _, group := range vtm.items[currentIdx].Group {
										canGroup = append(canGroup, *group)
									}

									if len(canGroup) > 0 {
										ccp.referCanGroupCb.SetModel(canGroup)
										ccp.referCanGroupCb.SetCurrentIndex(0)
									}
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true, Margins: Margins{Top: 3}},
								Children: []Widget{
									Label{
										AssignTo: &ccp.statTipsLbl,
										Font:     Font{PointSize: 10},
									},
									HSpacer{},
								},
							},
						},
					},
					Composite{
						Layout:        VBox{MarginsZero: true, Margins: Margins{Left: 3}},
						StretchFactor: 8,
						Children: []Widget{
							Composite{
								Layout:  HBox{MarginsZero: true, Margins: Margins{Top: 2, Left: 5}},
								MaxSize: Size{Height: 21},
								Children: []Widget{
									CheckBox{
										AssignTo: &ccp.allCheckedChk,
										Text:     "全选",
										OnClicked: func() {
											ccp.allCheckedChk.SetEnabled(false)
											defer ccp.allCheckedChk.SetEnabled(true)

											for _, v := range ccp.refercanListTvModel.items {
												v.Checked = ccp.allCheckedChk.Checked()
											}
											ccp.refercanListTvModel.PublishRowsReset()
										},
									},
									HSpacer{},
									Label{
										Text: "源车型：",
									},
									Label{
										AssignTo: &ccp.referVehicleTypeLbl,
										Text:     "未选择",
										MinSize:  Size{Width: 145},
									},
									ComboBox{
										AssignTo:      &ccp.referCanGroupCb,
										BindingMember: "Id",
										DisplayMember: "Name",
										MinSize:       Size{Width: 160},
										OnCurrentIndexChanged: func() {
											ccp.referCanGroupCb.SetEnabled(false)
											defer ccp.referCanGroupCb.SetEnabled(true)

											selectedIdx := ccp.referCanGroupCb.CurrentIndex()
											if selectedIdx == -1 {
												return
											}

											m := ccp.referCanGroupCb.Model().([]model.GroupStats)
											can, err := (&dataSource.VehicleTypeEntity{}).ListCan(m[selectedIdx].Id)
											if err != nil {
												ccp.referCanGroupCb.SetEnabled(true)
												mynotify.Error("查询失败：" + err.Error())
												return
											}

											ccp.refercanListTvModel.items = []*model.CanDetailTableAdapter{}
											for i, v := range can {
												temp := model.CanDetailTableAdapter{
													Index:     i + 1,
													Id:        v.Id,
													Key:       v.OutfieldId,
													Unit:      v.Unit,
													Sort:      v.OutfieldSn,
													GroupId:   v.GroupInfoId,
													Alias:     v.Chinesename,
													Formula:   v.Formula,
													FieldName: v.FieldName,
													Prec:      v.Decimals,
													DataScope: v.DataMap,
													Note:      v.Note,
												}

												_, isExist := ccp.targetCanMap[v.OutfieldId]
												if isExist {
													temp.Note = "已存在"
												}

												// 数据类型
												_dataType, err := strconv.Atoi(v.DataType)
												if err != nil {
													_dataType = 4
												}
												temp.DataType = utils.ToDataType(_dataType)

												// 是否可报警
												_alarm, err := strconv.Atoi(v.IsAlarm)
												if err != nil {
													_alarm = 0
												}
												temp.IsAlarm = utils.ToAlarm(_alarm)

												// 是否可分析
												temp.IsAnalysable = utils.ToAlayly(v.IsAnalysable)
												ccp.refercanListTvModel.items = append(ccp.refercanListTvModel.items, &temp)
											}
											ccp.refercanListTvModel.PublishRowsReset()
											ccp.referCanGroupCb.SetEnabled(true)
										},
									},
								},
							},
							TableView{
								AssignTo:         &ccp.referCanListTv,
								AlternatingRowBG: true,
								//AlternatingRowBGColor: walk.RGB(239, 239, 239),
								ColumnsOrderable: true,
								CheckBoxes:       true,
								MultiSelection:   true,
								Model:            ccp.refercanListTvModel,
								Columns: []TableViewColumn{
									{Title: "#", Frozen: true, Alignment: AlignCenter, Width: 60},
									{Title: "编号", Width: 60, Alignment: AlignFar},
									{Title: "中文名", Alignment: AlignCenter, Width: 120},
									{Title: "别名", Alignment: AlignCenter, Width: 120},
									{Title: "单位", Alignment: AlignCenter, Width: 50},
									{Title: "数据类型", Alignment: AlignCenter, Width: 80},
									{Title: "转换公式", Alignment: AlignCenter, Width: 80},
									{Title: "数值范围", Alignment: AlignCenter, Width: 160},
									{Title: "小数位", Alignment: AlignCenter, Width: 50},
									{Title: "软报警项", Alignment: AlignCenter, Width: 75},
									{Title: "可分析项", Alignment: AlignCenter, Width: 85},
									{Title: "排序", Alignment: AlignCenter, Width: 50},
								},
								OnSelectedIndexesChanged: func() {
									log.Println(fmt.Printf("SelectedIndexes: %v\n", ccp.referCanListTv.SelectedIndexes()))
								},
								StyleCell: func(style *walk.CellStyle) {
									item := ccp.refercanListTvModel.items[style.Row()]
									switch style.Col() {
									case 1:
										idx := strings.Index(item.Note, "已存在")
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
									ImageView{Image: "/img/warn.ico"},
									Label{Text: "：当前车系下，已存在该字段，可替换/跳过"},
									HSpacer{},
									Label{
										Text: "将勾选的行导入到 -->",
										//Font:      Font{PointSize: 8,},
										TextColor: walk.RGB(0xDC, 0x14, 0x3C),
									},
									Label{
										Text: fmt.Sprintf("目标车型：%s", vt.TypeName),
										Font: Font{PointSize: 10},
									},
									ComboBox{
										AssignTo:      &ccp.targetCanGroupCb,
										BindingMember: "Id",
										DisplayMember: "Name",
										MinSize:       Size{Width: 160},
										Model:         data.CanGroup,
										OnCurrentIndexChanged: func() {

											selectIdx := ccp.targetCanGroupCb.CurrentIndex()
											if selectIdx == -1 {
												return
											}

											cangroup := ccp.targetCanGroupCb.Model().([]*dataSource.CanGroupEntity)
											can, err := (&dataSource.VehicleTypeEntity{}).ListCan(cangroup[selectIdx].Id)
											if err != nil {
												mynotify.Error("查询失败：" + err.Error())
												return
											}

											ccp.targetCanMap = make(map[string]*model.CanDetail)
											for _, v := range can {
												ccp.targetCanMap[v.OutfieldId] = v
											}

											for _, v := range ccp.refercanListTvModel.items {
												_, isExist := ccp.targetCanMap[v.Key]
												if isExist {
													v.Note = "已存在"
												} else {
													v.Note = ""
												}
											}
											ccp.refercanListTvModel.PublishRowsReset()
										},
									},
									PushButton{
										AssignTo: &ccp.importCanPb,
										Image:    "/img/import.ico",
										Text:     "导入",
										OnClicked: func() {
											ccp.importCanPb.SetEnabled(false)
											defer ccp.importCanPb.SetEnabled(true)
											isSkip := false

											// 当前勾选的字段列表
											var selectField []model.CanDetail

											for _, v := range ccp.refercanListTvModel.items {
												if !v.Checked {
													continue
												}

												selectField = append(selectField, model.CanDetail{
													Id:         v.Id,
													OutfieldId: v.Key,
												})

												if isSkip {
													continue
												}

												if _, isExist := ccp.targetCanMap[v.Key]; isExist {
													isSkip = true
												}
											}

											if len(selectField) == 0 {
												ccp.importCanPb.SetEnabled(true)
												walk.MsgBox(ccp.mainWin, "", "请勾选导入项", walk.MsgBoxIconWarning)
												return
											}

											// 目标can分组ID
											var targetCanGroup = ccp.targetCanGroupCb.Model().([]*dataSource.CanGroupEntity)
											idx := ccp.targetCanGroupCb.CurrentIndex()
											if idx == -1 {
												ccp.importCanPb.SetEnabled(true)
												walk.MsgBox(ccp.mainWin, "", "请选择目标分组", walk.MsgBoxIconWarning)
												return
											}
											var groupId = targetCanGroup[idx].Id

											// 如果当前勾选字段，都不存在于当前分组字段内，直接执行插入操作
											var mode = 2
											if !isSkip {
												err := (&dataSource.VehicleTypeEntity{}).SyncCanFromExisting(groupId, selectField, mode)
												ccp.importCanPb.SetEnabled(true)
												if err != nil {
													walk.MsgBox(ccp.mainWin, "执行失败", err.Error(), walk.MsgBoxIconError)
													return
												}

												walk.MsgBox(ccp.mainWin, "", "执行成功", walk.MsgBoxIconInformation)
												return
											}

											// 否则，提醒用户进行分类处理
											cmd, err := RunConfirmDialog(ccp.mainWin)
											if err != nil {
												ccp.importCanPb.SetEnabled(true)
												mynotify.Error("确认窗口初始化失败：" + err.Error())
												return
											}

											if cmd == walk.DlgCmdCancel || cmd == walk.DlgCmdNone {
												ccp.importCanPb.SetEnabled(true)
												return
											}

											if cmd == walk.DlgCmdOK {
												mode = 1
											} else if cmd == walk.DlgCmdIgnore {
												mode = 2
											}

											err = (&dataSource.VehicleTypeEntity{}).SyncCanFromExisting(groupId, selectField, 1)
											if err != nil {
												ccp.importCanPb.SetEnabled(true)
												walk.MsgBox(ccp.mainWin, "执行失败", err.Error(), walk.MsgBoxIconError)
												return
											}

											ccp.importCanPb.SetEnabled(true)
											walk.MsgBox(ccp.mainWin, "", "执行成功", walk.MsgBoxIconInformation)
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	m := ccp.NewVehicleTypeModel()
	ccp.vehicleTypeTv.SetModel(m)

	return ccp, nil
}

type CanCopyPage struct {
	*walk.Composite

	mainWin           *TabMainWindow
	targetVehicleType *model.VehicleTypeStats

	// 左侧搜索框、车型列表
	searchStrLe   *walk.LineEdit
	queryPb       *walk.PushButton
	vehicleTypeTv *walk.TableView
	statTipsLbl   *walk.Label

	// 被选定的车型信息
	referVehicleTypeLbl *walk.Label
	referCanGroupCb     *walk.ComboBox

	// 全选按钮
	allCheckedChk *walk.CheckBox

	// 选定的分组CAN列表
	referCanListTv      *walk.TableView
	refercanListTvModel *CanListTableModel

	// 目标can分组
	targetCanGroupCb *walk.ComboBox

	// 目标分组CAN哈希
	targetCanMap map[string]*model.CanDetail
	importCanPb  *walk.PushButton
}

func (ccp *CanCopyPage) NewVehicleTypeModel() *VehicleTypeModel {
	searchStr := ccp.searchStrLe.Text()
	vtype, err := (&dataSource.OrganizationEntity{SearchKey: searchStr}).ListVehicleType()
	if err != nil {
		mynotify.Error("查询车型列表失败：" + err.Error())
		return nil
	}

	m := &VehicleTypeModel{}
	for _, v := range vtype {
		if ccp.targetVehicleType != nil && ccp.targetVehicleType.TypeId != v.TypeId {
			m.items = append(m.items, v)
		}
	}

	ccp.statTipsLbl.SetText(fmt.Sprintf("共 %d 项", len(m.items)))
	return m
}

type VehicleTypeModel struct {
	walk.SortedReflectTableModelBase
	items []*model.VehicleTypeStats
}

func (m *VehicleTypeModel) Items() interface{} {
	return m.items
}

type CanListTableModel struct {
	walk.TableModelBase
	items []*model.CanDetailTableAdapter
}

func (m *CanListTableModel) RowCount() int {
	return len(m.items)
}

func (m *CanListTableModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Key
	case 2:
		return item.FieldName
	case 3:
		return item.Alias
	case 4:
		return item.Unit
	case 5:
		return item.DataType
	case 6:
		return item.Formula
	case 7:
		return item.DataScope
	case 8:
		return item.Prec
	case 9:
		return item.IsAlarm
	case 10:
		return item.IsAnalysable
	case 11:
		return item.Sort
	case 12:
		return item.Note
	}
	return ""
}

func (m *CanListTableModel) Checked(row int) bool {
	return m.items[row].Checked
}

func (m *CanListTableModel) SetChecked(row int, checked bool) error {
	m.items[row].Checked = checked
	return nil
}

func RunConfirmDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		Title:         "提示",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,

		MinSize: Size{Width: 300},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text:      "部分字段已存在，是否全部替换？",
				TextColor: walk.RGB(0xDC, 0x14, 0x3C),
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{AssignTo: &acceptPB,
						Text: "全部替换",
						OnClicked: func() {
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "全部跳过",
						OnClicked: func() {
							dlg.Close(walk.DlgCmdIgnore)
						},
					},
					PushButton{
						Text: "取消",
						OnClicked: func() {
							dlg.Close(walk.DlgCmdCancel)
						},
					},
				},
			},
		},
	}.Run(owner)
}
