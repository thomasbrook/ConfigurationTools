package client

import (
	"ConfigurationTools/configurationManager"
	"ConfigurationTools/dataSource"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"reflect"
	"runtime"
	"strconv"
)

func NewVehicleTypeAddPage(parent walk.Container, mw *TabMainWindow) (*VehicleTypeAddForm, error) {
	vtaf := &VehicleTypeAddForm{}

	if err := (Composite{
		AssignTo: &vtaf.Composite,
		Layout:   HBox{},
		Children: []Widget{
			HSpacer{},
			TreeView{
				MaxSize:  Size{Width: 320},
				MinSize:  Size{Width: 320},
				AssignTo: &vtaf.treeView,
				//Model:    treeModel,
				OnCurrentItemChanged: func() {
					org := vtaf.treeView.CurrentItem().(*Organization)
					vtaf.orgNameEdit.SetText(org.orgName)
				},
			},
			HSpacer{},
			ScrollView{
				AssignTo:      &vtaf.scroll,
				StretchFactor: 7,
				Layout:        Flow{MarginsZero: true},
				Children: []Widget{
					Composite{
						Layout:  VBox{Margins: Margins{Top: 60}},
						MaxSize: Size{Width: 360},
						MinSize: Size{Width: 360},
						Children: []Widget{
							VSpacer{},
							Composite{
								MaxSize:   Size{Height: 30},
								MinSize:   Size{Height: 30},
								Layout:    HBox{MarginsZero: true},
								Alignment: AlignHNearVCenter,
								Children: []Widget{
									Label{
										Text:    "当前机构",
										MaxSize: Size{Width: 90},
										MinSize: Size{Width: 90},
									},
									LineEdit{
										AssignTo: &vtaf.orgNameEdit,
										Text:     "",
										ReadOnly: true,
									},
									HSpacer{},
								},
							},
							Composite{
								MaxSize:   Size{Height: 30},
								MinSize:   Size{Height: 30},
								Layout:    HBox{MarginsZero: true},
								Alignment: AlignHNearVCenter,
								Children: []Widget{
									Label{
										Text:    "车辆类型",
										MaxSize: Size{Width: 90},
										MinSize: Size{Width: 90},
									},
									LineEdit{
										AssignTo: &vtaf.vehicleTypeNameEdit,
										Text:     "",
									},
									HSpacer{},
								},
							},
							//Composite{
							//	MaxSize:   Size{Height: 30},
							//	MinSize:   Size{Height: 30},
							//	Layout:    HBox{MarginsZero: true},
							//	Alignment: AlignHNearVCenter,
							//	Children: []Widget{
							//		Label{
							//			Text:    "是否为智能机",
							//			MaxSize: Size{Width: 90},
							//			MinSize: Size{Width: 90},
							//		},
							//		RadioButtonGroup{
							//			Buttons: []RadioButton{
							//				{AssignTo: &vtaf.notIntelligentRB, Text: "否", Value: 0},
							//				{AssignTo: &vtaf.intelligentRB, Text: "是", Value: 1},
							//			},
							//		},
							//		HSpacer{},
							//	},
							//},
							Composite{
								MaxSize:   Size{Height: 30},
								MinSize:   Size{Height: 30},
								Layout:    HBox{MarginsZero: true},
								Alignment: AlignHNearVCenter,
								Children: []Widget{
									Label{
										Text:    "过滤不完整行",
										MaxSize: Size{Width: 90},
										MinSize: Size{Width: 90},
									},
									RadioButtonGroup{
										Buttons: []RadioButton{
											{AssignTo: &vtaf.notFilterMissingColumnRB, Text: "不过滤", Value: 0},
											{AssignTo: &vtaf.filterMissingColumnRB, Text: "过滤", Value: 1},
										},
									},
									HSpacer{},
								},
							},
							VSpacer{},
							VSeparator{
								MinSize: Size{Height: 1},
								MaxSize: Size{Height: 1},
							},
							Composite{
								AssignTo:  &vtaf.canGroup,
								Layout:    VBox{MarginsZero: true},
								Alignment: AlignHNearVCenter,
								Children: []Widget{
									Composite{
										Layout:    HBox{},
										Alignment: AlignHCenterVCenter,
										Children: []Widget{
											Label{
												Text:    "CAN类型",
												MaxSize: Size{Width: 80},
												MinSize: Size{Width: 80},
												Font:    Font{PointSize: 10, Bold: true},
											},
											Label{
												Text: "别名",
												Font: Font{PointSize: 10, Bold: true},
											},
											HSpacer{},
											Label{
												Text:    "排序",
												MaxSize: Size{Width: 28},
												MinSize: Size{Width: 28},
												Font:    Font{PointSize: 10, Bold: true},
											},
										},
									},
								},
							},
							VSpacer{},
							Composite{
								Layout:    HBox{},
								Alignment: AlignHNearVCenter,
								Children: []Widget{
									HSpacer{},
									PushButton{
										AssignTo: &vtaf.addVhBtn,
										Text:     "添加",
										//Image:    "/img/add.ico",
										OnClicked: func() {
											vtaf.addVhBtn.SetEnabled(false)

											org := vtaf.treeView.CurrentItem().(*Organization)

											vte := &dataSource.VehicleTypeEntity{}
											vte.OrgId = org.orgId
											vte.TypeName = vtaf.vehicleTypeNameEdit.Text()
											vte.CanGroup = []*dataSource.CanGroupEntity{}

											for i := 0; i < len(vtaf.chb); i++ {
												if vtaf.chb[i].Checked() {
													g := &dataSource.CanGroupEntity{}

													code, err := strconv.Atoi(vtaf.chb[i].Name())
													if err != nil {
														continue
													}

													g.Code = code
													g.Name = vtaf.le[i].Text()
													g.Sort = vtaf.ne[i].Value()

													vte.CanGroup = append(vte.CanGroup, g)
												}
											}

											//if vtaf.intelligentRB.Checked() {
											//	vte.IsIntelligent = 1
											//} else if vtaf.notIntelligentRB.Checked() {
											//	vte.IsIntelligent = 0
											//} else {
											//	vte.IsIntelligent = 0
											//}

											if vtaf.filterMissingColumnRB.Checked() {
												vte.IsFilterMissingColumn = 1
											} else if vtaf.notFilterMissingColumnRB.Checked() {
												vte.IsFilterMissingColumn = 0
											} else {
												vte.IsFilterMissingColumn = 0
											}

											err := vte.Add()
											if err != nil {
												vtaf.addVhBtn.SetEnabled(true)
												walk.MsgBox(vtaf.mainWin, "创建车型失败", err.Error(), walk.MsgBoxIconError)
												return
											}
											vtaf.addVhBtn.SetEnabled(true)
											walk.MsgBox(vtaf.mainWin, "", "创建车型成功", walk.MsgBoxIconInformation)
										},
									},
								},
							},
						},
					},
				},
			},
			HSpacer{},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	vtaf.mainWin = mw
	vtaf.parent = parent

	//vtaf.notIntelligentRB.SetChecked(true)
	vtaf.notFilterMissingColumnRB.SetChecked(true)

	for i := 0; i < len(configurationManager.CanGroup); i++ {

		chk := &walk.CheckBox{}
		le := &walk.LineEdit{}
		ne := &walk.NumberEdit{}

		can := Composite{
			Layout:    HBox{MarginsZero: true},
			MaxSize:   Size{Height: 30},
			MinSize:   Size{Height: 30},
			Alignment: AlignHNearVCenter,
			Children: []Widget{
				CheckBox{
					AssignTo: &chk,
					Name:     strconv.Itoa(configurationManager.CanGroup[i].Code),
					Text:     configurationManager.CanGroup[i].GroupName,

					MaxSize: Size{Width: 90},
					MinSize: Size{Width: 90},
				},
				LineEdit{
					AssignTo: &le,
					Name:     strconv.Itoa(configurationManager.CanGroup[i].Code),
					Text:     configurationManager.CanGroup[i].GroupName,
				},
				NumberEdit{
					AssignTo: &ne,
					Name:     strconv.Itoa(configurationManager.CanGroup[i].Code),
					Value:    float64(i),
					Decimals: 0,
				},
				HSpacer{},
			},
		}

		err := can.Create(NewBuilder(vtaf.canGroup))
		if err != nil {
			panic(err)
		}

		vtaf.chb = append(vtaf.chb, chk)
		vtaf.le = append(vtaf.le, le)
		vtaf.ne = append(vtaf.ne, ne)
	}

	go func() {
		treeModel, err := NewOrganizationTreeModel()
		if err != nil {
			log.Fatal(err)
		}
		vtaf.treeView.SetModel(treeModel)
	}()

	return vtaf, nil
}

type VehicleTypeAddForm struct {
	*walk.Composite

	parent  walk.Container
	mainWin *TabMainWindow

	scroll *walk.ScrollView

	treeView            *walk.TreeView
	orgNameEdit         *walk.LineEdit
	vehicleTypeNameEdit *walk.LineEdit
	//notIntelligentRB         *walk.RadioButton
	//intelligentRB            *walk.RadioButton
	notFilterMissingColumnRB *walk.RadioButton
	filterMissingColumnRB    *walk.RadioButton

	canGroup *walk.Composite

	chb []*walk.CheckBox
	le  []*walk.LineEdit
	ne  []*walk.NumberEdit

	addVhBtn *walk.PushButton
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
