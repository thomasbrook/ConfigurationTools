package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"ConfigurationTools/mynotify"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/op/go-logging"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var tlog = logging.MustGetLogger("NewExportHistoryCanPanel")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

// NewExportPanel 新建CAN历史信息导出页面
func NewExportHistoryCanPanel(parent walk.Container, mainWin *TabMainWindow) error {
	rand.Seed(time.Now().UnixNano())
	ecp := &exportCanPage{
		mainWin: mainWin,
	}

	if err := (ScrollView{
		AssignTo: &ecp.ScrollView,
		Layout:   HBox{},
		Children: []Widget{
			VSpacer{},
			Composite{
				Layout:  VBox{},
				MinSize: Size{Width: 260},
				MaxSize: Size{Width: 680},
				Children: []Widget{
					VSpacer{},
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							Label{
								Text: "导出单台设备的历史CAN信息",
								Font: Font{PointSize: 10, Bold: true},
							},
							VSeparator{
								MinSize: Size{Height: 1},
								MaxSize: Size{Height: 1},
							},
							Label{
								Text: "终端编号",
							},
							LineEdit{
								AssignTo: &ecp.ternoLe,
								Text:     "",
							},
							Composite{
								Layout: Grid{Columns: 3, MarginsZero: true},
								Children: []Widget{
									Composite{
										Layout:  VBox{MarginsZero: true, Margins: Margins{Right: 2}},
										MaxSize: Size{Width: 225},
										MinSize: Size{Width: 85},
										Children: []Widget{
											Label{
												Text: "导出数量",
											},
											NumberEdit{
												AssignTo: &ecp.exportNumNe,
												Value:    Bind("10000"),
												Decimals: 0,
												Suffix:   "条",
											},
										},
									},
									Composite{
										Layout:  VBox{MarginsZero: true, Margins: Margins{Right: 2}},
										MaxSize: Size{Width: 225},
										MinSize: Size{Width: 85},
										Children: []Widget{
											Label{Text: "开始时间"},
											DateEdit{
												AssignTo: &ecp.startDate,
												Format:   "yyyy/MM/dd",
												Date:     time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local),
											},
										},
									},
									Composite{
										Layout:  VBox{MarginsZero: true, Margins: Margins{Right: 2}},
										MaxSize: Size{Width: 225},
										MinSize: Size{Width: 85},
										Children: []Widget{
											Label{Text: "结束时间"},
											DateEdit{
												AssignTo: &ecp.endDate,
												Format:   "yyyy/MM/dd",
											},
										},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true, Margins: Margins{Top: 5}},
								Children: []Widget{
									Label{
										AssignTo: &ecp.exportFilePathLbl,
										Text:     "",
									},
									PushButton{
										AssignTo:    &ecp.exportHistoryCanPb,
										Text:        "开始导出",
										ToolTipText: "设备唯一识别号,GNSS时间,经度,纬度,GNSS速度,GNSS状态,海拔高度,水平精度因子,正在使用的卫星数量,发动机转速,机油压力,发动机工作时间,油耗消耗总量,每小时油耗,行驶总里程",
										Alignment:   AlignHFarVCenter,
										MaxSize:     Size{Width: 75},
										OnClicked: func() {
											ecp.exportHistoryCanPb.SetEnabled(false)
											err := ecp.exportCanByTer()
											ecp.exportHistoryCanPb.SetEnabled(true)
											if err != nil {
												mynotify.Error(err.Error())
												return
											}
											mynotify.Message("导出完毕")
										},
									},
								},
							},
						},
					},
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							Label{
								Text: "导出多台设备的历史CAN信息",
								Font: Font{PointSize: 10, Bold: true},
							},
							VSeparator{
								MinSize: Size{Height: 1},
								MaxSize: Size{Height: 1},
							},
							Label{
								Text: "设备清单",
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									LineEdit{
										AssignTo:           &ecp.deviceFilePathLe,
										AlwaysConsumeSpace: true,
										ReadOnly:           true,
									},
									ImageView{
										AssignTo: &ecp.openImageView,
										Image:    "/img/open.png",
										MinSize:  Size{Width: 22},
										OnMouseDown: func(x, y int, button walk.MouseButton) {
											dlg := new(walk.FileDialog)
											dlg.Filter = "文本文件 (*.txt)|*.txt"
											dlg.Title = "请选择设备清单"

											if ok, err := dlg.ShowOpen(ecp.mainWin); err != nil {
												mynotify.Error("初始化文件选择器失败：" + err.Error())
											} else if !ok {
												ecp.deviceFilePathLe.SetText("")
												return
											}
											ecp.deviceFilePathLe.SetText(dlg.FilePath)
										},
									},
								},
							},
							Composite{
								Layout: Grid{Columns: 3, MarginsZero: true},
								Children: []Widget{
									Composite{
										Layout:  VBox{MarginsZero: true, Margins: Margins{Right: 2}},
										MaxSize: Size{Width: 225},
										MinSize: Size{Width: 85},
										Children: []Widget{
											Label{Text: "单台导出数量"},
											NumberEdit{
												AssignTo: &ecp.exportNumNe1,
												Value:    Bind("10000"),
												Decimals: 0,
												Suffix:   "条",
											},
										},
									},
									Composite{
										Layout:  VBox{MarginsZero: true, Margins: Margins{Right: 2}},
										MaxSize: Size{Width: 225},
										MinSize: Size{Width: 85},
										Children: []Widget{
											Label{Text: "开始时间"},
											DateEdit{
												AssignTo: &ecp.startDate1,
												Format:   "yyyy/MM/dd",
												Date:     time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local),
											},
										},
									},
									Composite{
										Layout:  VBox{MarginsZero: true, Margins: Margins{Right: 2}},
										MaxSize: Size{Width: 225},
										MinSize: Size{Width: 85},
										Children: []Widget{
											Label{Text: "结束时间"},
											DateEdit{
												AssignTo: &ecp.endDate1,
												Format:   "yyyy/MM/dd",
											},
										},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true, Margins: Margins{Top: 5}},
								Children: []Widget{
									Label{
										AssignTo:  &ecp.exportFolderPathLbl,
										Text:      "",
										Alignment: AlignHNearVCenter,
									},
									PushButton{
										AssignTo:    &ecp.exportHistoryCanPb1,
										Text:        "开始导出",
										ToolTipText: "设备唯一识别号,GNSS时间,经度,纬度,GNSS速度,GNSS状态,海拔高度,水平精度因子,正在使用的卫星数量,发动机转速,机油压力,发动机工作时间,油耗消耗总量,每小时油耗,行驶总里程",
										Alignment:   AlignHFarVCenter,
										MaxSize:     Size{Width: 75},
										OnClicked: func() {
											ecp.exportHistoryCanPb1.SetEnabled(false)
											err := ecp.exportHistoryCan()
											ecp.exportHistoryCanPb1.SetEnabled(true)
											if err != nil {
												mynotify.Error(err.Error())
												return
											}
										},
									},
								},
							},
						},
					},
					Composite{
						Layout: VBox{MarginsZero: true},
						Children: []Widget{
							Label{
								Text: "导出多台设备的最后一次CAN信息",
								Font: Font{PointSize: 10, Bold: true},
							},
							VSeparator{
								MinSize: Size{Height: 1},
								MaxSize: Size{Height: 1},
							},
							Label{
								Text: "设备清单",
							},
							Composite{
								Layout: HBox{MarginsZero: true},
								Children: []Widget{
									LineEdit{
										AssignTo:           &ecp.deviceFilePathLe1,
										AlwaysConsumeSpace: true,
										ReadOnly:           true,
									},
									ImageView{
										AssignTo: &ecp.openImageView1,
										Image:    "/img/open.png",
										MinSize:  Size{Width: 22},
										OnMouseDown: func(x, y int, button walk.MouseButton) {
											dlg := new(walk.FileDialog)
											dlg.Filter = "文本文件 (*.txt)|*.txt"
											dlg.Title = "请选择设备清单"

											if ok, err := dlg.ShowOpen(ecp.mainWin); err != nil {
												mynotify.Error("初始化文件选择器失败：" + err.Error())
											} else if !ok {
												ecp.deviceFilePathLe1.SetText("")
												return
											}
											ecp.deviceFilePathLe1.SetText(dlg.FilePath)
										},
									},
								},
							},
							Composite{
								Layout: HBox{MarginsZero: true, Margins: Margins{Top: 5}},
								Children: []Widget{
									Label{
										AssignTo:  &ecp.exportFolderPathLbl1,
										Text:      "",
										Alignment: AlignHNearVCenter,
									},
									PushButton{
										AssignTo:    &ecp.exportLastInfoPb,
										Text:        "开始导出",
										ToolTipText: "终端编号,软件版本,硬件版本,电源电压,供电电压,备用电池电压,经度,纬度,地理位置,定位状态,上报时间,接收时间",
										MaxSize:     Size{Width: 75},
										Alignment:   AlignHFarVCenter,
										OnClicked: func() {
											ecp.exportLastInfoPb.SetEnabled(false)
											err := ecp.exportLastCan()
											ecp.exportLastInfoPb.SetEnabled(true)
											if err != nil {
												mynotify.Error(err.Error())
												return
											}
										},
									},
								},
							},
						},
					},
					VSpacer{},
				},
			},
			VSpacer{},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return err
	}

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	//tlog.Debugf("debug %s", Password("secret"))
	//tlog.Info("info")
	//tlog.Notice("notice")
	//tlog.Warning("warning")
	//tlog.Error("err")
	//tlog.Critical("crit")

	return nil
}

