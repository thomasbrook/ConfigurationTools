package client

import (
	"ConfigurationTools/model"
	"fmt"
	"log"
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type editForm struct {
	SearchStr string

	form *walk.Composite

	outfieldIdEdit *walk.LineEdit
	fieldNameEdit  *walk.LineEdit
	cnEdit         *walk.LineEdit
	unitEdit       *walk.LineEdit
	formulaEdit    *walk.LineEdit
	dataMapEdit    *walk.LineEdit
	decimalEdit    *walk.LineEdit
	groupIdEdit    *walk.LineEdit

	outfieldIdSnEdit *walk.NumberEdit

	dateValRb     *walk.RadioButton
	enumValRb     *walk.RadioButton
	numberValRb   *walk.RadioButton
	otherValRb    *walk.RadioButton
	enumTextValRb *walk.RadioButton

	alarmRb   *walk.RadioButton
	noAlarmRb *walk.RadioButton

	analyableRb         *walk.RadioButton
	noAnalyzableRb      *walk.RadioButton
	defaultAnalyzableRb *walk.RadioButton

	dataTypeRadioGroup *walk.GroupBox
}

func (ef *editForm) runExportDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		Title:         "查询",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		DataBinder: DataBinder{
			AssignTo:   &db,
			DataSource: ef,
		},
		MinSize: Size{300, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 1},
				Children: []Widget{
					Label{
						Text: "输入搜索关键词，多个换行显示",
					},
					HSpacer{
						Size: 5,
					},
					TextEdit{
						MinSize: Size{100, 50},
						Text:    Bind("SearchStr"),
						VScroll: true,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}

func (ef *editForm) canDetail(can *model.CanConfig) {
	ef.outfieldIdEdit.SetText(can.OutfieldId)
	ef.unitEdit.SetText(can.Unit)
	ef.cnEdit.SetText(can.Chinesename)
	ef.formulaEdit.SetText(can.Formula)
	ef.fieldNameEdit.SetText(can.FieldName)
	ef.decimalEdit.SetText(can.Decimals)
	ef.dataMapEdit.SetText(can.DataMap)
	ef.outfieldIdSnEdit.SetValue(can.OutfieldSn)

	var dt model.DataType
	b, err := strconv.Atoi(can.DataType)
	if err != nil {
		dt = model.NullValue
	} else {
		dt = model.DataType(b)
	}

	ef.dateValRb.SetChecked(false)
	ef.enumValRb.SetChecked(false)
	ef.numberValRb.SetChecked(false)
	ef.otherValRb.SetChecked(false)
	ef.enumTextValRb.SetChecked(false)

	switch dt {
	case model.DateValue:
		ef.dateValRb.SetChecked(true)
	case model.EnumValue:
		ef.enumValRb.SetChecked(true)
	case model.NumericValue:
		ef.numberValRb.SetChecked(true)
	case model.OtherValue:
		ef.otherValRb.SetChecked(true)
	case model.EnumText:
		ef.enumTextValRb.SetChecked(true)
	}

	ef.alarmRb.SetChecked(false)
	ef.noAlarmRb.SetChecked(false)
	switch can.IsAlarm {
	case "1":
		ef.alarmRb.SetChecked(true)
	case "0":
		ef.noAlarmRb.SetChecked(true)
	}

	ef.analyableRb.SetChecked(false)
	ef.defaultAnalyzableRb.SetChecked(false)
	ef.noAnalyzableRb.SetChecked(false)
	switch can.IsAnalysable {
	case 1:
		ef.analyableRb.SetChecked(true)
	case 2:
		ef.defaultAnalyzableRb.SetChecked(true)
	case 0:
		ef.noAnalyzableRb.SetChecked(true)
	}
}

func (ef *editForm) set() (can *model.CanConfig, err error) {
	can = new(model.CanConfig)
	can.OutfieldId = ef.outfieldIdEdit.Text()
	can.FieldName = ef.fieldNameEdit.Text()
	can.Chinesename = ef.cnEdit.Text()
	can.Unit = ef.unitEdit.Text()

	if ef.dateValRb.Checked() {
		can.DataType = strconv.Itoa(ef.dateValRb.Value().(int))
	}
	if ef.enumValRb.Checked() {
		can.DataType = strconv.Itoa(ef.enumValRb.Value().(int))
	}
	if ef.numberValRb.Checked() {
		can.DataType = strconv.Itoa(ef.numberValRb.Value().(int))
	}
	if ef.enumTextValRb.Checked() {
		can.DataType = strconv.Itoa(ef.enumTextValRb.Value().(int))
	}
	if ef.otherValRb.Checked() {
		can.DataType = strconv.Itoa(ef.otherValRb.Value().(int))
	}

	can.Formula = ef.formulaEdit.Text()
	can.DataMap = ef.dataMapEdit.Text()
	can.Decimals = ef.decimalEdit.Text()

	if ef.alarmRb.Checked() {
		can.IsAlarm = strconv.Itoa(ef.alarmRb.Value().(int))
	}
	if ef.noAlarmRb.Checked() {
		can.IsAlarm = strconv.Itoa(ef.alarmRb.Value().(int))
	}

	if ef.analyableRb.Checked() {
		can.IsAnalysable = ef.analyableRb.Value().(int)
	}
	if ef.defaultAnalyzableRb.Checked() {
		can.IsAnalysable = ef.defaultAnalyzableRb.Value().(int)
	}
	if ef.noAnalyzableRb.Checked() {
		can.IsAnalysable = ef.noAnalyzableRb.Value().(int)
	}

	can.OutfieldSn = ef.outfieldIdSnEdit.Value()

	return can, nil
}

func (ef *editForm) submit() (isOk bool, err error) {

	log.Print(fmt.Sprintf("%+v", ef.outfieldIdEdit.Text()))
	log.Print(fmt.Sprintf("%+v", ef.unitEdit.Text()))
	log.Print(fmt.Sprintf("%+v", ef.cnEdit.Text()))
	log.Print(fmt.Sprintf("%+v", ef.formulaEdit.Text()))
	log.Print(fmt.Sprintf("%+v", ef.fieldNameEdit.Text()))
	log.Print(fmt.Sprintf("%+v", ef.decimalEdit.Text()))
	log.Print(fmt.Sprintf("%+v", ef.dataMapEdit.Text()))
	log.Print(fmt.Sprintf("%+v", ef.outfieldIdSnEdit.Value()))

	log.Print(fmt.Sprintf("%+v", ef.dateValRb.Checked()))
	log.Print(fmt.Sprintf("%+v", ef.enumValRb.Checked()))
	log.Print(fmt.Sprintf("%+v", ef.numberValRb.Checked()))
	log.Print(fmt.Sprintf("%+v", ef.otherValRb.Checked()))

	return true, nil
}
