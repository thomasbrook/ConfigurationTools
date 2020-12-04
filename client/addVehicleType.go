package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

// NewAddPanel 新建配置页面1
func NewAddPanel(parent walk.Container, vt *model.VehicleType, mainWin *TabMainWindow) (*VehicleTypeAddPage, error) {
	rand.Seed(time.Now().UnixNano())

	vtec := &VehicleTypeAddPage{
		searchModel:    &model.SearchModel{Items: []model.SearchItem{}},
		targetCanModel: &TargetCanTable{items: []*model.CanConfig{}},
		searchCanModel: LoadSearchCan(),
		editCanModel:   &model.CanConfig{},
		vehicleType:    vt,
		mainWin:        mainWin,
	}

	ef := &editForm{}

	if err := (Composite{
		AssignTo: &vtec.Composite,
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Composite{
						Layout:        VBox{MarginsZero: true},
						StretchFactor: 2,
						Children: []Widget{
							PushButton{
								Text: "批量搜索",
								OnClicked: func() {
									if cmd, err := ef.runExportDialog(vtec.mainWin); err != nil {
										log.Print(err)
									} else if cmd == walk.DlgCmdOK {
										vtec.mappingCan(ef.SearchStr)
									}
								},
							},
							TableView{
								AssignTo: &vtec.searchCanTv,
								//AlternatingRowBGColor: walk.RGB(239, 239, 239),
								CheckBoxes:       true,
								ColumnsOrderable: true,
								MultiSelection:   false,
								Columns: []TableViewColumn{
									{Title: "#", Width: 50},
									{Title: "中文名称"},
									{Title: "编号", Width: 50},
									{Title: "英文名称"},
									{Title: "单位"},
								},
								StyleCell: func(style *walk.CellStyle) {
									item := vtec.searchCanModel.items[style.Row()]

									if item.Checked {
										if style.Row()%2 == 0 {
											style.BackgroundColor = walk.RGB(159, 215, 255)
										} else {
											style.BackgroundColor = walk.RGB(143, 199, 239)
										}
									}

								},
								Model: vtec.searchCanModel,
							},
						},
					},
					Composite{
						Layout:        VBox{MarginsZero: true},
						StretchFactor: 6,
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									ComboBox{
										AssignTo:      &vtec.groupSel,
										Alignment:     AlignHNearVCenter,
										StretchFactor: 10,
										BindingMember: "Id",
										DisplayMember: "Name",
									},
									PushButton{
										Text:          "保存",
										StretchFactor: 1,
										OnClicked: func() {
											groupType := vtec.canGroupForAdd()
											if len(groupType) > 0 {
												g := groupType[vtec.groupSel.CurrentIndex()]
												groupInfoId := ""
												if g.Id == 0 {
													groupInfoId = vtec.vehicleType.GeneralId
												} else if g.Id == 1 {
													groupInfoId = vtec.vehicleType.CanId
												}
												err := dataSource.BatchInsertCan(groupInfoId, vtec.targetCanModel.items)
												if err != nil {
													panic(err)
												}
											}
										},
									},
								},
							},
							TableView{
								AssignTo: &vtec.targetCanTv,
								//AlternatingRowBGColor: walk.RGB(239, 239, 239),
								CheckBoxes:       true,
								ColumnsOrderable: true,
								MultiSelection:   true,
								Columns: []TableViewColumn{
									{Title: "#", Width: 50},
									{Title: "编号", Width: 50},
									{Title: "英文名称"},
									{Title: "中文名称"},
									{Title: "单位"},
									{Title: "数据类型"},
									{Title: "转换公式"},
									{Title: "数值范围"},
									{Title: "小数位"},
									{Title: "软报警项"},
									{Title: "可分析项"},
									{Title: "排序"},
								},
								StyleCell: func(style *walk.CellStyle) {
									item := vtec.targetCanModel.items[style.Row()]

									if item.Checked {
										if style.Row()%2 == 0 {
											style.BackgroundColor = walk.RGB(159, 215, 255)
										} else {
											style.BackgroundColor = walk.RGB(143, 199, 239)
										}
									}

								},
								Model: vtec.targetCanModel,
								OnSelectedIndexesChanged: func() {
									idxs := vtec.targetCanTv.SelectedIndexes()

									if len(idxs) > 0 {
										editCanModel := model.CanConfig{
											Index:        vtec.targetCanModel.items[idxs[0]].Index,
											OutfieldId:   vtec.targetCanModel.items[idxs[0]].OutfieldId,
											Unit:         vtec.targetCanModel.items[idxs[0]].Unit,
											OutfieldSn:   vtec.targetCanModel.items[idxs[0]].OutfieldSn,
											Chinesename:  vtec.targetCanModel.items[idxs[0]].Chinesename,
											Formula:      vtec.targetCanModel.items[idxs[0]].Formula,
											DataType:     vtec.targetCanModel.items[idxs[0]].DataType,
											FieldName:    vtec.targetCanModel.items[idxs[0]].FieldName,
											Decimals:     vtec.targetCanModel.items[idxs[0]].Decimals,
											DataMap:      vtec.targetCanModel.items[idxs[0]].DataMap,
											IsAlarm:      vtec.targetCanModel.items[idxs[0]].IsAlarm,
											IsAnalysable: vtec.targetCanModel.items[idxs[0]].IsAnalysable,
										}

										err := vtec.editCanDb.SetDataSource(&editCanModel)

										if err != nil {
											log.Print(fmt.Sprintf("%+v", err))
											return
										}

										ef.canDetail(&editCanModel)
									}
								},
							},
						},
					},
					Composite{
						Layout:        VBox{MarginsZero: true},
						StretchFactor: 3,
						Children: []Widget{
							ScrollView{
								Layout: Flow{
									MarginsZero: true,
									Alignment:   AlignHNearVNear,
								},
								Children: []Widget{
									Composite{
										Layout:    Grid{Columns: 2, Margins: Margins{5, 5, 5, 5}},
										Alignment: AlignHNearVNear,
										DataBinder: DataBinder{
											AssignTo:   &vtec.editCanDb,
											Name:       "editCanModel",
											DataSource: vtec.editCanModel,
										},
										Children: []Widget{
											Label{
												Text: "编号",
											},
											LineEdit{
												AssignTo: &ef.outfieldIdEdit,
												Text:     Bind("OutfieldId"),
												ReadOnly: true,
											},
											Label{
												Text: "英文名称",
											},
											LineEdit{
												AssignTo: &ef.fieldNameEdit,
												Text:     Bind("FieldName"),
												ReadOnly: true,
											},
											Label{
												Text: "中文名称",
											},
											LineEdit{
												AssignTo: &ef.cnEdit,
												Text:     Bind("Chinesename"),
											},
											Label{
												Text: "单位",
											},
											LineEdit{
												AssignTo: &ef.unitEdit,
												Text:     Bind("Unit"),
											},
											Label{
												Text: "数据类型",
											},
											RadioButtonGroupBox{
												Layout:   VBox{MarginsZero: true, SpacingZero: true},
												AssignTo: &ef.dataTypeRadioGroup,
												Buttons: []RadioButton{
													{AssignTo: &ef.dateValRb, Text: "日期时间", Value: 1},
													{AssignTo: &ef.enumValRb, Text: "枚举值", Value: 2},
													{AssignTo: &ef.numberValRb, Text: "数值", Value: 3},
													{AssignTo: &ef.enumTextValRb, Text: "枚举文本", Value: 5},
													{AssignTo: &ef.otherValRb, Text: "其他", Value: 4},
												},
											},
											Label{
												Text: "转换公式",
											},
											LineEdit{
												AssignTo: &ef.formulaEdit,
												Text:     Bind("Formula"),
											},
											Label{
												Text: "数值范围",
											},
											LineEdit{
												AssignTo: &ef.dataMapEdit,
												Text:     Bind("DataMap"),
											},
											Label{
												Text: "小数位",
											},
											LineEdit{
												AssignTo: &ef.decimalEdit,
												Text:     Bind("Decimals"),
											},
											Label{
												Text: "软报警项",
											},
											RadioButtonGroupBox{
												Layout: VBox{MarginsZero: true, SpacingZero: true},
												Buttons: []RadioButton{
													{AssignTo: &ef.alarmRb, Text: "报警项", Value: 1},
													{AssignTo: &ef.noAlarmRb, Text: "非报警项", Value: 0},
												},
											},
											Label{
												Text: "可分析项",
											},
											RadioButtonGroupBox{
												Layout: VBox{MarginsZero: true, SpacingZero: true},
												Buttons: []RadioButton{
													{AssignTo: &ef.analyableRb, Text: "可分析", Value: 1},
													{AssignTo: &ef.defaultAnalyzableRb, Text: "默认可分析", Value: 2},
													{AssignTo: &ef.noAnalyzableRb, Text: "不可分析", Value: 0},
												},
											},
											Label{
												Text: "排序号",
											},
											NumberEdit{
												AssignTo: &ef.outfieldIdSnEdit,
												Value:    Bind("OutfieldSn"),
											},
										},
									}},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text: "设置",
										OnClicked: func() {
											can, err := ef.set()
											if err != nil {
												panic(err)
											}
											vtec.refreshItem(can)
										},
									},
								},
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{
					MarginsZero: true,
					SpacingZero: true,
				},
				Children: []Widget{
					Label{
						AssignTo: &vtec.searchStatLbl,
						Text:     strconv.Itoa(len(vtec.searchCanModel.items)) + " 个搜索项",
					},
					Label{
						Text: ",  ",
					},
					Label{
						AssignTo: &vtec.targetStatLbl,
						Text:     "0 个选定项",
					},
					HSpacer{},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	group := vtec.canGroupForAdd()
	if len(group) > 0 {
		vtec.groupSel.SetModel(group)
		vtec.groupSel.SetCurrentIndex(0)
	}

	vtec.searchCanModel.parent = vtec
	vtec.targetCanModel.parent = vtec

	return vtec, nil
}

// VehicleTypeAddPage 页面对象
type VehicleTypeAddPage struct {
	*walk.Composite

	vehicleType *model.VehicleType
	mainWin     *TabMainWindow

	searchCanTv    *walk.TableView
	searchModel    *model.SearchModel
	searchCanModel *SearchCanTable

	groupSel       *walk.ComboBox
	targetCanTv    *walk.TableView
	targetCanModel *TargetCanTable

	searchStatLbl *walk.Label
	targetStatLbl *walk.Label

	editCanForm  *walk.Composite
	editCanModel *model.CanConfig
	editCanDb    *walk.DataBinder
}

func (vtec *VehicleTypeAddPage) canGroupForAdd() []*model.CanGroup {
	group := []*model.CanGroup{}

	if strings.TrimSpace(vtec.vehicleType.GeneralId) != "" {
		group = append(group, &model.CanGroup{
			Id:   0,
			Name: "常规信息",
		})
	}

	if strings.TrimSpace(vtec.vehicleType.CanId) != "" {
		group = append(group, &model.CanGroup{
			Id:   1,
			Name: "CAN信息",
		})
	}

	return group
}
func (vtec *VehicleTypeAddPage) mappingCan(searchStr string) {
	temp := strings.Split(searchStr, "\r\n")
	vtec.searchModel.Items = []model.SearchItem{}

	for _, item := range temp {
		if strings.TrimSpace(item) == "" {
			continue
		}

		vtec.searchModel.Items = append(vtec.searchModel.Items, model.SearchItem{item, ""})
	}

	result := dataSource.ListCanConfig()
	matchCan := []*model.CanConfig{}
	index := 0

	for i := 0; i < len(vtec.searchModel.Items); i++ {

		for j := 0; j < len(result.Props); j++ {
			idx := strings.Index(result.Props[j].Cn, vtec.searchModel.Items[i].Name)
			isEqual := false
			if idx == -1 {
				isEqual = strings.EqualFold(strings.TrimSpace(result.Props[j].Id), strings.TrimSpace(vtec.searchModel.Items[i].Name))
			}

			if idx != -1 || isEqual {
				index++
				matchCan = append(matchCan, &model.CanConfig{
					Index:       index,
					OutfieldId:  result.Props[j].Id,
					Unit:        result.Props[j].Unit,
					Chinesename: result.Props[j].Cn,
					FieldName:   result.Props[j].En})
			}
		}
	}

	if len(vtec.searchModel.Items) == 0 {

		idx := 0

		for j := 0; j < len(result.Props); j++ {
			idx++
			matchCan = append(matchCan, &model.CanConfig{
				Index:       idx,
				OutfieldId:  result.Props[j].Id,
				Unit:        result.Props[j].Unit,
				Chinesename: result.Props[j].Cn,
				FieldName:   result.Props[j].En})
		}
	}

	vtec.searchCanModel.items = matchCan
	vtec.searchCanTv.SetModel(vtec.searchCanModel)
	vtec.searchCanModel.ResetRows()
	vtec.searchStatLbl.SetText(strconv.Itoa(len(matchCan)) + "个搜索项")
}

func (vtec *VehicleTypeAddPage) addOrRemoveTargetItem(rowIndex int, isChecked bool) {
	item := vtec.searchCanModel.items[rowIndex]

	if isChecked {

		vtec.targetCanModel.items = append(vtec.targetCanModel.items, &model.CanConfig{
			Index:       len(vtec.targetCanModel.items) + 1,
			OutfieldId:  item.OutfieldId,
			Unit:        item.Unit,
			Chinesename: item.Chinesename,
			FieldName:   item.FieldName,
			Checked:     true,
		})
	} else {
		targetCan := []*model.CanConfig{}
		idx := 0
		for i := 0; i < len(vtec.targetCanModel.items); i++ {
			if vtec.targetCanModel.items[i].OutfieldId != item.OutfieldId {
				idx = idx + 1
				vtec.targetCanModel.items[i].Index = idx
				vtec.targetCanModel.items[i].Checked = true
				targetCan = append(targetCan, vtec.targetCanModel.items[i])
			}
		}
		vtec.targetCanModel.items = targetCan
	}
	vtec.targetCanTv.SetModel(vtec.targetCanModel)
	vtec.targetCanModel.ResetRows()
	vtec.targetStatLbl.SetText(strconv.Itoa(len(vtec.targetCanModel.items)) + "个选定项")
}

func (vtec *VehicleTypeAddPage) removeTargetItem(rowIndex int, isChecked bool) {
	item := vtec.targetCanModel.items[rowIndex]
	if !isChecked {
		// 删除勾选项
		targetCan := []*model.CanConfig{}
		idx := 0
		for i := 0; i < len(vtec.targetCanModel.items); i++ {
			if vtec.targetCanModel.items[i].OutfieldId != item.OutfieldId {
				idx = idx + 1
				vtec.targetCanModel.items[i].Index = idx
				targetCan = append(targetCan, vtec.targetCanModel.items[i])
			}
		}
		vtec.targetCanModel.items = targetCan
		vtec.targetCanTv.SetModel(vtec.targetCanModel)
		vtec.targetCanModel.ResetRows()
		vtec.targetStatLbl.SetText(strconv.Itoa(len(vtec.targetCanModel.items)) + "个选定项")

		// 重置选定状态
		for i := 0; i < len(vtec.searchCanModel.items); i++ {
			if vtec.searchCanModel.items[i].OutfieldId == item.OutfieldId {
				vtec.searchCanModel.items[i].Checked = false
			}
		}
		vtec.searchCanModel.ResetRows()
	}
}

func (vtec *VehicleTypeAddPage) refreshItem(can *model.CanConfig) {
	for i := 0; i < len(vtec.targetCanModel.items); i++ {
		if vtec.targetCanModel.items[i].OutfieldId == can.OutfieldId {
			vtec.targetCanModel.items[i].Chinesename = can.Chinesename
			vtec.targetCanModel.items[i].Unit = can.Unit
			vtec.targetCanModel.items[i].DataType = can.DataType
			vtec.targetCanModel.items[i].Formula = can.Formula
			vtec.targetCanModel.items[i].DataMap = can.DataMap
			vtec.targetCanModel.items[i].Decimals = can.Decimals
			vtec.targetCanModel.items[i].IsAlarm = can.IsAlarm
			vtec.targetCanModel.items[i].IsAnalysable = can.IsAnalysable
			vtec.targetCanModel.items[i].OutfieldSn = can.OutfieldSn

			break
		}
	}
	vtec.targetCanModel.ResetRows()
}

// 过滤can
type SearchCanTable struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*model.CanConfig

	parent *VehicleTypeAddPage
}