type exportCanPage struct {
	*walk.ScrollView

	mainWin *TabMainWindow

	ternoLe            *walk.LineEdit
	exportNumNe        *walk.NumberEdit
	startDate          *walk.DateEdit
	endDate            *walk.DateEdit
	exportFilePathLbl  *walk.Label
	exportHistoryCanPb *walk.PushButton

	deviceFilePathLe    *walk.LineEdit
	openImageView       *walk.ImageView
	exportNumNe1        *walk.NumberEdit
	startDate1          *walk.DateEdit
	endDate1            *walk.DateEdit
	exportFolderPathLbl *walk.Label
	exportHistoryCanPb1 *walk.PushButton

	deviceFilePathLe1    *walk.LineEdit
	openImageView1       *walk.ImageView
	exportFolderPathLbl1 *walk.Label
	exportLastInfoPb     *walk.PushButton
}

// exportCanByTer 导出单台设备CAN历史
func (ecp *exportCanPage) exportCanByTer() error {

	terNo := ecp.ternoLe.Text()
	terNo = strings.Trim(terNo, " ")
	if len(terNo) == 0 {
		return errors.New("请输入终端编号")
	}

	count := ecp.exportNumNe.Value()
	if count <= 0 {
		return errors.New("导出数量应大于0")
	}

	dlg := new(walk.FileDialog)
	dlg.Title = "位置"
	if ok, err := dlg.ShowBrowseFolder(ecp.mainWin); err != nil {
		return errors.New("打开文件选择器失败：" + err.Error())
	} else if !ok {
		ecp.exportFilePathLbl.SetText("")
		return nil
	}
	ecp.exportFilePathLbl.SetText(fmt.Sprintf("位置：%s", dlg.FilePath))

	sdate := ecp.startDate.Date().Format("2006-01-02 15:04:05")
	edate := ecp.endDate.Date().AddDate(0, 0, 1).Format("2006-01-02 15:04:05")
	respBytes, err := dataSource.ListCanHistory(terNo, sdate, edate, int(count))
	if err != nil {
		return err
	}

	// Json解析
	pageInfo := model.PageInfoModel{}
	err = json.Unmarshal(respBytes, &pageInfo)
	if err != nil {
		tlog.Fatal(err)
	}

	// 声明文件名、表头，并创建文件
	filePath := fmt.Sprintf("%s\\%s_%s.csv", dlg.FilePath, terNo, time.Now().Format("20060102150405"))

	title := []string{"设备唯一识别号", "GNSS时间", "经度", "纬度", "GNSS速度", "GNSS状态", "海拔高度",
		"水平精度因子", "正在使用的卫星数量", "发动机转速", "机油压力", "发动机工作时间", "油耗消耗总量",
		"每小时油耗", "行驶总里程"}

	_, err = createFile(filePath, title)
	if err != nil {
		return err
	}

	// 形成二维数组，并输出到文件
	buffer := [][]string{}
	for j := 0; j < len(pageInfo.Rows); j++ {
		if pageInfo.Rows[j] == nil {
			continue
		}

		gpsDate, _ := time.Parse("20060102150405", pageInfo.Rows[j].GnssTime)

		ewt, err := strconv.ParseFloat(pageInfo.Rows[j].EngineWorkTime, 64)
		if err != nil {
			ewt = 0
		}
		ewt1, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", ewt), 64)

		op, err := strconv.ParseFloat(pageInfo.Rows[j].OilPressure, 64)
		if err != nil {
			op = 0
		}
		op1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", op/1000), 64)

		gpsStatus := pageInfo.Rows[j].Status
		switch gpsStatus {
		case "0":
			gpsStatus = "2"
			break
		case "1":
			gpsStatus = "0"
			break
		case "2":
			gpsStatus = "1"
			break
		case "3":
			gpsStatus = "2"
			break
		case "4":
			gpsStatus = "6"
			break
		}

		al, err := strconv.ParseFloat(pageInfo.Rows[j].Altitude, 64)
		if err != nil {
			al = 0
		}
		al1, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", al), 64)

		gwsz, err := strconv.ParseFloat(pageInfo.Rows[j].GWSZ, 64)
		if err != nil {
			gwsz = 0
		}

		gwsz1, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", gwsz), 64)

		data := []string{terNo,
			gpsDate.Format("2006-01-02 15:04:05"),
			pageInfo.Rows[j].Longitude,
			pageInfo.Rows[j].Latitude,
			pageInfo.Rows[j].Speed,
			gpsStatus,
			strconv.FormatFloat(al1, 'f', 0, 64),
			pageInfo.Rows[j].Hdop,
			pageInfo.Rows[j].UsedSatelliteNumber,
			pageInfo.Rows[j].EngineSpeed,
			strconv.FormatFloat(op1, 'f', 2, 64),
			strconv.FormatFloat(ewt1, 'f', 1, 64),
			pageInfo.Rows[j].TotalFuelConsumption,
			pageInfo.Rows[j].FuelConsumptionPerHour,
			strconv.FormatFloat(gwsz1, 'f', 0, 64)}

		buffer = append(buffer, data)
	}

	// 批量输出信息到csv文件
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	w.WriteAll(buffer)
	w.Flush()

	file.Close()

	return nil
}

