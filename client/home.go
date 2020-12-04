package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
	"strconv"
	"time"
)

type VehicleTypeCom struct {
	*walk.Composite
	parent        walk.Container
	mainWin       *TabMainWindow
	vehicleSeries *walk.ScrollView
	//addVehicleType *walk.ImageView
	statLbl *walk.Label

	searchKeyEdit *walk.LineEdit
}

func NewVehicleTypeList(parent walk.Container, mw *TabMainWindow) (*VehicleTypeCom, error) {
	vtc, err := initList(parent)
	if err != nil {
		return nil, err
	}
	vtc.mainWin = mw
	vtc.parent = parent

	err = vtc.loadVehicleType()
	if err != nil {
		return vtc, err
	}
	return vtc, nil
}

// 初始化首页
func initList(parent walk.Container) (*VehicleTypeCom, error) {
	vtc := &VehicleTypeCom{}

	if err := (Composite{
		AssignTo: &vtc.Composite,
		Layout:   VBox{
			//MarginsZero: true,
		},
		Children: []Widget{
			Composite{
				Layout: HBox{
					MarginsZero: true,
					SpacingZero: true,
				},
				Children: []Widget{
					Composite{
						Layout: HBox{
							Margins:     Margins{Right: 10},
							MarginsZero: true,
							SpacingZero: true,
						},
						Children: []Widget{
							Label{
								Text: "车辆类型",
								Font: Font{PointSize: 11},
							},
						},
					},
					ImageView{
						Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						Image:      "img/add.ico",
						Margin:     2,
						OnMouseDown: func(x, y int, button walk.MouseButton) {
							vtc.mainWin.AddVehicleType()
						},
					},
					HSpacer{},
					Composite{
						Layout: HBox{MarginsZero: true, Margins: Margins{Right: 5}},
						Children: []Widget{
							LineEdit{
								AssignTo:  &vtc.searchKeyEdit,
								Alignment: AlignHFarVCenter,
								MaxSize:   Size{Width: 180},
								OnKeyDown: func(key walk.Key) {
									if key == walk.KeyReturn {
										vtc.loadVehicleType()
									}
								},
							},
						},
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							PushButton{
								Text: "查询",
								//Image: "/img/search.ico",
								OnClicked: func() {
									vtc.loadVehicleType()
								},
							},
						},
					},
				},
			},
			ScrollView{
				AssignTo: &vtc.vehicleSeries,
				Layout: Flow{
					//MarginsZero: true,
					Alignment: AlignHNearVNear,
				},
				Children: []Widget{},
			},
			Composite{
				Layout: HBox{
					Margins:     Margins{Top: 3},
					MarginsZero: true,
					SpacingZero: true,
				},
				Children: []Widget{
					Label{
						AssignTo: &vtc.statLbl,
						Text:     "",
						Font:     Font{PointSize: 10},
					},
					HSpacer{},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	//vtc.addVehicleType.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
	//	if button == walk.LeftButton {
	//		//vtc.AddVehicleType()
	//		vtc.mainWin.AddVehicleType()
	//	}
	//})

	return vtc, nil
}

// 初始化列表项
func (vtc *VehicleTypeCom) loadVehicleType() error {

	searchKey := vtc.searchKeyEdit.Text()
	vts, err := (&dataSource.OrganizationEntity{SearchKey: searchKey}).ListVehicleType()
	if err != nil {
		return err
	}

	vtc.vehicleSeries.SetSuspended(true)
	defer vtc.vehicleSeries.SetSuspended(false)

	c := vtc.vehicleSeries.Children()
	for i := c.Len() - 1; i >= 0; i-- {
		item := c.At(i)
		item.SetParent(nil)
		item.Dispose()
	}

	for i := 0; i < len(vts); i++ {

		s := new(VehicleTypeItem)

		var db *walk.DataBinder
		data := &model.VehicleTypeStats{
			TypeId:   vts[i].TypeId,
			TypeName: vts[i].TypeName,
			OrgId:    vts[i].OrgId,
			OrgName:  vts[i].OrgName,
			Group:    vts[i].Group,
		}

		contextMenu := []MenuItem{
			Action{
				Text: "编辑车系",
				OnTriggered: func() {
					vtc.mainWin.EditVehicleType(db.DataSource().(*model.VehicleTypeStats))
				},
			},
			Action{
				Text: "CAN 管理",
				OnTriggered: func() {
					vtc.mainWin.CanManage(db.DataSource().(*model.VehicleTypeStats))
				},
			},
			Separator{},
			Action{
				Text: "批量编辑",
				OnTriggered: func() {
					vtc.mainWin.EditGroupCan(db.DataSource().(*model.VehicleTypeStats))
				},
			},
			Separator{},
			Action{
				Text: "从大数据导入",
				OnTriggered: func() {
					vtc.mainWin.ImportCanFromBigData(db.DataSource().(*model.VehicleTypeStats))
				},
			},
			Action{
				Text: "从其他车型导入",
				OnTriggered: func() {
					vtc.mainWin.ImportCanFromExisting(db.DataSource().(*model.VehicleTypeStats))
				},
			},
			Separator{},
			Action{
				Text: "复制至剪贴板",
				OnTriggered: func() {
					vt := db.DataSource().(*model.VehicleTypeStats)

					data, err := (&dataSource.VehicleTypeEntity{}).ListAllCan(vt.TypeId)
					if err != nil {
						mynotify.Error(err.Error())
						return
					}

					jsonData, err := json.MarshalIndent(data, "", "\t")
					if err != nil {
						mynotify.Error(err.Error())
						return
					}

					err = walk.Clipboard().SetText(string(jsonData))
					if err != nil {
						mynotify.Error(err.Error())
						return
					}

					mynotify.Info("已复制至剪贴板")
				},
			},
			Action{
				Text: "从剪贴板导入",
				OnTriggered: func() {
					vtc.mainWin.NewImportCanFromClipboard(db.DataSource().(*model.VehicleTypeStats))
				},
			},
			Separator{},
			Action{
				Text: "导出CSV文件",
				OnTriggered: func() {
					vt := db.DataSource().(*model.VehicleTypeStats)
					err := vtc.exportCan(vt)
					if err != nil {
						mynotify.Error(err.Error())
						return
					}
				},
			},
			Action{
				Text: "从CSV文件导入",
				OnTriggered: func() {
					vtc.mainWin.NewImportCanFromCsvFile(db.DataSource().(*model.VehicleTypeStats))
				},
			},
		}

		w := Composite{
			AssignTo: &s.Composite,
			DataBinder: DataBinder{
				AssignTo:   &db,
				DataSource: data,
			},

			ContextMenuItems: contextMenu,
			Layout:           VBox{},
			Background:       SolidColorBrush{Color: walk.RGB(0x6F, 0xAF, 0x9D)},
			MaxSize:          Size{Width: 410, Height: 155},
			MinSize:          Size{Width: 200, Height: 155},

			Children: []Widget{
				Composite{
					Layout:    HBox{MarginsZero: true, SpacingZero: true},
					Alignment: AlignHNearVCenter,
					Children: []Widget{
						Label{
							AssignTo:  &s.VehicleTypeNameLbl,
							Font:      Font{PointSize: 11},
							Text:      Bind("TypeName"),
							TextColor: walk.RGB(0xF5, 0xF5, 0xF5),
						},
					},
				},
				Composite{
					AssignTo:  &s.scrollCmp,
					Alignment: AlignHNearVNear,
					//MinSize:   Size{Height: 25},
					//MaxSize:   Size{Height: 80},
					Layout:   Flow{MarginsZero: true, SpacingZero: true},
					Children: []Widget{},
				},
				Composite{
					Layout:    HBox{MarginsZero: true, SpacingZero: true},
					Alignment: AlignHFarVCenter,
					//Border:    true,
					Children: []Widget{
						HSpacer{},
						Label{
							AssignTo:     &s.OrgNameLbl,
							Font:         Font{PointSize: 10},
							MaxSize:      Size{Width: 160},
							Text:         Bind("OrgName"),
							TextColor:    walk.RGB(0xF5, 0xF5, 0xF5),
							EllipsisMode: EllipsisEnd,
						},
					},
				},
			},
		}

		err := w.Create(NewBuilder(vtc.vehicleSeries))
		if err != nil {
			panic(err)
		}

		for _, v := range vts[i].Group {
			cangroup := Composite{
				Layout:  HBox{Margins: Margins{Top: 5, Right: 5}, MarginsZero: true, SpacingZero: true},
				MinSize: Size{Width: 60},
				MaxSize: Size{Width: 100},
				Children: []Widget{
					Label{
						Text:         v.Name,
						Font:         Font{PointSize: 9},
						MinSize:      Size{Width: 55},
						TextColor:    walk.RGB(0xF8, 0xF8, 0xFF),
						EllipsisMode: EllipsisEnd,
					},
					Label{
						Font:      Font{PointSize: 10, Underline: true},
						Text:      strconv.Itoa(v.Count),
						TextColor: walk.RGB(0xFf, 0xFF, 0xFf),
					},
				},
			}
			cangroup.Create(NewBuilder(s.scrollCmp))
		}

		HSpacer{}.Create(NewBuilder(s.scrollCmp))
	}

	HSpacer{}.Create(NewBuilder(vtc.vehicleSeries))

	vtc.statLbl.SetText(fmt.Sprintf("共 %d 项", len(vts)))
	return nil
}

type VehicleTypeItem struct {
	*walk.Composite

	VehicleTypeNameLbl *walk.Label
	OrgNameLbl         *walk.Label

	scrollCmp *walk.Composite

	GeneralCmp *walk.Composite
	GeneralLbl *walk.Label

	CanCmp *walk.Composite
	CanLbl *walk.Label

	IssueCmp *walk.Composite
	IssueLbl *walk.Label
}

func (vtc *VehicleTypeCom) exportCan(vt *model.VehicleTypeStats) error {

	dlg := new(walk.FileDialog)
	dlg.Title = "位置"

	if ok, err := dlg.ShowBrowseFolder(vtc.mainWin); err != nil {
		return errors.New("打开文件选择器失败：" + err.Error())
	} else if !ok {
		return nil
	}

	data, err := (&dataSource.VehicleTypeEntity{}).ListAllCan(vt.TypeId)
	if err != nil {
		return err
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
		return err
	}
	defer file.Close()

	file.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(file)
	writer.WriteAll(buffer)
	writer.Flush()

	file.Close()

	mynotify.Message("导出完毕")
	return nil
}
