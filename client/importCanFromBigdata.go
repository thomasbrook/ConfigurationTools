package client

import (
	"ConfigurationTools/configurationManager"
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"ConfigurationTools/utils"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

// ImportCanFromBigdataPanel 从大数据XML导入CAN
func ImportCanFromBigdataPanel(parent walk.Container, vt *model.VehicleTypeStats, mainWin *TabMainWindow) (*CanAddPage, error) {
	rand.Seed(time.Now().UnixNano())

	var bigdataCan *SearchCanTable
	var url string
	if strings.TrimSpace(configurationManager.CANConfig.TestUrl) != "" {
		url = configurationManager.CANConfig.TestUrl
	} else if strings.TrimSpace(configurationManager.CANConfig.ProUrl) != "" {
		url = configurationManager.CANConfig.ProUrl
	} else {
		url = "http://192.168.11.8/config.xml"
	}

	can, err := LoadSearchCan(url)
	if err != nil {
		walk.MsgBox(mainWin, "", err.Error(), walk.MsgBoxIconError)
		bigdataCan = new(SearchCanTable)
	} else {
		bigdataCan = can
	}

	vtec := &CanAddPage{
		searchModel:    &model.SearchModel{Items: []model.SearchItem{}},
		targetCanModel: &TargetCanTable{items: []*model.CanDetail{}},
		searchCanModel: bigdataCan,
		vehicleType:    vt,
		mainWin:        mainWin,
	}

	ef := &editForm{}

	fieldExistIcon, _ := walk.Resources.Icon("/img/warn.ico")

	if err := (Composite{
		AssignTo: &vtec.Composite,
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Composite{
						Layout:        VBox{MarginsZero: true, Margins: Margins{Right: 3}},
						StretchFactor: 2,
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									ComboBox{
										AssignTo:           &vtec.urlCb,
										BindingMember:      "Value",
										DisplayMember:      "Key",
										CurrentIndex:       0,
										Model:              utils.KnownConfigUrl(),
										Alignment:          AlignHNearVCenter,
										AlwaysConsumeSpace: true,
										OnCurrentIndexChanged: func() {

										},
									},
									ImageView{
										Image:       "/img/search1.ico",
										ToolTipText: "搜索",
										Alignment:   AlignHFarVCenter,
										OnMouseDown: func(x, y int, button walk.MouseButton) {

											urls := vtec.urlCb.Model().([]*utils.KeyValuePair2)
											ef.ConfigUrl = urls[vtec.urlCb.CurrentIndex()].Value

											if cmd, err := ef.runExportDialog(vtec.mainWin); err != nil {
												mynotify.Error("CAN关键词搜索：" + err.Error())
											} else if cmd == walk.DlgCmdOK {
												vtec.mappingCan(ef.SearchStr, ef.ConfigUrl)
											}
										},
									},
								}},
							TableView{
								AssignTo:         &vtec.searchCanTv,
								AlternatingRowBG: true,
								//AlternatingRowBGColor: walk.RGB(239, 239, 239),
								CheckBoxes:       true,
								ColumnsOrderable: true,
								MultiSelection:   false,
								Columns: []TableViewColumn{
									{Title: "#", Width: 50, Frozen: true},
									{Title: "中文名"},
									{Title: "编号", Width: 50},
									{Title: "英文名"},
									{Title: "单位", Width: 50},
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
						Layout:        VBox{MarginsZero: true, Margins: Margins{Left: 3}},
						StretchFactor: 8,
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{
										Text: fmt.Sprintf("%s - %s", vt.OrgName, vt.TypeName),
									},
									HSpacer{},
									ComboBox{
										AssignTo:      &vtec.groupSel,
										BindingMember: "Id",
										DisplayMember: "Name",
										MinSize:       Size{Width: 160},
										OnCurrentIndexChanged: func() {
											vtec.listCan()
											vtec.reloadTargetList()
										},
									},
									PushButton{
										Image: "/img/import.ico",
										Text:  "导入",
										OnClicked: func() {

											//判断是否存在分组
											if vtec.groupSel.CurrentIndex() == -1 {
												walk.MsgBox(vtec.mainWin, "", "请选择分组", walk.MsgBoxIconWarning)
												return
											}

											// 分组ID
											ddlCanGroup := vtec.groupSel.Model().([]*dataSource.CanGroupEntity)
											g := ddlCanGroup[vtec.groupSel.CurrentIndex()]

											// 过滤掉已存在的字段
											temp := []*model.CanDetail{}
											for _, v := range vtec.targetCanModel.items {
												if _, isExist := vtec.groupCanMap[v.OutfieldId]; !isExist {
													temp = append(temp, v)
												}
											}

											if len(temp) == 0 {
												walk.MsgBox(vtec.mainWin, "", "无新增字段", walk.MsgBoxIconWarning)
												return
											}

											// 批量插入字段
											err := (&dataSource.VehicleTypeEntity{}).InsertCanFromBigdata(g.Id, temp)
											if err != nil {
												walk.MsgBox(vtec.mainWin, "执行失败", err.Error(), walk.MsgBoxIconError)
												return
											}

											walk.MsgBox(vtec.mainWin, "", "执行成功，如需要编辑，请进编辑界面继续操作", walk.MsgBoxIconInformation)
										},
									},
								},
							},
							TableView{
								AssignTo:           &vtec.targetCanTv,
								AlwaysConsumeSpace: true,
								AlternatingRowBG:   true,
								//AlternatingRowBGColor: walk.RGB(239, 239, 239),
								CheckBoxes:       true,
								ColumnsOrderable: true,
								MultiSelection:   true,
								Columns: []TableViewColumn{
									{Title: "#", Width: 50, Frozen: true},
									{Title: "编号", Width: 60, Alignment: AlignFar},
									{Title: "中文名", Width: 120, Alignment: AlignCenter},
									{Title: "别名", Width: 120, Alignment: AlignCenter},
									{Title: "单位", Width: 50, Alignment: AlignCenter},
									{Title: "数据类型", Width: 80, Alignment: AlignCenter},
									{Title: "转换公式", Width: 80, Alignment: AlignCenter},
									{Title: "数值范围", Width: 160, Alignment: AlignCenter},
									{Title: "小数位", Width: 50, Alignment: AlignCenter},
									{Title: "软报警项", Width: 75, Alignment: AlignCenter, FormatFunc: func(value interface{}) string {
										switch value {
										case "0":
											return ""
										case "1":
											return "报警项"
										default:
											return ""
										}
									}},
									{Title: "可分析项", Width: 85, Alignment: AlignCenter, FormatFunc: func(value interface{}) string {
										switch value {
										case 0:
											return ""
										case 1:
											return "可分析"
										default:
											return ""
										}
									}},
									{Title: "排序", Width: 50, Alignment: AlignCenter},
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
									switch style.Col() {
									case 1:
										idx := strings.Index(item.Note, "已存在")
										if idx != -1 {
											//style.TextColor = walk.RGB(255, 0, 0)
											style.Image = fieldExistIcon
										}
									}
								},
								Model: vtec.targetCanModel,
							},
						},
					},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true, Margins: Margins{Top: 3}},
				Children: []Widget{
					Label{
						AssignTo: &vtec.searchStatLbl,
						Text:     fmt.Sprintf("共 %d 项", len(vtec.searchCanModel.items)),
						Font:     Font{PointSize: 10},
					},
					HSpacer{},
					ImageView{Image: "/img/warn.ico"},
					Label{Text: "：当前车系下，已存在该字段，将不会导入"},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	data, err := (&dataSource.VehicleTypeEntity{TypeId: vt.TypeId}).GetVehicleType()
	if err != nil {
		return vtec, err
	}

	if len(data.CanGroup) > 0 {
		vtec.groupSel.SetModel(data.CanGroup)
		vtec.groupSel.SetCurrentIndex(0)
	}

	vtec.searchCanModel.parent = vtec
	vtec.targetCanModel.parent = vtec

	return vtec, nil
}