func LoadSearchCan() *SearchCanTable {
	m := new(SearchCanTable)

	result := dataSource.ListCanConfig()
	matchCan := []*model.CanConfig{}
	index := 0

	for j := 0; j < len(result.Props); j++ {
		index++
		matchCan = append(matchCan, &model.CanConfig{
			Index:       index,
			OutfieldId:  result.Props[j].Id,
			Unit:        result.Props[j].Unit,
			Chinesename: result.Props[j].Cn,
			FieldName:   result.Props[j].En})
	}

	m.items = matchCan

	m.ResetRows()
	return m
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *SearchCanTable) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *SearchCanTable) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Chinesename
	case 2:
		return item.OutfieldId
	case 3:
		return item.FieldName
	case 4:
		return item.Unit
	}

	panic("unexpected col")
}

// Checked Called by the TableView to retrieve if a given row is checked.
func (m *SearchCanTable) Checked(row int) bool {
	return m.items[row].Checked
}

// SetChecked Called by the TableView when the user toggled the check box of a given row.
func (m *SearchCanTable) SetChecked(row int, checked bool) error {
	m.items[row].Checked = checked
	m.parent.addOrRemoveTargetItem(row, checked)
	return nil
}

// Sort Called by the TableView to sort the model.
func (m *SearchCanTable) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.Chinesename < b.Chinesename)

		case 2:
			return c(a.OutfieldId < b.OutfieldId)

		case 3:
			return c(a.FieldName < b.FieldName)

		case 4:
			return c(a.Unit < a.Unit)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

