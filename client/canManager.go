package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"ConfigurationTools/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func NewCanManagePanel(parent walk.Container, vt *model.VehicleTypeStats, mainWin *TabMainWindow) (*CanManageForm, error) {
	cmf := &CanManageForm{
		mainWin: mainWin,
	}

	var copyAllCanIv *walk.ImageView
	var exportAllCanIv *walk.ImageView

	if err := (Composite{
		Layout: VBox{},
		//MaxSize: Size{Width: 1280},
		Children: []Widget{
			Composite{
				Layout: HBox{Margins: Margins{Left: 10, Top: 0, Right: 10, Bottom: 0}},
				Children: []Widget{
					ImageView{
						AssignTo:    &copyAllCanIv,
						Image:       "/img/copy.ico",
						ToolTipText: "复制到剪贴板",
						MinSize:     Size{Width: 25},
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							copyAllCanIv.SetEnabled(false)
							defer copyAllCanIv.SetEnabled(true)

							data, err := (&dataSource.VehicleTypeEntity{}).ListAllCan(vt.TypeId)
							if err != nil {
								walk.MsgBox(cmf.mainWin, "", err.Error(), walk.MsgBoxIconError)
								return
							}

							jsonData, err := json.MarshalIndent(data, "", "\t")
							if err != nil {
								copyAllCanIv.SetEnabled(true)
								walk.MsgBox(cmf.mainWin, "", err.Error(), walk.MsgBoxIconError)
								return
							}

							err = walk.Clipboard().SetText(string(jsonData))
							if err != nil {
								copyAllCanIv.SetEnabled(true)
								walk.MsgBox(cmf.mainWin, "", err.Error(), walk.MsgBoxIconError)
								return
							}

							copyAllCanIv.SetEnabled(true)
							mynotify.Info("已复制至剪贴板")
						},
					},
					ImageView{
						AssignTo:    &exportAllCanIv,
						Image:       "/img/export.ico",
						ToolTipText: "导出CSV",
						MinSize:     Size{Width: 25},
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							exportAllCanIv.SetEnabled(false)

							dlg := new(walk.FileDialog)
							dlg.Title = "位置"

							if ok, err := dlg.ShowBrowseFolder(cmf.mainWin); err != nil {
								exportAllCanIv.SetEnabled(true)
								walk.MsgBox(cmf.mainWin, "", "打开文件选择器失败："+err.Error(), walk.MsgBoxIconError)
								return
							} else if !ok {
								exportAllCanIv.SetEnabled(true)
								return
							}

							data, err := (&dataSource.VehicleTypeEntity{}).ListAllCan(vt.TypeId)
							if err != nil {
								exportAllCanIv.SetEnabled(true)
								walk.MsgBox(cmf.mainWin, "", err.Error(), walk.MsgBoxIconError)
							}

							buffer := [][]string{}
							buffer = append(buffer, []string{"分组名称", "分组排序", "分组标识", "别名", "字段名", "CAN编码", "单位", "数据类型", "转换公式",
								"小数位", "值域", "是否为软报警字段", "是否为可分析字段", "是否删除", "CAN排序"})

							for i := 0; i < len(data); i++ {
								data := []string{
									data[i].GroupName,
									strconv.Itoa(data[i].Sort),
									data[i].Remark,
									data[i].Chinesename,
									data[i].FieldName,
									data[i].OutfieldId,
									data[i].Unit,
									data[i].DataType,
									data[i].Formula,
									data[i].Decimals,
									data[i].DataMap,
									data[i].IsAlarm,
									strconv.Itoa(data[i].IsAnalysable),
									strconv.Itoa(data[i].IsDelete),
									strconv.FormatFloat(data[i].OutfieldSn, 'f', 2, 64),
								}
								buffer = append(buffer, data)
							}

							fileName := fmt.Sprintf("%s\\%s_%s_%s.csv", dlg.FilePath, vt.TypeName, vt.OrgName, time.Now().Format("20060102150405"))
							file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0666)
							if err != nil {
								exportAllCanIv.SetEnabled(true)
								walk.MsgBox(cmf.mainWin, "", err.Error(), walk.MsgBoxIconError)
								return
							}
							defer file.Close()

							file.WriteString("\xEF\xBB\xBF")

							writer := csv.NewWriter(file)
							writer.WriteAll(buffer)
							writer.Flush()
							file.Close()

							exportAllCanIv.SetEnabled(true)
							mynotify.Message("导出完毕")
						},
					},
					ImageView{
						Image:       "/img/add.ico",
						ToolTipText: "添加 Can",
						MinSize:     Size{Width: 25},
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							cd := &canDetail{}
							if cmd, err := cd.runCanDetailDialog(cmf.mainWin); err != nil {
								walk.MsgBox(cmf.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
								return
							} else if cmd == walk.DlgCmdOK {

								canGroupModel := cmf.canGroupCB.Model().([]*dataSource.CanGroupEntity)
								groupId := canGroupModel[cmf.canGroupCB.CurrentIndex()].Id

								_, err := (&dataSource.GroupCanDetail{
									Id:           cd.Id,
									OutfieldId:   cd.Key,
									Alias:        cd.Alias,
									FieldName:    cd.Cname,
									Formula:      cd.Formula,
									DataType:     cd.DataType,
									Unit:         cd.Unit,
									DataScope:    cd.DataScope,
									Decimals:     cd.Prec,
									Sort:         cd.Sort,
									IsAlarm:      cd.IsAlarm == 1,
									IsAnalysable: cd.IsAnalysable == 1,
									IsDelete:     cd.IsDelete == 1,
									GroupInfoId:  groupId,
								}).Add()

								if err != nil {
									mynotify.Error("添加失败，" + err.Error())
									return
								}

								cmf.loadCan()
							}
						}},
					ImageView{
						Image:       "/img/delete.ico",
						ToolTipText: "删除选定行",
						MinSize:     Size{Width: 25},
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							selectedIndex := cmf.canListTv.SelectedIndexes()
							if len(selectedIndex) == 0 {
								walk.MsgBox(cmf.mainWin, "", "请选择数据行", walk.MsgBoxIconInformation)
								return
							}

							model := cmf.canListTv.Model().([]*model.CanDetail)
							canId := []string{}
							for _, idx := range selectedIndex {
								item := model[idx]
								canId = append(canId, item.Id)
							}

							_, err := (&dataSource.VehicleTypeEntity{}).DeleteCanField(canId, true)
							if err != nil {
								walk.MsgBox(cmf.mainWin, "删除失败", err.Error(), walk.MsgBoxIconError)
								return
							}

							cmf.loadCan()
						},
					},
					ImageView{
						Image:       "/img/undo.ico",
						ToolTipText: "撤销删除",
						MinSize:     Size{Width: 25},
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							selectedIndex := cmf.canListTv.SelectedIndexes()
							if len(selectedIndex) == 0 {
								walk.MsgBox(cmf.mainWin, "", "请选择数据行", walk.MsgBoxIconInformation)
								return
							}

							model := cmf.canListTv.Model().([]*model.CanDetail)
							canId := []string{}
							for _, idx := range selectedIndex {
								item := model[idx]
								canId = append(canId, item.Id)
							}

							_, err := (&dataSource.VehicleTypeEntity{}).CancelDelete(canId)
							if err != nil {
								walk.MsgBox(cmf.mainWin, "撤销失败", err.Error(), walk.MsgBoxIconError)
								return
							}

							cmf.loadCan()
						},
					},
					ImageView{
						Image:       "/img/cancel.ico",
						ToolTipText: "彻底删除选定行",
						MinSize:     Size{Width: 25},
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							selectedIndex := cmf.canListTv.SelectedIndexes()
							if len(selectedIndex) == 0 {
								walk.MsgBox(cmf.mainWin, "", "请选择数据行", walk.MsgBoxIconInformation)
								return
							}

							model := cmf.canListTv.Model().([]*model.CanDetail)
							canId := []string{}
							for _, idx := range selectedIndex {
								item := model[idx]
								canId = append(canId, item.Id)
							}

							_, err := (&dataSource.VehicleTypeEntity{}).DeleteCanField(canId, false)
							if err != nil {
								walk.MsgBox(cmf.mainWin, "删除失败", err.Error(), walk.MsgBoxIconError)
								return
							}

							cmf.loadCan()
						},
					},
					HSpacer{},
					ComboBox{
						AssignTo:      &cmf.canGroupCB,
						BindingMember: "Id",
						DisplayMember: "Name",
						MinSize:       Size{Width: 160},
						OnCurrentIndexChanged: func() {
							cmf.loadCan()
						},
					},
					PushButton{
						Text: "查询",
						OnClicked: func() {
							cmf.loadCan()
						},
					},
				},
			},
			TableView{
				AssignTo:         &cmf.canListTv,
				AlternatingRowBG: true,
				//AlternatingRowBGColor: walk.RGB(239, 239, 239),
				ColumnsOrderable: true,
				//CheckBoxes:            true,
				MultiSelection: true,
				Columns: []TableViewColumn{
					{Name: "Index", Title: "#", Frozen: true, Alignment: AlignCenter, Width: 60},
					{Name: "OutfieldId", Title: "编号", Width: 60, Alignment: AlignFar},
					{Name: "FieldName", Title: "中文名", Alignment: AlignCenter, Width: 120},
					{Name: "Chinesename", Title: "别名", Alignment: AlignCenter, Width: 120},
					{Name: "Unit", Title: "单位", Alignment: AlignCenter, Width: 50},
					{Name: "DataType", Title: "数据类型", Alignment: AlignCenter, Width: 80, FormatFunc: func(value interface{}) string {
						switch value.(string) {
						case "1":
							return "日期时间"
						case "2":
							return "数字枚举"
						case "3":
							return "数据"
						case "5":
							return "文本枚举"
						case "6":
							return "文本多枚举"
						case "7":
							return "多字段组合多枚举"
						default:
							return "原始值"
						}
					}},
					{Name: "Formula", Title: "转换公式", Alignment: AlignCenter, Width: 80},
					{Name: "DataMap", Title: "数值范围", Alignment: AlignCenter, Width: 160},
					{Name: "Decimals", Title: "小数位", Alignment: AlignCenter, Width: 50, FormatFunc: func(value interface{}) string {
						switch value.(string) {
						case "":
							return ""
						case "0":
							return ""
						default:
							return value.(string)
						}
					}},
					{Name: "IsAlarm", Title: "软报警项", Alignment: AlignCenter, Width: 75, FormatFunc: func(value interface{}) string {
						switch value.(string) {
						case "0":
							return ""
						case "1":
							return "√"
						default:
							return ""
						}
					}},
					{Name: "IsAnalysable", Title: "可分析项", Alignment: AlignCenter, Width: 85, FormatFunc: func(value interface{}) string {
						switch value.(int) {
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
					{Name: "IsDelete", Title: "是否删除", Alignment: AlignCenter, Width: 85, FormatFunc: func(value interface{}) string {
						switch value.(int) {
						case 0:
							return ""
						case 1:
							return "√"
						default:
							return ""
						}
					}},
					{Name: "OutfieldSn", Title: "排序", Alignment: AlignCenter, Width: 50},
				},
				OnItemActivated: func() {
					model := cmf.canListTv.Model().([]*model.CanDetail)
					detail := model[cmf.canListTv.CurrentIndex()]

					dt, err := strconv.Atoi(detail.DataType)
					if err != nil {
						dt = 4
					}

					prec, err := strconv.Atoi(detail.Decimals)
					if err != nil {
						prec = 0
					}

					isAlarm, err := strconv.Atoi(detail.IsAlarm)
					if err != nil {
						isAlarm = 0
					}

					cd := &canDetail{
						Id:           detail.Id,
						Key:          detail.OutfieldId,
						Unit:         detail.Unit,
						Sort:         detail.OutfieldSn,
						GroupInfoId:  detail.GroupInfoId,
						Cname:        detail.FieldName,
						Alias:        detail.Chinesename,
						Formula:      detail.Formula,
						DataType:     dt,
						Prec:         prec,
						DataScope:    detail.DataMap,
						IsAlarm:      isAlarm,
						IsAnalysable: detail.IsAnalysable,
						IsDelete:     detail.IsDelete,
					}
					if cmd, err := cd.runCanDetailDialog(cmf.mainWin); err != nil {
						walk.MsgBox(cmf.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
						return
					} else if cmd == walk.DlgCmdOK {

						_, err := (&dataSource.GroupCanDetail{
							Id:           cd.Id,
							OutfieldId:   cd.Key,
							Alias:        cd.Alias,
							FieldName:    cd.Cname,
							Formula:      cd.Formula,
							DataType:     cd.DataType,
							Unit:         cd.Unit,
							DataScope:    cd.DataScope,
							Decimals:     cd.Prec,
							Sort:         cd.Sort,
							IsAlarm:      cd.IsAlarm == 1,
							IsAnalysable: cd.IsAnalysable == 1,
							IsDelete:     cd.IsDelete == 1,
							GroupInfoId:  cd.GroupInfoId,
						}).Update()

						if err != nil {
							walk.MsgBox(cmf.mainWin, "更新失败", err.Error(), walk.MsgBoxIconError)
							return
						}

						cmf.loadCan()
					}
				},
				StyleCell: func(style *walk.CellStyle) {
					model := cmf.canListTv.Model().([]*model.CanDetail)
					item := model[style.Row()]

					if item.IsDelete != 1 {
						return
					}

					style.TextColor = walk.RGB(255, 0, 0)
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	group, err := (&dataSource.VehicleTypeEntity{TypeId: vt.TypeId}).GetVehicleType()
	if err != nil {
		return cmf, err
	}

	if len(group.CanGroup) > 0 {
		cmf.canGroupCB.SetModel(group.CanGroup)
		cmf.canGroupCB.SetCurrentIndex(0)

		cmf.loadCan()
	}

	return cmf, nil
}

type CanManageForm struct {
	*walk.Composite
	mainWin *TabMainWindow

	canGroupCB *walk.ComboBox
	canListTv  *walk.TableView

	statusLbl *walk.Label
	detail    *model.CanDetail
}

func (cmf *CanManageForm) loadCan() {
	canGroupModel := cmf.canGroupCB.Model().([]*dataSource.CanGroupEntity)
	groupId := canGroupModel[cmf.canGroupCB.CurrentIndex()].Id
	cans, err := (&dataSource.VehicleTypeEntity{}).ListCan(groupId)
	if err != nil {
		mynotify.Error("can字段加载失败：" + err.Error())
		return
	}
	cmf.canListTv.SetModel(cans)
}

type canDetail struct {
	Id           string
	Key          string
	Unit         string
	Sort         float64
	GroupInfoId  string
	Cname        string
	Alias        string
	Formula      string
	DataType     int
	Prec         int
	DataScope    string
	IsAlarm      int
	IsAnalysable int
	IsDelete     int
}

func (cd *canDetail) runCanDetailDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	var title string
	if cd.Id == "" {
		title = "添加 CAN"
	} else {
		title = "编辑 CAN"
	}

	return Dialog{
		AssignTo:      &dlg,
		Title:         title,
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "canDetail",
			DataSource:     cd,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{Width: 500, Height: 460},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: "CAN 编号"},
					LineEdit{Text: Bind("Key", Regexp{"^[\\s|\\S]{1,}$"})},
					Label{Text: "中文名"},
					LineEdit{Text: Bind("Cname")},
					Label{Text: "别名"},
					LineEdit{Text: Bind("Alias")},
					Label{Text: "单位"},
					LineEdit{Text: Bind("Unit")},
					Label{Text: "数据类型"},
					ComboBox{
						BindingMember: "Code",
						DisplayMember: "Name",
						Model:         utils.KnownDataType(),
						Value:         Bind("DataType"),
						MaxSize:       Size{Width: 120},
						MinSize:       Size{Width: 120},
					},
					Label{Text: "转换公式"},
					LineEdit{Text: Bind("Formula")},
					Label{Text: "数值范围"},
					LineEdit{Text: Bind("DataScope")},
					Label{Text: "小数位"},
					NumberEdit{
						Value: Bind("Prec"), //cans[i].Decimals,
					},
					Label{Text: "报警项"},
					ComboBox{
						BindingMember: "Code",
						DisplayMember: "Name",
						Model:         utils.KnownAlarm(),
						Value:         Bind("IsAlarm"),
					},
					Label{Text: "可分析项"},
					ComboBox{
						BindingMember: "Code",
						DisplayMember: "Name",
						Model:         utils.KnownAnaly(),
						Value:         Bind("IsAnalysable"),
					},
					Label{Text: "排序"},
					NumberEdit{
						Value:    Bind("Sort"),
						Decimals: 2,
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
