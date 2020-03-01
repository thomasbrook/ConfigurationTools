package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func NewVehicleTypeList(parent walk.Container, mw *TabMainWindow) (*VehicleTypeCom, error) {
	vtc, err := initList(parent)
	if err != nil {
		return nil, err
	}
	vtc.mainWin = mw
	vtc.parent = parent

	err = vtc.loadVehicleType()
	if err != nil {
		return nil, err
	}
	return vtc, nil
}

func initList(parent walk.Container) (*VehicleTypeCom, error) {
	vtc := &VehicleTypeCom{}

	if err := (Composite{
		AssignTo: &vtc.Composite,
		Layout: VBox{
			MarginsZero: true,
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
						AssignTo:   &vtc.addVehicleType,
						Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						Image:      "img/add.png",
						Margin:     2,
					},
					HSpacer{},
					Composite{
						Layout: HBox{MarginsZero: true, Margins: Margins{Right: 5}},
						Children: []Widget{
							LineEdit{
								AssignTo:  &vtc.searchKeyEdit,
								Alignment: AlignHFarVCenter,
								MaxSize:   Size{Width: 180},
							},
						},
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							PushButton{
								Text:       "查询",
								Background: SolidColorBrush{Color: walk.RGB(0xFF, 0xFF, 0xFF)},
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
					MarginsZero: true,
					Alignment:   AlignHNearVNear,
				},
				Children: []Widget{},
			},
			Composite{
				Layout: HBox{
					Margins:     Margins{Top: 5},
					MarginsZero: true,
					SpacingZero: true,
				},
				Children: []Widget{
					Label{
						AssignTo: &vtc.statLbl,
						Text:     "",
					},
					HSpacer{},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return nil, err
	}

	vtc.addVehicleType.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			vtc.AddVehicleType()
		}
	})

	return vtc, nil
}

type VehicleTypeCom struct {
	*walk.Composite
	parent         walk.Container
	mainWin        *TabMainWindow
	vehicleSeries  *walk.ScrollView
	addVehicleType *walk.ImageView
	statLbl        *walk.Label

	searchKeyEdit *walk.LineEdit
}