// ResetRows 排序
func (m *SearchCanTable) ResetRows() {

	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}

// TargetCanTable 目标can
type TargetCanTable struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*model.CanConfig

	parent *VehicleTypeAddPage
}

// LoadTargetCan 目标Table
func LoadTargetCan() *TargetCanTable {
	m := new(TargetCanTable)
	m.ResetRows()
	return m
}

// NewTargetCanModel
func NewTargetCanModel() *TargetCanTable {
	m := new(TargetCanTable)
	m.ResetRows()
	return m
}

// RowCount Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *TargetCanTable) RowCount() int {
	return len(m.items)
}

// Value Called by the TableView when it needs the text to display for a given cell.
func (m *TargetCanTable) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index
	case 1:
		return item.OutfieldId
	case 2:
		return item.FieldName
	case 3:
		return item.Chinesename
	case 4:
		return item.Unit
	case 5:
		return item.DataType
	case 6:
		return item.Formula
	case 7:
		return item.DataMap
	case 8:
		return item.Decimals
	case 9:
		return item.IsAlarm
	case 10:
		return item.IsAnalysable
	case 11:
		return item.OutfieldSn
	}

	panic("unexpected col")
}

// Checked Called by the TableView to retrieve if a given row is checked.
func (m *TargetCanTable) Checked(row int) bool {
	return m.items[row].Checked
}

// SetChecked Called by the TableView when the user toggled the check box of a given row.
func (m *TargetCanTable) SetChecked(row int, checked bool) error {
	m.items[row].Checked = checked
	m.parent.removeTargetItem(row, checked)
	return nil
}

// Sort Called by the TableView to sort the model.
func (m *TargetCanTable) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)

		case 1:
			return c(a.OutfieldId < b.OutfieldId)

		case 2:
			return c(a.FieldName < b.FieldName)

		case 3:
			return c(a.Chinesename < b.Chinesename)

		case 4:
			return c(a.Unit < b.Unit)

		case 5:
			return c(a.DataType < b.DataType)

		case 6:
			return c(a.Formula < b.Formula)

		case 7:
			return c(a.DataMap < b.DataMap)

		case 8:
			return c(a.Decimals < b.Decimals)

		case 9:
			return c(a.IsAlarm < a.IsAlarm)

		case 10:
			return c(a.IsAnalysable < a.IsAnalysable)

		case 11:
			return c(a.OutfieldSn < b.OutfieldSn)
		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

// ResetRows
func (m *TargetCanTable) ResetRows() {
	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}