// VehicleTypeAddPage 页面对象
type CanAddPage struct {
	*walk.Composite

	vehicleType *model.VehicleTypeStats
	mainWin     *TabMainWindow

	urlCb *walk.ComboBox
	// 搜索关键词模型
	searchModel *model.SearchModel

	// 字段搜索组件
	searchCanTv    *walk.TableView
	searchCanModel *SearchCanTable

	// 目标组件
	groupSel       *walk.ComboBox
	targetCanTv    *walk.TableView
	targetCanModel *TargetCanTable

	// 底部汇总组件
	searchStatLbl *walk.Label

	// 存储当前分组下已存在的Can信息 outfieldId -> *model.CanDetail
	groupCanMap map[string]*model.CanDetail
}

func (vtec *CanAddPage) listCan() {
	ddlData := vtec.groupSel.Model().([]*dataSource.CanGroupEntity)
	groupId := ddlData[vtec.groupSel.CurrentIndex()].Id
	groupCan, _ := (&dataSource.VehicleTypeEntity{}).ListCan(groupId)

	vtec.groupCanMap = make(map[string]*model.CanDetail)
	for _, v := range groupCan {
		vtec.groupCanMap[v.OutfieldId] = v
	}
}

func (vtec *CanAddPage) mappingCan(searchStr string, url string) {
	temp := strings.Split(searchStr, "\r\n")
	vtec.searchModel.Items = []model.SearchItem{}

	for _, item := range temp {
		if strings.TrimSpace(item) == "" {
			continue
		}

		vtec.searchModel.Items = append(vtec.searchModel.Items, model.SearchItem{item, ""})
	}

	result, err := dataSource.ListCanConfig(url)
	if err != nil {
		walk.MsgBox(vtec.mainWin, "", err.Error(), walk.MsgBoxIconError)
		return
	}

	matchCan := []*model.CanDetail{}
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
				matchCan = append(matchCan, &model.CanDetail{
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
			matchCan = append(matchCan, &model.CanDetail{
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
	vtec.searchStatLbl.SetText(fmt.Sprintf("共 %d 项", len(matchCan)))
}

func (vtec *CanAddPage) reloadTargetList() {

	targetCan := []*model.CanDetail{}
	idx := 0
	for i := 0; i < len(vtec.targetCanModel.items); i++ {
		idx = idx + 1
		vtec.targetCanModel.items[i].Index = idx
		vtec.targetCanModel.items[i].Checked = true

		remark := ""
		if _, isExist := vtec.groupCanMap[vtec.targetCanModel.items[i].OutfieldId]; isExist {
			remark = "已存在"
		}
		vtec.targetCanModel.items[i].Note = remark

		targetCan = append(targetCan, vtec.targetCanModel.items[i])
	}

	vtec.targetCanModel.items = targetCan
	vtec.targetCanModel.ResetRows()
}

// 搜索列表项切换选定状态，并同步目标列表
func (vtec *CanAddPage) addOrRemoveTargetItem(rowIndex int, isChecked bool) {
	item := vtec.searchCanModel.items[rowIndex]

	if isChecked {

		remark := ""
		if _, isExist := vtec.groupCanMap[item.OutfieldId]; isExist {
			remark = "已存在"
		}

		vtec.targetCanModel.items = append(vtec.targetCanModel.items, &model.CanDetail{
			Index:       len(vtec.targetCanModel.items) + 1,
			OutfieldId:  item.OutfieldId,
			Unit:        item.Unit,
			Chinesename: item.Chinesename,
			FieldName:   item.FieldName,
			Checked:     true,
			OutfieldSn:  float64(rowIndex),
			Note:        remark,
		})
	} else {
		targetCan := []*model.CanDetail{}
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
}

// 目标列表主动移除行，兵取消搜索列表的勾选状态
func (vtec *CanAddPage) removeTargetItem(rowIndex int, isChecked bool) {
	item := vtec.targetCanModel.items[rowIndex]
	if !isChecked {
		// 删除勾选项
		targetCan := []*model.CanDetail{}
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

		// 重置选定状态
		for i := 0; i < len(vtec.searchCanModel.items); i++ {
			if vtec.searchCanModel.items[i].OutfieldId == item.OutfieldId {
				vtec.searchCanModel.items[i].Checked = false
			}
		}
		vtec.searchCanModel.ResetRows()
	}
}

func (vtec *CanAddPage) refreshItem(can *model.CanDetail) {
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
	items      []*model.CanDetail

	parent *CanAddPage
}

func LoadSearchCan(url string) (*SearchCanTable, error) {
	m := new(SearchCanTable)

	result, err := dataSource.ListCanConfig(url)
	if err != nil {
		return nil, err
	}
	var matchCan []*model.CanDetail
	index := 0

	for j := 0; j < len(result.Props); j++ {
		index++
		matchCan = append(matchCan, &model.CanDetail{
			Index:       index,
			OutfieldId:  result.Props[j].Id,
			Unit:        result.Props[j].Unit,
			Chinesename: result.Props[j].Cn,
			FieldName:   result.Props[j].En})
	}

	m.items = matchCan

	m.ResetRows()
	return m, nil
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
	items      []*model.CanDetail

	parent *CanAddPage
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