func (vtc *VehicleTypeCom) EditVehicleType(data *model.VehicleType) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder

	return Dialog{
		AssignTo:   &dlg,
		Title:      "编辑 车辆类型",
		MinSize:    Size{420, 360},
		Background: SolidColorBrush{Color: walk.RGB(0xFF, 0xFF, 0xFF)},
		Layout:     VBox{},
		DataBinder: DataBinder{
			AssignTo:   &db,
			DataSource: data,
		},
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					Composite{
						StretchFactor: 7,
						Layout:        VBox{MarginsZero: true},
						Children: []Widget{
							HSpacer{},
							Composite{
								Layout:    HBox{},
								Alignment: AlignHCenterVNear,
								Children: []Widget{
									HSpacer{},
									Label{
										Text: "当前机构",
									},
									LineEdit{
										Text:     Bind("OrgName"),
										MaxSize:  Size{200, 30},
										ReadOnly: true,
									},
									HSpacer{},
								},
							},
							Composite{
								Layout:    HBox{},
								Alignment: AlignHCenterVNear,
								Children: []Widget{
									HSpacer{},
									Label{
										Text: "车辆类型",
									},
									LineEdit{
										Text:    Bind("TypeName"),
										MaxSize: Size{200, 30},
									},
									HSpacer{},
								},
							},
							Composite{
								Layout:    HBox{},
								Alignment: AlignHCenterVNear,
								Children: []Widget{
									HSpacer{},
									CheckBox{
										Text:    "常规信息",
										Checked: Bind("HasGeneral"),
										MaxSize: Size{200, 30},
										Enabled: !data.HasGeneral,
									},
									CheckBox{
										Text:    "CAN信息",
										Checked: Bind("HasCan"),
										MaxSize: Size{200, 30},
										Enabled: !data.HasCan,
									},
									CheckBox{
										Text:    "指令下发",
										Checked: Bind("HasCmd"),
										MaxSize: Size{200, 30},
										Enabled: !data.HasCmd,
									},
									HSpacer{},
								},
							},
							HSpacer{},
						},
					},
					Composite{
						StretchFactor: 3,
						Layout:        VBox{MarginsZero: true},
						Children: []Widget{
							Composite{
								Layout:    HBox{Margins: Margins{0, 0, 55, 0}},
								Alignment: AlignHCenterVCenter,
								Children: []Widget{
									HSpacer{},
									PushButton{
										Text:    "更新",
										MaxSize: Size{60, 60},
										OnClicked: func() {
											if err := db.Submit(); err != nil {
												log.Print(err)
												return
											}
											ds := db.DataSource().(*model.VehicleType)
											log.Print(fmt.Sprintf("%+v", ds))

											// 已存在的分组，暂时不允许删除。只允许创建不存在的分组
											isCreateGeneral := ds.HasGeneral
											if strings.TrimSpace(ds.GeneralId) != "" {
												isCreateGeneral = false
											}

											isCreateCan := ds.HasCan
											if strings.TrimSpace(ds.CanId) != "" {
												isCreateCan = false
											}

											isCreateCmd := ds.HasCmd
											if strings.TrimSpace(ds.CmdId) != "" {
												isCreateCmd = false
											}

											err := dataSource.EditVehicleType(ds.TypeId, ds.TypeName, isCreateGeneral, isCreateCan, isCreateCmd)
											if err != nil {
												log.Fatal(err)
												return
											}

											dlg.Accept()
											vtc.loadVehicleType()
										},
									},
									PushButton{
										Text:    "取消",
										MaxSize: Size{60, 60},
										OnClicked: func() {
											dlg.Cancel()
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}.Run(vtc.mainWin)
}

func (vtc *VehicleTypeCom) AddVehicleType() (int, error) {
	var dlg *walk.Dialog
	var treeView *walk.TreeView
	var orgNameEdit, vehicleTypeNameEdit *walk.LineEdit
	var generalCb, canCb, cmdCd *walk.CheckBox

	treeModel, err := NewOrganizationTreeModel()
	if err != nil {
		log.Fatal(err)
	}

	return Dialog{
		AssignTo:   &dlg,
		Title:      "创建 车辆类型",
		MinSize:    Size{700, 450},
		Background: SolidColorBrush{Color: walk.RGB(0xFF, 0xFF, 0xFF)},
		Layout:     VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TreeView{
						AssignTo:      &treeView,
						StretchFactor: 5,
						Model:         treeModel,
						OnCurrentItemChanged: func() {
							org := treeView.CurrentItem().(*Organization)
							orgNameEdit.SetText(org.orgName)
						},
					},
					Composite{
						StretchFactor: 6,
						Layout:        VBox{MarginsZero: true},
						Children: []Widget{
							VSplitter{
								Children: []Widget{
									Composite{
										StretchFactor: 7,
										Layout:        VBox{MarginsZero: true},
										Children: []Widget{
											HSpacer{},
											Composite{
												Layout:    HBox{},
												Alignment: AlignHCenterVNear,
												Children: []Widget{
													HSpacer{},
													Label{
														Text: "当前机构",
													},
													LineEdit{
														AssignTo: &orgNameEdit,
														Text:     "",
														MaxSize:  Size{200, 30},
														ReadOnly: true,
													},
													HSpacer{},
												},
											},
											Composite{
												Layout:    HBox{},
												Alignment: AlignHCenterVNear,
												Children: []Widget{
													HSpacer{},
													Label{
														Text: "车辆类型",
													},
													LineEdit{
														AssignTo: &vehicleTypeNameEdit,

														Text:    "",
														MaxSize: Size{200, 30},
													},
													HSpacer{},
												},
											},
											Composite{
												Layout:    HBox{},
												Alignment: AlignHCenterVNear,
												Children: []Widget{
													HSpacer{},
													CheckBox{
														AssignTo: &generalCb,
														Text:     "常规信息",
														Checked:  true,
														MaxSize:  Size{200, 30},
													},
													CheckBox{
														AssignTo: &canCb,
														Text:     "CAN信息",
														Checked:  true,
														MaxSize:  Size{200, 30},
													},
													CheckBox{
														AssignTo: &cmdCd,
														Text:     "指令下发",
														Checked:  true,
														MaxSize:  Size{200, 30},
													},
													HSpacer{},
												},
											},
											HSpacer{},
										},
									},
									Composite{
										StretchFactor: 3,
										Layout:        VBox{MarginsZero: true},
										Children: []Widget{
											Composite{
												Layout:    HBox{Margins: Margins{0, 0, 55, 0}},
												Alignment: AlignHCenterVCenter,
												Children: []Widget{
													HSpacer{},
													PushButton{
														Text:    "新建",
														MaxSize: Size{60, 60},
														OnClicked: func() {
															org := treeView.CurrentItem().(*Organization)
															err := dataSource.AddVehicleType(org.orgId, vehicleTypeNameEdit.Text(), generalCb.Checked(), canCb.Checked(), cmdCd.Checked())
															if err != nil {
																log.Fatal(err)
																return
															}

															dlg.Cancel()

															vtc.loadVehicleType()
														},
													},
													PushButton{
														Text:    "取消",
														MaxSize: Size{60, 60},
														OnClicked: func() {
															dlg.Cancel()
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}.Run(vtc.mainWin)
}

func (vtc *VehicleTypeCom) loadVehicleType() error {
	vtc.vehicleSeries.SetSuspended(true)
	defer vtc.vehicleSeries.SetSuspended(false)

	c := vtc.vehicleSeries.Children()
	for i := c.Len() - 1; i >= 0; i-- {
		item := c.At(i)
		item.SetParent(nil)
		item.Dispose()
	}

	//	vtc.vehicleSeries.Children().Clear()

	searchKey := vtc.searchKeyEdit.Text()
	vts, err := dataSource.ListVehicleType(searchKey)
	if err != nil {
		return err
	}

	for i := 0; i < len(vts); i++ {

		s := new(VehicleTypeItem)

		var db *walk.DataBinder
		data := &model.VehicleType{
			TypeId:           vts[i].TypeId,
			TypeName:         vts[i].TypeName,
			OrgName:          vts[i].OrgName,
			GeneralId:        vts[i].GeneralId,
			GeneralInfoCount: vts[i].GeneralInfoCount,
			CanId:            vts[i].CanId,
			CanCount:         vts[i].CanCount,
			CmdId:            vts[i].CmdId,
			CmdCount:         vts[i].CmdCount,
			HasGeneral:       strings.TrimSpace(vts[i].GeneralId) != "",
			HasCan:           strings.TrimSpace(vts[i].CanId) != "",
			HasCmd:           strings.TrimSpace(vts[i].CmdId) != "",
		}

		w := Composite{
			AssignTo: &s.Composite,
			DataBinder: DataBinder{
				AssignTo:   &db,
				DataSource: data,
			},
			ContextMenuItems: []MenuItem{
				Action{
					Text: "Edit Can",
					OnTriggered: func() {
						vtc.mainWin.EditGroupCan(db.DataSource().(*model.VehicleType))
					},
				},
				Action{
					Text: "Add Can",
					OnTriggered: func() {
						vtc.mainWin.AddGroupCan(db.DataSource().(*model.VehicleType))
					},
				},
				Action{
					Text: "Edit 车系",
					OnTriggered: func() {
						vtc.EditVehicleType(db.DataSource().(*model.VehicleType))
					},
				},
			},
			Layout:     VBox{},
			Background: SolidColorBrush{Color: walk.RGB(0x00, 0x8B, 0x8B)},
			MaxSize:    Size{220, 160},
			MinSize:    Size{220, 160},
			Children: []Widget{
				Label{
					AssignTo:      &s.VehicleTypeNameLbl,
					StretchFactor: 2,
					Alignment:     AlignHNearVCenter,
					Font:          Font{PointSize: 11},
					Text:          Bind("TypeName"),
					TextColor:     walk.RGB(0xFF, 0xFF, 0xFF),
				},
				Composite{
					AssignTo:      &s.groupCmp,
					Layout:        VBox{MarginsZero: true, SpacingZero: true},
					StretchFactor: 6,
					Children: []Widget{
						Composite{
							AssignTo: &s.GeneralCmp,
							Visible:  Bind("HasGeneral"),
							Layout:   HBox{MarginsZero: true, SpacingZero: true},
							Children: []Widget{
								Label{
									Text:      "常规信息",
									Alignment: AlignHNearVCenter,
									Font:      Font{PointSize: 9},
									TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
								},
								Label{
									AssignTo:  &s.GeneralLbl,
									Alignment: AlignHNearVCenter,
									Font:      Font{PointSize: 9},
									Text:      Bind("GeneralInfoCount"),
									TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
								},
							},
						},
						Composite{
							AssignTo: &s.CanCmp,
							Visible:  Bind("HasCan"),
							Layout:   HBox{MarginsZero: true, SpacingZero: true},
							Children: []Widget{
								Label{
									Text:      "CAN信息",
									Alignment: AlignHNearVCenter,
									Font:      Font{PointSize: 9},
									TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
								},
								Label{
									AssignTo:  &s.CanLbl,
									Alignment: AlignHNearVCenter,
									Font:      Font{PointSize: 9},
									Text:      Bind("CanCount"),
									TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
								},
							},
						},
						Composite{
							AssignTo: &s.IssueCmp,
							Visible:  Bind("HasCan"),
							Layout:   HBox{MarginsZero: true, SpacingZero: true},
							Children: []Widget{
								Label{
									Text:      "指令下发",
									Alignment: AlignHNearVCenter,
									Font:      Font{PointSize: 9},
									TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
								},
								Label{
									AssignTo:  &s.IssueLbl,
									Alignment: AlignHNearVCenter,
									Font:      Font{PointSize: 9},
									Text:      Bind("CmdCount"),
									TextColor: walk.RGB(0xFF, 0xFF, 0xFF),
								},
							},
						},
					},
				},
				Label{
					AssignTo:      &s.OrgNameLbl,
					StretchFactor: 2,
					Alignment:     AlignHFarVCenter,
					Font:          Font{PointSize: 10},
					Text:          Bind("OrgName"),
					TextColor:     walk.RGB(0xFF, 0xFF, 0xFF),
				},
			},
		}

		err := w.Create(NewBuilder(vtc.vehicleSeries))
		if err != nil {
			panic(err)
		}

		s.GeneralCmp.SetVisible(data.HasGeneral)
		s.CanCmp.SetVisible(data.HasCan)
		s.IssueCmp.SetVisible(data.HasCmd)

		s.GeneralLbl.SetText(strconv.Itoa(data.GeneralInfoCount))
		s.CanLbl.SetText(strconv.Itoa(data.CanCount))
		s.IssueLbl.SetText(strconv.Itoa(data.CmdCount))

		//		s.Composite.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		//			if button == walk.LeftButton {
		//				ds := s.Composite.DataBinder().DataSource()
		//				vtc.MainWin.EditVehicle(ds.(*model.VehicleType))
		//			}
		//		})
	}

	vtc.statLbl.SetText(strconv.Itoa(len(vts)) + " 个")
	return nil
}

type VehicleTypeItem struct {
	*walk.Composite

	VehicleTypeNameLbl *walk.Label
	OrgNameLbl         *walk.Label

	groupCmp *walk.Composite

	GeneralCmp *walk.Composite
	GeneralLbl *walk.Label

	CanCmp *walk.Composite
	CanLbl *walk.Label

	IssueCmp *walk.Composite
	IssueLbl *walk.Label
}