// exportHistoryCan 批量导出多台设备CAN历史
func (ecp *exportCanPage) exportHistoryCan() error {

	deviceFilePath := ecp.deviceFilePathLe.Text()
	deviceFilePath = strings.Trim(deviceFilePath, " ")
	if len(deviceFilePath) == 0 {
		return errors.New("请选择设备清单")
	}

	count := ecp.exportNumNe1.Value()
	if count <= 0 {
		return errors.New("导出数量应大于0")
	}

	dlg := new(walk.FileDialog)
	dlg.Title = "位置"
	if ok, err := dlg.ShowBrowseFolder(ecp.mainWin); err != nil {
		return errors.New("打开文件选择器失败：" + err.Error())
	} else if !ok {
		ecp.exportFolderPathLbl.SetText("")
		return nil
	}
	ecp.exportFolderPathLbl.SetText(fmt.Sprintf("位置：%s", dlg.FilePath))

	pathSection := strings.Split(deviceFilePath, "\\")
	name := pathSection[len(pathSection)-1]

	subStr := strings.Split(name, ".")
	topicName := ""
	if len(subStr) > 0 {
		topicName = subStr[0]
	}

	// 检查是否存在统计文件，不存在则创建
	statHeader := []string{"设备唯一识别号", "CAN数量", "文件索引"}
	statFileName := fmt.Sprintf("%s\\%s_stat.csv", dlg.FilePath, topicName)
	_, err := createFile(statFileName, statHeader)
	if err != nil {
		return err
	}

	// 读取统计信息
	statFile, err := os.Open(statFileName)
	if err != nil {
		return err
	}
	defer statFile.Close()

	r := csv.NewReader(statFile)
	terTxt, err := r.ReadAll()
	if err != nil {
		return err
	}

	// 统计已处理过的终端数据
	processedTers := make(map[string]string)
	lines := 0
	fileNum := -1
	for idx, row := range terTxt {
		if idx == 0 {
			continue
		}

		if len(row) == 3 {
			_, isExist := processedTers[row[0]]
			if !isExist {
				processedTers[row[0]] = row[0]
			}

			size, err := strconv.Atoi(row[1])
			if err == nil {
				lines += size
			}
			index := strings.Split(row[2], ",")
			if len(index) > 0 {
				idx, err := strconv.Atoi(index[len(index)-1])
				if err != nil {
					fileNum = idx
				}
			}
		}
	}

	// 读取设备清单

	deviceFile, err := os.Open(deviceFilePath)
	if err != nil {
		return err
	}
	defer deviceFile.Close()

	terFile := csv.NewReader(deviceFile)
	terTxt, err = terFile.ReadAll()
	if err != nil {
		return err
	}

	// 忽略已完成导出的终端设备
	ters := []*string{}
	for _, row := range terTxt {
		if len(row) == 1 {
			_, isExist := processedTers[row[0]]
			if !isExist {
				ters = append(ters, &row[0])
			}
		}
	}

	sdate := ecp.startDate1.Date().Format("2006-01-02 15:04:05")
	edate := ecp.endDate1.Date().AddDate(0, 0, 1).Format("2006-01-02 15:04:05")

	if len(ters) == 0 {
		walk.MsgBox(ecp.mainWin, "", "导出完毕，如需重新导出，请删除导出的文件", walk.MsgBoxIconInformation)
		return nil
	}

	go func() {
		exportCsv(ters, statFileName, fmt.Sprintf("%s\\%s", dlg.FilePath, topicName), sdate, edate, int(count), lines, fileNum)
		mynotify.Message("导出完毕")
	}()

	return nil
}

