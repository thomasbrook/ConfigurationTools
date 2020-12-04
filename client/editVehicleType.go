package client

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func NewEditPanel(parent walk.Container, vt *model.VehicleType, mainWin *TabMainWindow) (*VehicleTypeEditPage, error) {
	rand.Seed(time.Now().UnixNano())

	vtec := &VehicleTypeEditPage{
		canCfgTable:  &CanConfigTable{items: []*model.CanConfig{}},
		editCanModel: &model.CanConfig{},
		vehicleType:  vt,
		mainWin:      mainWin,
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
						StretchFactor: 7,
						Children: []Widget{Composite{
							Layout: VBox{MarginsZero: true},
							Children: []Widget{
								Composite{
									Layout: HBox{MarginsZero: true, SpacingZero: true},
									Children: []Widget{
										ComboBox{
											AssignTo:      &vtec.groupSel,
											Alignment:     AlignHNearVCenter,
											BindingMember: "Id",
											DisplayMember: "Name",
										},
									},
								},
								TableView{
									AssignTo: &vtec.vehicleTypeCanTv,
									//AlternatingRowBGColor: walk.RGB(239, 239, 239),
									CheckBoxes:       false,
									ColumnsOrderable: true,
									MultiSelection:   true,
									Columns: []TableViewColumn{
										{Title: "#", Width: 35},
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
										item := vtec.canCfgTable.items[style.Row()]

										if item.Checked {
											if style.Row()%2 == 0 {
												style.BackgroundColor = walk.RGB(159, 215, 255)
											} else {
												style.BackgroundColor = walk.RGB(143, 199, 239)
											}
										}

									},
									Model: vtec.canCfgTable,
									OnSelectedIndexesChanged: func() {
										idxs := vtec.vehicleTypeCanTv.SelectedIndexes()

										if len(idxs) > 0 {
											editCanModel := model.CanConfig{
												Id:           vtec.canCfgTable.items[idxs[0]].Id,
												Index:        vtec.canCfgTable.items[idxs[0]].Index,
												OutfieldId:   vtec.canCfgTable.items[idxs[0]].OutfieldId,
												Unit:         vtec.canCfgTable.items[idxs[0]].Unit,
												OutfieldSn:   vtec.canCfgTable.items[idxs[0]].OutfieldSn,
												Chinesename:  vtec.canCfgTable.items[idxs[0]].Chinesename,
												Formula:      vtec.canCfgTable.items[idxs[0]].Formula,
												DataType:     vtec.canCfgTable.items[idxs[0]].DataType,
												FieldName:    vtec.canCfgTable.items[idxs[0]].FieldName,
												Decimals:     vtec.canCfgTable.items[idxs[0]].Decimals,
												DataMap:      vtec.canCfgTable.items[idxs[0]].DataMap,
												IsAlarm:      vtec.canCfgTable.items[idxs[0]].IsAlarm,
												IsAnalysable: vtec.canCfgTable.items[idxs[0]].IsAnalysable,
											}

											//	kind := reflect.ValueOf(editCanModel).Kind()

											//	log.Print(fmt.Sprintf("%+v", kind))
											//	log.Print(fmt.Sprintf("%+v", kind != reflect.Func))
											//	log.Print(fmt.Sprintf("%+v", kind != reflect.Map))
											//	log.Print(fmt.Sprintf("%+v", kind != reflect.Slice))
											//	log.Print(fmt.Sprintf("%+v", kind == reflect.ValueOf(mw.editCanDb.DataSource()).Kind()))
											//	log.Print(fmt.Sprintf("%+v", editCanModel == mw.editCanDb.DataSource()))

											//	kind != reflect.Func && kind != reflect.Map && kind != reflect.Slice &&
											//	kind == reflect.ValueOf(db.dataSource).Kind() && dataSource == db.dataSource

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
						},
					},
					Composite{
						Layout:        VBox{Margins: Margins{5, 5, 5, 5}},
						StretchFactor: 3,
						Children: []Widget{
							ScrollView{
								Layout: Flow{
									MarginsZero: false,
								},
								Children: []Widget{
									Composite{
										AssignTo:  &ef.form,
										Layout:    Grid{Columns: 2, MarginsZero: false},
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
									},
								},
							},
							Composite{
								Layout: HBox{},
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text: "删除",
										OnClicked: func() {
											cmd := walk.MsgBox(vtec.mainWin, "提示", "是否确认删除？", walk.MsgBoxOKCancel)

											if cmd == walk.DlgCmdOK {

												ds := vtec.editCanDb.DataSource()
												data := ds.(*model.CanConfig)

												_, err := dataSource.DeleteCan(data.Id)
												if err != nil {
													panic(err)
												}
												vtec.loadCan()
											}
										},
									},
									PushButton{
										Text: "更新",
										OnClicked: func() {
											can, err := ef.set()
											if err != nil {
												panic(err)
											}

											ds := vtec.editCanDb.DataSource()
											data := ds.(*model.CanConfig)
											can.Id = data.Id

											_, err = dataSource.UpdateCan(can)
											if err != nil {
												panic(err)
											}
											vtec.loadCan()
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
					},
					HSpacer{},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	vtec.groupSel.CurrentIndexChanged().Attach(func() {
		vtec.loadCan()
	})

	group := vtec.canGroupForEdit()
	if len(group) > 0 {
		vtec.groupSel.SetModel(group)
		vtec.groupSel.SetCurrentIndex(0)
	}

	return vtec, nil
}

type VehicleTypeEditPage struct {
	*walk.Composite

	// 外部数据
	vehicleType *model.VehicleType
	mainWin     *TabMainWindow

	// 页面组件
	groupSel         *walk.ComboBox
	vehicleTypeCanTv *walk.TableView
	editCanForm      *walk.Composite
	editCanDb        *walk.DataBinder
	searchStatLbl    *walk.Label

	// 数据模型
	canCfgTable  *CanConfigTable
	editCanModel *model.CanConfig
}

func (vtec *VehicleTypeEditPage) canGroupForEdit() []*model.CanGroup {
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

func (vtec *VehicleTypeEditPage) loadCan() {
	g := vtec.canGroupForEdit()[vtec.groupSel.CurrentIndex()]
	cans, err := dataSource.ListCan(vtec.vehicleType.TypeId, g.Id)
	if err != nil {
		log.Fatal(err)
	}

	bigCan := dataSource.ListCanConfig()

	kv := make(map[string]string)

	// 填充can表格
	matchCan := []*model.CanConfig{}
	index := 1

	for j := 0; j < len(cans); j++ {
		for i := 0; i < len(bigCan.Props); i++ {
			if bigCan.Props[i].Id == cans[j].OutfieldId {

				c := &model.CanConfig{
					Id:           cans[j].Id,
					Index:        index,
					OutfieldId:   cans[j].OutfieldId,
					OutfieldSn:   cans[j].OutfieldSn,
					Formula:      cans[j].Formula,
					DataType:     cans[j].DataType,
					FieldName:    bigCan.Props[i].En,
					Decimals:     cans[j].Decimals,
					DataMap:      cans[j].DataMap,
					IsAlarm:      cans[j].IsAlarm,
					IsAnalysable: cans[j].IsAnalysable,
				}

				if strings.TrimSpace(cans[j].Unit) == "" {
					c.Unit = bigCan.Props[i].Unit
				} else {
					c.Unit = cans[j].Unit
				}

				if strings.TrimSpace(cans[j].Chinesename) == "" {
					c.Chinesename = bigCan.Props[i].Cn
				} else {
					c.Chinesename = cans[j].Chinesename
				}

				matchCan = append(matchCan, c)
				_, ok := kv[c.OutfieldId]
				if !ok {
					kv[c.OutfieldId] = c.Chinesename
				}
				index++
				break
			}
		}
	}
	vtec.canCfgTable.items = matchCan
	vtec.vehicleTypeCanTv.SetModel(vtec.canCfgTable)
	vtec.searchStatLbl.SetText(strconv.Itoa(len(matchCan)) + " 项")
}

type CanConfigTable struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*model.CanConfig
}

// Called by the TableView from SetModel and every time the model publishes a
// RowsReset event.
func (m *CanConfigTable) RowCount() int {
	return len(m.items)
}

// Called by the TableView when it needs the text to display for a given cell.
func (m *CanConfigTable) Value(row, col int) interface{} {
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

// Called by the TableView to retrieve if a given row is checked.
func (m *CanConfigTable) Checked(row int) bool {
	return m.items[row].Checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *CanConfigTable) SetChecked(row int, checked bool) error {
	m.items[row].Checked = checked

	return nil
}

// Called by the TableView to sort the model.
func (m *CanConfigTable) Sort(col int, order walk.SortOrder) error {
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

func (m *CanConfigTable) ResetRows() {
	// Notify TableView and other interested parties about the reset.
	m.PublishRowsReset()

	m.Sort(m.sortColumn, m.sortOrder)
}
