package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"ConfigurationTools/utils"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func NewEditCanPanel(parent walk.Container, vt *model.VehicleTypeStats, mainWin *TabMainWindow) (*CanListForm, error) {
	rand.Seed(time.Now().UnixNano())

	clf := &CanListForm{}

	if err := (Composite{
		Layout:  VBox{},
		MaxSize: Size{Width: 1280},
		Children: []Widget{
			Composite{
				Layout: HBox{Margins: Margins{Left: 20, Top: 0, Right: 10, Bottom: 0}},
				Children: []Widget{
					HSpacer{},
					ComboBox{
						AssignTo:      &clf.CanGroupDDL,
						BindingMember: "Id",
						DisplayMember: "Name",
						MinSize:       Size{Width: 160},
						OnKeyDown: func(key walk.Key) {
							if key == walk.KeyReturn {
								clf.renderCanList()
							}
						},
					},
					PushButton{
						Text: "查询",
						//Image: "/img/search.ico",
						OnClicked: func() {
							clf.renderCanList()
						},
					},
				},
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					HSpacer{},
					Label{
						Text:    "#",
						MaxSize: Size{Width: 30},
						MinSize: Size{Width: 30},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "编号",
						MaxSize: Size{Width: 60},
						MinSize: Size{Width: 60},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "字段名称",
						MaxSize: Size{Width: 120},
						MinSize: Size{Width: 120},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "别名",
						MaxSize: Size{Width: 120},
						MinSize: Size{Width: 120},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "单位",
						MaxSize: Size{Width: 50},
						MinSize: Size{Width: 50},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "数据类型",
						MaxSize: Size{Width: 120},
						MinSize: Size{Width: 120},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "转换公式",
						MaxSize: Size{Width: 80},
						MinSize: Size{Width: 80},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "数值范围",
						MaxSize: Size{Width: 200},
						MinSize: Size{Width: 200},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "小数位",
						MaxSize: Size{Width: 50},
						MinSize: Size{Width: 50},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "软报警项",
						MaxSize: Size{Width: 75},
						MinSize: Size{Width: 75},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "可分析项",
						MaxSize: Size{Width: 85},
						MinSize: Size{Width: 85},
						Font:    Font{PointSize: 10, Bold: true},
					},
					Label{
						Text:    "排序",
						MaxSize: Size{Width: 50},
						MinSize: Size{Width: 50},
						Font:    Font{PointSize: 10, Bold: true},
					},
					HSpacer{},
				},
			},
			VSeparator{
				MinSize: Size{Height: 1},
				MaxSize: Size{Height: 1},
			},
			ScrollView{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					ScrollView{
						AssignTo:      &clf.CanListScroll,
						Layout:        VBox{},
						VerticalFixed: true,
						Children:      []Widget{},
					},
				},
			},
			Composite{
				AssignTo: &clf.btnCom,
				Visible:  false,
				Layout:   HBox{Margins: Margins{Left: 20, Top: 10, Right: 10, Bottom: 0}},
				Children: []Widget{
					Label{
						AssignTo: &clf.statsLbl,
						Text:     "",
						Font:     Font{PointSize: 10},
					},
					HSpacer{},
					PushButton{
						Text: "保存修改",
						//Image: "/img/update.ico",
						OnClicked: func() {

							if len(clf.Table) == 0 {
								walk.MsgBox(clf.mainWin, "", "数据不存在", walk.MsgBoxIconWarning)
								return
							}

							can := []*model.CanDetail{}
							for _, v := range clf.Table {

								idx := v.DataTypeCb.CurrentIndex()

								dt := ""
								if idx != -1 {
									dt = strconv.Itoa(v.DataTypeCb.Model().([]*utils.KeyValuePair)[idx].Code)
								}

								decimal := v.DecimalNe.Value()

								idx = v.IsAlarmCb.CurrentIndex()
								isAlarm := "0"
								if idx != -1 {
									isAlarm = strconv.Itoa(v.IsAlarmCb.Model().([]*utils.KeyValuePair)[idx].Code)
								}

								idx = v.IsAnalyableCb.CurrentIndex()
								isAnalysable := 0
								if idx != -1 {
									isAnalysable = v.IsAnalyableCb.Model().([]*utils.KeyValuePair)[idx].Code
								}

								can = append(can, &model.CanDetail{
									Id:           v.data.Id,
									OutfieldId:   v.KeyLe.Text(),
									Chinesename:  v.AliasLe.Text(),
									Unit:         v.UnitLe.Text(),
									DataType:     dt,
									Formula:      v.ConvertFormulaLe.Text(),
									DataMap:      v.DataScopeLe.Text(),
									Decimals:     strconv.FormatFloat(decimal, 'f', 0, 64),
									IsAlarm:      isAlarm,
									IsAnalysable: isAnalysable,
									OutfieldSn:   v.SortNe.Value(),
								})
							}

							err := (&dataSource.VehicleTypeEntity{}).BatchUpdateCanDetail(can)
							if err != nil {
								walk.MsgBox(clf.mainWin, "更新失败", err.Error(), walk.MsgBoxIconError)
								return
							}
							mynotify.Info("已保存")
						},
					},
				},
			},
			VSpacer{},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	clf.mainWin = mainWin

	group, err := (&dataSource.VehicleTypeEntity{TypeId: vt.TypeId}).GetVehicleType()
	if err != nil {
		return clf, err
	}

	if len(group.CanGroup) > 0 {
		clf.CanGroupDDL.SetModel(group.CanGroup)
		clf.CanGroupDDL.SetCurrentIndex(0)

		clf.renderCanList()
	}

	return clf, nil
}

type CanListForm struct {
	*walk.Composite
	mainWin *TabMainWindow

	statsLbl *walk.Label

	CanGroupDDL   *walk.ComboBox
	CanListScroll *walk.ScrollView
	btnCom        *walk.Composite

	Table []*Row
}

type Row struct {
	db   *walk.DataBinder
	data *rowData

	rowCmp *walk.Composite

	IdLbl            *walk.Label
	KeyLe            *walk.LineEdit
	RawFieldNameLe   *walk.Label
	AliasLe          *walk.LineEdit
	UnitLe           *walk.LineEdit
	DataTypeCb       *walk.ComboBox
	ConvertFormulaLe *walk.LineEdit
	DataScopeLe      *walk.LineEdit
	DecimalNe        *walk.NumberEdit
	IsAlarmCb        *walk.ComboBox
	IsAnalyableCb    *walk.ComboBox
	SortNe           *walk.NumberEdit
}

type rowData struct {
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

func (clf *CanListForm) renderCanList() {
	ddlGroup := clf.CanGroupDDL.Model().([]*dataSource.CanGroupEntity)
	groupId := ddlGroup[clf.CanGroupDDL.CurrentIndex()].Id
	cans, err := (&dataSource.VehicleTypeEntity{}).ListCan(groupId)
	if err != nil {
		mynotify.Error("CAN字段加载失败：" + err.Error())
		return
	}

	clf.CanListScroll.SetSuspended(true)
	defer clf.CanListScroll.SetSuspended(false)

	c := clf.CanListScroll.Children()
	for i := c.Len() - 1; i >= 0; i-- {
		item := c.At(i)
		item.SetParent(nil)
		item.Dispose()
	}
	clf.Table = nil

	for i := 0; i < len(cans); i++ {
		row := &Row{}

		dt, err := strconv.Atoi(cans[i].DataType)
		if err != nil {
			dt = 4
		}

		prec, err := strconv.Atoi(cans[i].Decimals)
		if err != nil {
			prec = 0
		}

		isAlarm, err := strconv.Atoi(cans[i].IsAlarm)
		if err != nil {
			isAlarm = 0
		}

		row.data = &rowData{
			Id:           cans[i].Id,
			Key:          cans[i].OutfieldId,
			Unit:         cans[i].Unit,
			Sort:         cans[i].OutfieldSn,
			GroupInfoId:  cans[i].GroupInfoId,
			Cname:        cans[i].FieldName,
			Alias:        cans[i].Chinesename,
			Formula:      cans[i].Formula,
			DataType:     dt,
			Prec:         prec,
			DataScope:    cans[i].DataMap,
			IsAlarm:      isAlarm,
			IsAnalysable: cans[i].IsAnalysable,
			IsDelete:     cans[i].IsDelete,
		}

		if strings.TrimSpace(row.data.Cname) == "" {
			row.data.Cname = "未命名"
		}

		isDelete := false
		colorflag := walk.RGB(0x00, 0x00, 0x00)
		if cans[i].IsDelete == 1 {
			isDelete = true
			colorflag = walk.RGB(0xDC, 0x14, 0x3C)
		}

		if err := (Composite{
			AssignTo: &row.rowCmp,
			Layout:   HBox{MarginsZero: true},
			OnMouseDown: func(x, y int, button walk.MouseButton) {
				bgColor := walk.RGB(0x6F, 0xAF, 0x9D)
				scb, _ := SolidColorBrush{bgColor}.Create()
				row.rowCmp.SetBackground(scb)
			},
			OnMouseUp: func(x, y int, button walk.MouseButton) {
				row.rowCmp.SetBackground(walk.NullBrush())
			},
			DataBinder: DataBinder{
				AssignTo:   &row.db,
				DataSource: row.data,
			},
			Children: []Widget{
				HSpacer{},
				Label{
					AssignTo:  &row.IdLbl,
					Name:      cans[i].Id,
					Text:      strconv.Itoa(i + 1),
					MaxSize:   Size{Width: 30},
					MinSize:   Size{Width: 30},
					Font:      Font{PointSize: 10, StrikeOut: isDelete},
					TextColor: colorflag,
				},
				LineEdit{
					AssignTo: &row.KeyLe,
					MaxSize:  Size{Width: 60},
					MinSize:  Size{Width: 60},
					Text:     Bind("Key"),
				},
				Label{
					AssignTo: &row.RawFieldNameLe,
					Text:     Bind("Cname"),
					MaxSize:  Size{Width: 120},
					MinSize:  Size{Width: 120},
				},
				LineEdit{
					AssignTo:  &row.AliasLe,
					Text:      Bind("Alias"),
					MaxSize:   Size{Width: 120},
					MinSize:   Size{Width: 120},
					MaxLength: 100,
				},
				LineEdit{
					AssignTo:  &row.UnitLe,
					Text:      Bind("Unit"), // cans[i].Unit,
					MaxSize:   Size{Width: 50},
					MinSize:   Size{Width: 50},
					MaxLength: 50,
				},
				ComboBox{
					AssignTo:      &row.DataTypeCb,
					BindingMember: "Code",
					DisplayMember: "Name",
					Model:         utils.KnownDataType(),
					Value:         Bind("DataType"),
					MaxSize:       Size{Width: 120},
					MinSize:       Size{Width: 120},
				},
				LineEdit{
					AssignTo:  &row.ConvertFormulaLe,
					Text:      Bind("Formula"), // cans[i].Formula,
					MaxSize:   Size{Width: 80},
					MinSize:   Size{Width: 80},
					MaxLength: 200,
				},
				LineEdit{
					AssignTo:  &row.DataScopeLe,
					Text:      Bind("DataScope"), // cans[i].DataMap,
					MaxSize:   Size{Width: 200},
					MinSize:   Size{Width: 200},
					MaxLength: 1255,
				},
				NumberEdit{
					AssignTo: &row.DecimalNe,
					Value:    Bind("Prec"), //cans[i].Decimals,
					MaxSize:  Size{Width: 50},
					MinSize:  Size{Width: 50},
				},
				ComboBox{
					AssignTo:      &row.IsAlarmCb,
					BindingMember: "Code",
					DisplayMember: "Name",
					Model:         utils.KnownAlarm(),
					Value:         Bind("IsAlarm"),
					MaxSize:       Size{Width: 75},
					MinSize:       Size{Width: 75},
				},
				ComboBox{
					AssignTo:      &row.IsAnalyableCb,
					BindingMember: "Code",
					DisplayMember: "Name",
					Model:         utils.KnownAnaly(),
					Value:         Bind("IsAnalysable"),
					MaxSize:       Size{Width: 85},
					MinSize:       Size{Width: 85},
				},
				NumberEdit{
					AssignTo: &row.SortNe,
					Value:    Bind("Sort"),
					MaxSize:  Size{Width: 50},
					MinSize:  Size{Width: 50},
					Decimals: 2,
				},
				HSpacer{},
			},
		}).Create(NewBuilder(clf.CanListScroll)); err != nil {
			log.Panic(err)
		}

		clf.Table = append(clf.Table, row)
	}

	if len(clf.Table) > 0 {
		clf.statsLbl.SetText(fmt.Sprintf("共 %d 项", len(clf.Table)))
		clf.btnCom.SetVisible(true)
		return
	}

	clf.btnCom.SetVisible(false)
}