// exportLastCan 批量导出多台设备最后一次CAN信息
func (ecp *exportCanPage) exportLastCan() error {

	terFilePath := ecp.deviceFilePathLe1.Text()
	terFilePath = strings.Trim(terFilePath, " ")

	if len(terFilePath) == 0 {
		return errors.New("请选择设备清单文件")
	}

	dlg := new(walk.FileDialog)
	dlg.Title = "位置"

	if ok, err := dlg.ShowBrowseFolder(ecp.mainWin); err != nil {
		return errors.New("打开文件选择器失败：" + err.Error())
	} else if !ok {
		ecp.exportFolderPathLbl1.SetText("")
		return nil
	}
	ecp.exportFolderPathLbl1.SetText(fmt.Sprintf("位置：%s", dlg.FilePath))

	pathSection := strings.Split(terFilePath, "\\")
	name := pathSection[len(pathSection)-1]

	subStr := strings.Split(name, ".")
	topicName := ""
	if len(subStr) > 0 {
		topicName = subStr[0]
	}

	// 检查是否存在统计文件，不存在则创建
	statHeader := []string{"设备唯一识别号", "数量"}
	statFileName := fmt.Sprintf("%s\\%s_stat.csv", dlg.FilePath, topicName)
	_, err := createFile(statFileName, statHeader)
	if err != nil {
		return err
	}

	// 读取统计信息
	statFile, err := os.Open(statFileName)
	statFileReader := csv.NewReader(statFile)
	terTxt, err := statFileReader.ReadAll()
	if err != nil {
		return err
	}

	// 将已处理的设备，哈希存储
	processedTers := make(map[string]string)
	for idx, row := range terTxt {
		if idx == 0 {
			continue
		}

		if len(row) == 2 {
			_, isExist := processedTers[row[0]]
			if !isExist {
				processedTers[row[0]] = row[0]
			}
		}
	}

	// 读取设备清单
	terFile, _ := os.Open(terFilePath)
	terFileReader := csv.NewReader(terFile)
	terTxt, err = terFileReader.ReadAll()
	if err != nil {
		return err
	}

	device := []*string{}
	for _, row := range terTxt {
		if len(row) == 1 {
			_, isExist := processedTers[row[0]]
			if !isExist {
				device = append(device, &row[0])
			}
		}
	}

	if len(device) == 0 {
		walk.MsgBox(ecp.mainWin, "", "导出完毕，如需重新导出，请删除导出的文件", walk.MsgBoxIconInformation)
		return nil
	}

	go func() {
		exportLastCan(device, statFileName, fmt.Sprintf("%s\\%s", dlg.FilePath, topicName))
		mynotify.Message("导出完毕")
	}()
	return nil
}

