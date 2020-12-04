package client

import (
	"ConfigurationTools/configurationManager"
	"ConfigurationTools/dataSource"
	"bytes"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func NewSendInstructionsPanel(parent walk.Container, mainWin *TabMainWindow) (*SendInstructionsPage, error) {
	sip := &SendInstructionsPage{
		mainWin: mainWin,
	}

	terType, err := new(dataSource.TerminalTypeEntity).ListTerminalType()
	if err != nil {
		walk.MsgBox(mainWin, "", err.Error(), walk.MsgBoxIconError)
		terType = []*dataSource.TerminalTypeEntity{}
	}

	now := time.Now()
	startDate := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if err := (Composite{
		AssignTo: &sip.Composite,
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Composite{
						StretchFactor: 5,
						Layout:        VBox{MarginsZero: true, Margins: Margins{Right: 5}},
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{Text: "选择终端类型："},
									ComboBox{
										AssignTo:      &sip.terTypeCb,
										BindingMember: "Id",
										DisplayMember: "TypeName",
										MaxSize:       Size{Width: 120},
										Model:         terType,
										ToolTipText:   "请选择终端类型",
										OnCurrentIndexChanged: func() {
											currIdx := sip.terTypeCb.CurrentIndex()
											if currIdx == -1 {
												return
											}

											_ = sip.terTypeCb.Model().([]*dataSource.TerminalTypeEntity)
										},
									},
									HSpacer{},
									Label{Text: "请输入终端号："},
									LineEdit{
										AssignTo:    &sip.ternoLe,
										ToolTipText: "请输入终端编号",
										MaxSize:     Size{Width: 120},
									},
								},
							},
							ScrollView{
								AssignTo: &sip.instructionsContanier,
								Layout:   VBox{MarginsZero: true},
								Children: []Widget{},
							},
						},
					},
					Composite{
						StretchFactor: 5,
						Layout:        VBox{MarginsZero: true, Margins: Margins{Left: 5}},
						Children: []Widget{
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									Label{Text: "下发历史："},
									DateEdit{
										MaxSize:     Size{Width: 78},
										Date:        startDate,
										Format:      "yyyy/MM/dd",
										ToolTipText: "请选择查询开始日期",
									},
									Label{Text: "至"},
									DateEdit{
										MaxSize:     Size{Width: 78},
										Date:        endDate,
										Format:      "yyyy/MM/dd",
										ToolTipText: "请选择查询结束日期",
									},
									HSpacer{},
									LineEdit{
										Text:        "",
										MaxSize:     Size{Width: 120},
										ToolTipText: "请输入完整的终端编号",
									},
									PushButton{
										Text: "查询",
										OnClicked: func() {

										},
									},
								},
							},
							TableView{
								AlternatingRowBG: true,
								//AlternatingRowBGColor: walk.RGB(239, 239, 239),
								ColumnsOrderable: true,
								CheckBoxes:       true,
								MultiSelection:   true,
								Columns: []TableViewColumn{
									{Title: "#", Frozen: true, Alignment: AlignCenter, Width: 35},
									{Title: "终端编号", Width: 120, Alignment: AlignFar},
									{Title: "指令内容", Alignment: AlignCenter, Width: 200},
									{Title: "执行时间", Alignment: AlignCenter, Width: 80},
									{Title: "执行状态", Alignment: AlignCenter, Width: 80},
								},
								OnSelectedIndexesChanged: func() {

								},
								StyleCell: func(style *walk.CellStyle) {

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

	sip.renderForm()

	return sip, nil
}

type SendInstructionsPage struct {
	*walk.Composite
	mainWin *TabMainWindow

	terTypeCb             *walk.ComboBox
	ternoLe               *walk.LineEdit
	instructionsContanier *walk.ScrollView

	titleHearderCb []*walk.CheckBox
	paramLe        []*walk.LineEdit
	cmdCmp         []*walk.Composite
}

func (sip *SendInstructionsPage) renderForm() {
	cmd, err := new(dataSource.CommandConfigEntity).ListCommand()
	if err != nil {
		walk.MsgBox(sip.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}

	sip.instructionsContanier.SetSuspended(true)
	defer sip.instructionsContanier.SetSuspended(false)

	sip.titleHearderCb = nil
	sip.paramLe = nil
	sip.cmdCmp = nil

	host := configurationManager.AppSetting("CmdIp")

	for idx, item := range cmd {

		title := fmt.Sprintf("%d、%s", idx+1, item.CmdName)

		var bodyCmp *walk.Composite
		var titleCb *walk.CheckBox
		var urlLbl *walk.Label
		if err := (Composite{
			Layout: VBox{Margins: Margins{2, 2, 2, 2}},
			Children: []Widget{
				Composite{
					Layout: HBox{MarginsZero: true},
					Children: []Widget{
						CheckBox{
							Text:     title,
							AssignTo: &titleCb,
							Name:     fmt.Sprintf("%d", idx),
							OnClicked: func() {
								bodyCmp.SetVisible(titleCb.Checked())

								for i, _ := range sip.titleHearderCb {
									if fmt.Sprintf("%d", i) != titleCb.Name() {
										sip.titleHearderCb[i].SetChecked(false)
										sip.cmdCmp[i].SetVisible(false)
									}
								}
							},
						},
						HSpacer{},
						Label{AssignTo: &urlLbl, Text: item.Url, Font: Font{Bold: true}},
					},
				},
				VSeparator{
					MinSize: Size{Height: 1},
					MaxSize: Size{Height: 1},
				},
				Composite{
					AssignTo: &bodyCmp,
					Visible:  false,
					Layout:   VBox{MarginsZero: true},
					Children: []Widget{
						Composite{
							Layout: HBox{MarginsZero: true, Margins: Margins{Bottom: 5}},
							Children: []Widget{
								Label{Text: item.Desc, Alignment: AlignHFarVCenter},
							},
						},
					},
				},
			},
		}).Create(NewBuilder(sip.instructionsContanier)); err != nil {
			log.Fatal(err)
		}

		bodyBuilder := NewBuilder(bodyCmp)
		ele := strings.Split(item.RenderTemplate, ",")
		for _, val := range ele {
			kv := strings.Split(val, "|")
			if len(kv) == 3 {
				if strings.ToLower(strings.TrimSpace(kv[2])) == "input" {
					(Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							Label{Text: kv[0], MinSize: Size{Width: 80}},
							LineEdit{Name: kv[1], MinSize: Size{Width: 120}},
						},
					}).Create(bodyBuilder)

					item.ParamTemplate = strings.ReplaceAll(item.ParamTemplate, fmt.Sprintf("{%s}", kv[1]), fmt.Sprintf("' + %s.Text + '", kv[1]))
				}
			}
		}

		var paramLe *walk.LineEdit
		(LineEdit{
			Text:     Bind("'" + item.ParamTemplate + "'"),
			ReadOnly: true,
			AssignTo: &paramLe,
			MinSize:  Size{Width: 80},
		}).Create(bodyBuilder)

		(Composite{
			Layout: HBox{MarginsZero: true, Margins: Margins{Top: 5}},
			Children: []Widget{
				HSpacer{},
				PushButton{
					Text: "Try it out!",
					OnClicked: func() {
						urlSection := urlLbl.Text()

						idx := strings.Index(urlSection, "/")
						if idx != 0 {
							urlSection = fmt.Sprintf("/%s", urlSection)
						}
						url := fmt.Sprintf("%s%s", host, urlSection)

						sip.sendInstructions(url, paramLe.Text())
					},
				},
			},
		}).Create(bodyBuilder)

		sip.titleHearderCb = append(sip.titleHearderCb, titleCb)
		sip.paramLe = append(sip.paramLe, paramLe)
		sip.cmdCmp = append(sip.cmdCmp, bodyCmp)
	}
}

func (sip *SendInstructionsPage) sendInstructions(url, param string) {

	terno := strings.TrimSpace(sip.ternoLe.Text())
	if terno == "" {
		walk.MsgBox(sip.mainWin, "提醒", "请输入终端编号", walk.MsgBoxIconWarning)
		return
	}

	param = fmt.Sprintf("%s&environment=farm&dids=%s", param, terno)

	bytesData := []byte(param)

	fmt.Println(url)
	fmt.Println(param)
	reader := bytes.NewReader(bytesData)

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		walk.MsgBox(sip.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	client := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, time.Second*60)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 60))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 60,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		walk.MsgBox(sip.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		walk.MsgBox(sip.mainWin, "错误", err.Error(), walk.MsgBoxIconError)
		return
	}

	result := string(respBytes)
	result = strings.ReplaceAll(result, "\"", "")

	walk.MsgBox(sip.mainWin, "结果", result, walk.MsgBoxIconInformation)

	//s := []byte("")
	//fmt.Println(cap(s), len(s))
	//fmt.Println(unsafe.Pointer(&s))
	//
	//s1 := append(s, 'a')
	//fmt.Println(cap(s1), len(s1))
	//fmt.Println(unsafe.Pointer(&s1))
	//
	//s2 := append(s, 'b')
	//fmt.Println(cap(s2), len(s2))
	//fmt.Println(unsafe.Pointer(&s2))
	//fmt.Println(string(s1), "====", string(s2))
}
