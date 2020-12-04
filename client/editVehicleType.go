package client

import (
	"ConfigurationTools/configurationManager"
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"encoding/json"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"strconv"
)

func NewVehicleTypeEditPage(parent walk.Container, vt *model.VehicleTypeStats, mw *TabMainWindow) (*VehicleTypeForm, error) {

	var orgLbl *walk.Label
	vtf := &VehicleTypeForm{}
	vtf.canGroupMap = make(map[int]dataSource.CanGroupEntity)

	var db *walk.DataBinder

	if err := (Composite{
		AssignTo:  &vtf.Composite,
		Layout:    HBox{},
		Alignment: AlignHCenterVCenter,
		Children: []Widget{
			ScrollView{
				Layout: Flow{MarginsZero: true},
				Children: []Widget{
					Composite{
						Layout: VBox{Margins: Margins{Top: 60}},
						DataBinder: DataBinder{
							AssignTo:   &db,
							DataSource: vt,
						},
						MaxSize: Size{Width: 360},
						MinSize: Size{Width: 360},
						Children: []Widget{
							Composite{
								MaxSize:   Size{Height: 30},
								MinSize:   Size{Height: 30},
								Layout:    HBox{MarginsZero: true},
								Alignment: AlignHNearVCenter,
								Children: []Widget{
									Label{
										Text:    "所属机构",
										MaxSize: Size{Width: 90},
										MinSize: Size{Width: 90},
									},
									Label{
										AssignTo: &orgLbl,
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
										Text:    "车系名称",
										MaxSize: Size{Width: 90},
										MinSize: Size{Width: 90},
									},
									LineEdit{
										AssignTo: &vtf.vtnamele,
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
							//			Text:        "是否为智能机",
							//			ToolTipText: "可自动计亩的设备",
							//			MaxSize:     Size{Width: 90},
							//			MinSize:     Size{Width: 90},
							//		},
							//		RadioButtonGroup{
							//			Buttons: []RadioButton{
							//				{AssignTo: &vtf.notIntelligentRB, Text: "否", Value: 0},
							//				{AssignTo: &vtf.intelligentRB, Text: "是", Value: 1},
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
										Text:        "过滤不完整行",
										ToolTipText: "包含空字段的数据行",
										MaxSize:     Size{Width: 90},
										MinSize:     Size{Width: 90},
									},
									RadioButtonGroup{
										Buttons: []RadioButton{
											{AssignTo: &vtf.notFilterMissingColumnRB, Text: "不过滤", Value: 0},
											{AssignTo: &vtf.filterMissingColumnRB, Text: "过滤", Value: 1},
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
								AssignTo:  &vtf.canGroupCom,
								Layout:    VBox{MarginsZero: true},
								Alignment: AlignHNearVCenter,
								Children: []Widget{
									Composite{
										Layout:    HBox{MarginsZero: true, Margins: Margins{Top: 3, Right: 8, Bottom: 3, Left: 3}},
										Alignment: AlignHCenterVCenter,
										Children: []Widget{
											Label{
												Text:    "分组类型",
												MaxSize: Size{Width: 90},
												MinSize: Size{Width: 90},
												Font:    Font{PointSize: 10, Bold: true},
											},
											Label{
												Text: "分组别名",
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
							HSpacer{},
							Composite{
								Layout:    HBox{},
								Alignment: AlignHNearVCenter,
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text: "复制",
										//Image: "/img/copy.ico",
										OnClicked: func() {
											vtype, err := (&dataSource.VehicleTypeEntity{TypeId: vt.TypeId}).GetVehicleType()
											if err != nil {
												mynotify.Error(err.Error())
												return
											}

											jsonData, err := json.MarshalIndent(vtype, "", "\t")
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
									PushButton{
										AssignTo: &vtf.updateBtn,
										//Image:    "/img/update.ico",
										Text: "更新",
										OnClicked: func() {
											vtf.updateBtn.SetEnabled(false)

											vte := &dataSource.VehicleTypeEntity{}
											vte.CanGroup = []*dataSource.CanGroupEntity{}

											for i := 0; i < len(vtf.chb); i++ {
												if vtf.chb[i].Checked() {
													g := &dataSource.CanGroupEntity{}

													code, err := strconv.Atoi(vtf.chb[i].Name())
													if err != nil {
														continue
													}

													if val, isExist := vtf.canGroupMap[code]; isExist {
														g.Id = val.Id
													}

													g.Code = code
													g.Name = vtf.le[i].Text()
													g.Sort = vtf.ne[i].Value()

													vte.CanGroup = append(vte.CanGroup, g)
												}
											}

											vte.TypeId = vt.TypeId
											vte.TypeName = vtf.vtnamele.Text()

											//if vtf.intelligentRB.Checked() {
											//	vte.IsIntelligent = 1
											//} else if vtf.notIntelligentRB.Checked() {
											//	vte.IsIntelligent = 0
											//} else {
											//	vte.IsIntelligent = 0
											//}

											if vtf.filterMissingColumnRB.Checked() {
												vte.IsFilterMissingColumn = 1
											} else if vtf.notFilterMissingColumnRB.Checked() {
												vte.IsFilterMissingColumn = 0
											} else {
												vte.IsFilterMissingColumn = 0
											}

											err := vte.Update()
											if err != nil {
												vtf.updateBtn.SetEnabled(true)
												walk.MsgBox(vtf.mainWin, "更新失败", err.Error(), walk.MsgBoxIconError)
												return
											}
											vtf.updateBtn.SetEnabled(true)
											walk.MsgBox(vtf.mainWin, "", "更新成功", walk.MsgBoxIconInformation)
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

	vtf.mainWin = mw
	vtf.parent = parent

	vtype, err := (&dataSource.VehicleTypeEntity{TypeId: vt.TypeId}).GetVehicleType()
	if err != nil {
		return vtf, err
	}

	orgLbl.SetText(vtype.OrgName)
	vtf.vtnamele.SetText(vtype.TypeName)

	//if vtype.IsIntelligent == 1 {
	//	vtf.intelligentRB.SetChecked(true)
	//} else {
	//	vtf.notIntelligentRB.SetChecked(true)
	//}

	if vtype.IsFilterMissingColumn == 1 {
		vtf.filterMissingColumnRB.SetChecked(true)
	} else {
		vtf.notFilterMissingColumnRB.SetChecked(true)
	}

	for _, val := range vtype.CanGroup {
		vtf.canGroupMap[val.Code] = *val
	}

	for i := 0; i < len(configurationManager.CanGroup); i++ {

		group, isExist := vtf.canGroupMap[configurationManager.CanGroup[i].Code]

		gname := configurationManager.CanGroup[i].GroupName
		sort := float64(i)
		if isExist {
			gname = group.Name
			sort = group.Sort
		}

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
					Checked:  isExist,
					Enabled:  !isExist,
					MaxSize:  Size{Width: 90},
					MinSize:  Size{Width: 90},
				},
				LineEdit{
					AssignTo: &le,
					Name:     strconv.Itoa(configurationManager.CanGroup[i].Code),
					Text:     gname,
				},
				NumberEdit{
					AssignTo: &ne,
					Name:     strconv.Itoa(configurationManager.CanGroup[i].Code),
					Value:    sort,
					Decimals: 0,
				},
				HSpacer{},
			},
		}
		err := can.Create(NewBuilder(vtf.canGroupCom))
		if err != nil {
			panic(err)
		}

		vtf.chb = append(vtf.chb, chk)
		vtf.le = append(vtf.le, le)
		vtf.ne = append(vtf.ne, ne)
	}

	return vtf, nil
}

type VehicleTypeForm struct {
	*walk.Composite

	parent  walk.Container
	mainWin *TabMainWindow

	vtnamele *walk.LineEdit
	//notIntelligentRB         *walk.RadioButton
	//intelligentRB            *walk.RadioButton
	notFilterMissingColumnRB *walk.RadioButton
	filterMissingColumnRB    *walk.RadioButton

	canGroupCom *walk.Composite

	canGroupMap map[int]dataSource.CanGroupEntity
	chb         []*walk.CheckBox
	le          []*walk.LineEdit
	ne          []*walk.NumberEdit

	updateBtn *walk.PushButton
}