func exportLastCan(ternos []*string, statFilePath string, dataFilePath string) error {
	tlog.Notice(fmt.Sprintf("共 %d 台设备\r\n", len(ternos)))

	// 查看文件是否存在，不存在则创建
	title := []string{"终端编号", "软件版本", "硬件版本", "电源电压", "供电电压", "备用电池电压", "经度", "纬度", "地理位置", "定位状态", "上报时间", "接收时间"}
	dataFileName := fmt.Sprintf("%s.csv", dataFilePath)
	_, err := createFile(dataFileName, title)
	if err != nil {
		return err
	}

	// 打开数据文件
	dataFile, err := os.OpenFile(dataFileName, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer dataFile.Close()
	dataWriter := csv.NewWriter(dataFile)

	// 打开统计文件
	statFile, err := os.OpenFile(statFilePath, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer statFile.Close()
	statWriter := csv.NewWriter(statFile)

	for i := 0; i < len(ternos); i++ {

		var respBytes []byte
		var err error

		for retryIdx := 0; retryIdx < 5; retryIdx++ {

			tlog.Notice(fmt.Sprintf("%s at %d/%d", *ternos[i], i+1, len(ternos)))

			t := time.Now()
			respBytes, err = dataSource.GetLastestCan(*ternos[i])
			elapsed := time.Since(t)
			tlog.Info(" Api request:", elapsed)

			if err != nil {
				tlog.Error(fmt.Sprintf("%s", err))
				tlog.Notice(fmt.Sprintf("6秒后重试"))
				time.Sleep(6 * time.Second)
				continue
			}

			break
		}

		if respBytes == nil {
			break
		}

		// Json解析
		t := time.Now()
		pageInfo := model.LastRowModel{}
		err1 := json.Unmarshal(respBytes, &pageInfo)
		if err1 != nil {
			tlog.Error(err1.Error())
		}
		elapsed := time.Since(t)
		tlog.Info(" json parse:", elapsed)

		size := len(pageInfo.Rows)
		tlog.Info(fmt.Sprintf(" total %d", size))

		// 形成二维数组
		data := []string{}
		if size == 0 {
			data = []string{*ternos[i],
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				""}
		} else {
			gpsDate, _ := time.Parse("20060102150405", pageInfo.Rows[0].GPSDateTime)
			reviceTime, _ := time.Parse("20060102150405", pageInfo.Rows[0].TIME)

			//0 GPS已定位；1 GPS未定位；2非差分定位; 3 差分定位; 4正在估算
			gpsStatus := pageInfo.Rows[0].Status
			switch gpsStatus {
			case "0":
				gpsStatus = "GPS已定位"
				break
			case "1":
				gpsStatus = "GPS未定位"
				break
			case "2":
				gpsStatus = "非差分定位"
				break
			case "3":
				gpsStatus = "差分定位"
				break
			case "4":
				gpsStatus = "正在估算"
				break
			}

			addr := ""
			lat, err1 := strconv.ParseFloat(pageInfo.Rows[0].Latitude, 64)
			lng, err2 := strconv.ParseFloat(pageInfo.Rows[0].Longitude, 64)
			if err1 == nil && err2 == nil {
				_addr, err3 := getAddressByCoordinate(lat, lng)
				if err3 != nil {
					addr = err3.Error()
					tlog.Error(fmt.Sprintf("获取位置失败：%s", err3.Error()))
					for i := 0; i < 5; i++ {
						tlog.Notice(fmt.Sprintf("3秒后重试"))
						time.Sleep(3 * time.Millisecond)
						_addr, err3 = getAddressByCoordinate(lat, lng)
						if err3 == nil {
							addr = _addr
							tlog.Info(addr)
							break
						}
					}
				} else {
					addr = _addr
					tlog.Info(addr)
				}
			}

			data = []string{*ternos[i],
				pageInfo.Rows[0].SoftwareVersion,
				pageInfo.Rows[0].HardwareVersion,
				pageInfo.Rows[0].PowerSupplyVolt,
				pageInfo.Rows[0].InputPowerVoltage,
				pageInfo.Rows[0].InputBatteryVoltage,
				pageInfo.Rows[0].Longitude,
				pageInfo.Rows[0].Latitude,
				addr,
				gpsStatus,
				gpsDate.Format("2006/01/02 15:04:05"),
				reviceTime.Format("2006/01/02 15:04:05")}
		}

		dataWriter.Write(data)
		dataWriter.Flush()

		statWriter.Write([]string{*ternos[i], strconv.Itoa(size)})
		statWriter.Flush()

		time.Sleep(100 * time.Millisecond)
		tlog.Debug("\r\n")
	}

	dataFile.Close()
	statFile.Close()

	return nil
}

// getAddressByCoordinate 经纬度转地址
func getAddressByCoordinate(lat float64, lng float64) (string, error) {
	respBytes, err := dataSource.GetAddressByCoordinate(lat, lng)

	if err != nil {
		fmt.Println(fmt.Sprintf("%s", err))
	}

	//addr := model.BaiduAddress{}
	addr := model.BcldAddress{}
	err1 := json.Unmarshal(respBytes, &addr)
	if err1 != nil {
		fmt.Println(fmt.Sprintf("%s", err1))
		return "", err1
	}

	if addr.Status == 1 {
		//return fmt.Sprintf("%s,%s,%s,%s,%s",
		//	addr.Result.AddressComponent.Province,
		//	addr.Result.AddressComponent.City,
		//	addr.Result.AddressComponent.District,
		//	addr.Result.AddressComponent.Town,
		//	addr.Result.Sematic_Description), nil

		if len(addr.Data) > 0 {
			return addr.Data[0].Value, nil
		}
		return "", nil
	} else {
		return "", errors.New("解析失败")
	}
}

func exportCsv(ters []*string, statFilePath string, dataFilePath string, startDate string, endDate string, count int, lines int, fileNum int) error {
	tlog.Notice(fmt.Sprintf("共 %d 台设备\r\n", len(ters)))

	title := []string{"设备唯一识别号", "GNSS时间", "经度", "纬度", "GNSS速度", "GNSS状态", "海拔高度",
		"水平精度因子", "正在使用的卫星数量", "发动机转速", "机油压力", "发动机工作时间", "油耗消耗总量",
		"每小时油耗", "行驶总里程"}

	for i := 0; i < len(ters); i++ {
		var respBytes []byte
		var err error
		for retryIdx := 0; retryIdx < 5; retryIdx++ {
			tlog.Notice(fmt.Sprintf("%s at %d/%d", *ters[i], i+1, len(ters)))
			t := time.Now()
			respBytes, err = dataSource.ListCanHistory(*ters[i], startDate, endDate, count)
			elapsed := time.Since(t)
			tlog.Info(" Api request:", elapsed)

			if err != nil {
				tlog.Error(fmt.Sprintf("%s", err))
				tlog.Info(fmt.Sprintf("6秒后重试"))
				time.Sleep(6 * time.Second)
				continue
			}

			break
		}

		if respBytes == nil {
			break
		}

		// Json解析
		t := time.Now()
		pageInfo := model.PageInfoModel{}
		err1 := json.Unmarshal(respBytes, &pageInfo)
		if err1 != nil {
			tlog.Error(err1)
		}
		elapsed := time.Since(t)
		tlog.Info(" json parse:", elapsed)

		size := len(pageInfo.Rows)
		tlog.Info(fmt.Sprintf(" total %d", size))

		// 1、如果CAN数量为0，输出统计，并继续循环处理
		if size == 0 {
			statTxt, err := os.OpenFile(statFilePath, os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				tlog.Fatal(err)
			}
			wstat := csv.NewWriter(statTxt)
			wstat.Write([]string{*ters[i], strconv.Itoa(0), strconv.Itoa(fileNum)})
			wstat.Flush()
			continue
		}

		// 2、处理CAN数据，输出二维数组
		buffer := [][]string{}
		for j := 0; j < len(pageInfo.Rows); j++ {

			gpsDate, _ := time.Parse("20060102150405", pageInfo.Rows[j].GnssTime)

			ewt, err := strconv.ParseFloat(pageInfo.Rows[j].EngineWorkTime, 64)
			if err != nil {
				ewt = 0
			}

			ewt1, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", ewt), 64)

			op, err := strconv.ParseFloat(pageInfo.Rows[j].OilPressure, 64)
			if err != nil {
				op = 0
			}
			op1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", op/1000), 64)

			gpsStatus := pageInfo.Rows[j].Status
			switch gpsStatus {
			case "0":
				gpsStatus = "2"
				break
			case "1":
				gpsStatus = "0"
				break
			case "2":
				gpsStatus = "1"
				break
			case "3":
				gpsStatus = "2"
				break
			case "4":
				gpsStatus = "6"
				break
			}

			al, err := strconv.ParseFloat(pageInfo.Rows[j].Altitude, 64)
			if err != nil {
				al = 0
			}
			al1, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", al), 64)

			gwsz, err := strconv.ParseFloat(pageInfo.Rows[j].GWSZ, 64)
			if err != nil {
				gwsz = 0
			}

			gwsz1, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", gwsz), 64)

			data := []string{*ters[i],
				gpsDate.Format("2006-01-02 15:04:05"),
				pageInfo.Rows[j].Longitude,
				pageInfo.Rows[j].Latitude,
				pageInfo.Rows[j].Speed,
				gpsStatus,
				strconv.FormatFloat(al1, 'f', 0, 64),
				pageInfo.Rows[j].Hdop,
				pageInfo.Rows[j].UsedSatelliteNumber,
				pageInfo.Rows[j].EngineSpeed,
				strconv.FormatFloat(op1, 'f', 2, 64),
				strconv.FormatFloat(ewt1, 'f', 1, 64),
				pageInfo.Rows[j].TotalFuelConsumption,
				pageInfo.Rows[j].FuelConsumptionPerHour,
				strconv.FormatFloat(gwsz1, 'f', 0, 64)}

			buffer = append(buffer, data)
		}

		// 拆分数据，输出到不同的csv文件
		buf := [][][]string{}
		if (lines/100000+1)*100000-lines >= size {
			buf = append(buf, buffer)
		} else {
			diff := (lines/100000+1)*100000 - lines
			buf = append(buf, buffer[:diff])
			buf = append(buf, buffer[diff:])
		}

		// 批量输出信息到csv文件
		total := 0
		index := []string{}
		for k := 0; k < len(buf); k++ {
			csvName := fmt.Sprintf("%s_%d.csv", dataFilePath, lines/100000)

			// 如果当前索引号 与 实际索引号不一致，则需要创建一个新文件
			if fileNum != (lines / 100000) {
				affectRows, err := createFile(csvName, title)

				if err != nil {
					tlog.Fatal(err)
				}

				fileNum = lines / 100000
				lines += affectRows
			}

			txt, err := os.OpenFile(csvName, os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				tlog.Fatal(err)
			}
			w := csv.NewWriter(txt)
			w.WriteAll(buf[k])
			w.Flush()
			txt.Close()

			tlog.Info(fmt.Sprintf(" total size at index of %d, %d", k, len(buf[k])))
			lines = lines + len(buf[k])

			total += len(buf[k])
			index = append(index, strconv.Itoa(fileNum))
		}

		// 输出索引统计
		statFile, err := os.OpenFile(statFilePath, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			tlog.Error(err)
		}
		wstat := csv.NewWriter(statFile)
		wstat.Write([]string{*ters[i], strconv.Itoa(total), strings.Join(index, ",")})
		wstat.Flush()
		statFile.Close()

		fmt.Println()
	}

	return nil
}

// createFile 创建带表头的csv文件
func createFile(csvName string, title []string) (int, error) {
	file, err := os.Open(csvName)
	defer file.Close()

	if err == nil {
		return 0, nil
	}

	if os.IsNotExist(err) {

		file, err := os.Create(csvName)
		if err != nil {
			return 0, err
		}
		defer file.Close()

		file.WriteString("\xEF\xBB\xBF")

		w := csv.NewWriter(file)
		w.Write(title)
		w.Flush()

		return 1, nil
	}
	return 0, err
}
