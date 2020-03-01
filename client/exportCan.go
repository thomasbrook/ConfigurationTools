package client

import (
	"ConfigurationTools/dataSource"
	"ConfigurationTools/model"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// NewExportPanel 新建CAN历史信息导出页面
func NewExportPanel(parent walk.Container, mainWin *TabMainWindow) error {
	rand.Seed(time.Now().UnixNano())
	ecp := &exportCanPage{
		mainWin: mainWin,
	}

	if err := (Composite{
		AssignTo: &ecp.Composite,
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			ScrollView{
				Layout:VBox{},
				Children:[]Widget{
					Composite{
						Layout: Grid{Columns: 1},
						Children: []Widget{
							Label{
								Text: "车辆类型ID",
							},
							LineEdit{
								AssignTo: &ecp.vehicleIDEdit,
							},
							Label{
								Text: "文件名称",
							},
							LineEdit{
								AssignTo: &ecp.fileNameEdit,
							},
							PushButton{
								AssignTo: &ecp.exportBtn,
								Text:     "按车辆类型导出",
								OnClicked: func() {
									err := ecp.exportCan()
									if err != nil {
										log.Print(err)
									}
								},
							},
						},
					},
					Composite{
						Layout: Grid{Columns: 1},
						Children: []Widget{
							Label{
								Text: "输入机构ID",
							},
							LineEdit{
								AssignTo: &ecp.orgIdEdit,
							},
							PushButton{
								AssignTo: &ecp.exportBtnByOrg,
								Text:     "按机构ID导出",
								OnClicked: func() {
									err := ecp.exportCanByOrg()
									if err != nil {
										log.Print(err)
									}
								},
							},
						},
					},
					Composite{
						Layout: Grid{Columns: 1},
						Children: []Widget{
							Label{
								Text: "终端编号",
							},
							LineEdit{
								AssignTo: &ecp.ternoEdit,
							},
							Label{
								Text: "文件路径",
							},
							LineEdit{
								AssignTo: &ecp.filePathEdit,
							},
							PushButton{
								AssignTo: &ecp.exportBtnByTer,
								Text:     "单终端导出",
								OnClicked: func() {
									err := ecp.exportCanByTer()
									if err != nil {
										log.Print(err)
									}
								},
							},
						},
					},
					Composite{
						Layout: Grid{Columns: 1},
						Children: []Widget{
							Label{
								Text: "源终端清单路径",
							},
							LineEdit{
								AssignTo: &ecp.sourcFilePath,
							},
							Label{
								Text: "文件名称",
							},
							LineEdit{
								AssignTo: &ecp.sfileNameEdit,
							},
							PushButton{
								AssignTo: &ecp.exportBtnByMultiTer,
								Text:     "多终端导出",
								OnClicked: func() {
									err := ecp.exportCanByMultiTer()
									if err != nil {
										log.Println(err)
									}
								},
							},
						},
					},
					Composite{
						Layout: Grid{Columns: 1},
						Children: []Widget{
							Label{
								Text: "源终端清单路径",
							},
							LineEdit{
								AssignTo: &ecp.sourcFilePath1,
							},
							Label{
								Text: "文件名称",
							},
							LineEdit{
								AssignTo: &ecp.sfileNameEdit1,
							},
							PushButton{
								AssignTo: &ecp.exportLastInfo,
								Text:     "导出最后一次定位信息",
								OnClicked: func() {
									err := ecp.exportLastCan()
									if err != nil {
										log.Println(err)
									}
								},
							},
						},
					},
				},
			},
		},
	}).Create(NewBuilder(parent)); err != nil {
		return err
	}

	return nil
}

type exportCanPage struct {
	*walk.Composite

	mainWin *TabMainWindow

	vehicleIDEdit *walk.LineEdit
	fileNameEdit  *walk.LineEdit
	exportBtn     *walk.PushButton

	orgIdEdit      *walk.LineEdit
	exportBtnByOrg *walk.PushButton

	ternoEdit      *walk.LineEdit
	filePathEdit   *walk.LineEdit
	exportBtnByTer *walk.PushButton

	sourcFilePath       *walk.LineEdit
	sfileNameEdit       *walk.LineEdit
	exportBtnByMultiTer *walk.PushButton

	sourcFilePath1 *walk.LineEdit
	sfileNameEdit1 *walk.LineEdit
	exportLastInfo *walk.PushButton
}

type statInfo struct {
	did     string
	canSize int
	index   int
}

func (ecp *exportCanPage) exportLastCan() error {
	ecp.exportLastInfo.SetEnabled(false)

	sfilePath := ecp.sourcFilePath1.Text()
	sfilePath = strings.Trim(sfilePath, " ")

	if len(sfilePath) == 0 {
		fmt.Println("请输入设备清单路径")
		ecp.exportLastInfo.SetEnabled(true)
		return nil
	}

	sfileName := ecp.sfileNameEdit1.Text()
	sfileName = strings.Trim(sfileName, " ")
	if len(sfileName) == 0 {
		fmt.Println("请输入文件名")
		ecp.exportLastInfo.SetEnabled(true)
		return nil
	}

	// 检查是否存在统计文件，不存在则创建
	stat := []string{"设备唯一识别号", "CAN数量", "文件索引"}
	statfname := fmt.Sprintf("export/%s_stat.csv", sfileName)
	_, err := createFile(statfname, stat)
	if err != nil {
		ecp.exportLastInfo.SetEnabled(true)
		log.Fatal(err)
	}

	// 读取统计信息
	fs, _ := os.Open(statfname)
	r := csv.NewReader(fs)
	content, err := r.ReadAll()
	if err != nil {
		ecp.exportLastInfo.SetEnabled(true)
		log.Fatal(err)
	}

	processedTers := make(map[string]string)
	lines := 0
	fileNum := -1
	for idx, row := range content {
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

			index, err := strconv.Atoi(row[2])
			if err != nil {
				fileNum = index
			}
		}
	}

	// 读取设备清单
	devicefs, _ := os.Open(sfilePath)
	devicer := csv.NewReader(devicefs)
	content, err = devicer.ReadAll()
	if err != nil {
		ecp.exportLastInfo.SetEnabled(true)
		log.Fatal(err)
	}

	device := []*string{}
	for _, row := range content {
		if len(row) == 1 {
			_, isExist := processedTers[row[0]]
			if !isExist {
				device = append(device, &row[0])
			}
		}
	}

	fmt.Println(len(device))
	fmt.Println()

	exportLastCan(device, statfname, fmt.Sprintf("export/%s", sfileName), lines, fileNum)
	ecp.exportLastInfo.SetEnabled(true)
	return nil
}

func (ecp *exportCanPage) exportCanByMultiTer() error {
	ecp.exportBtnByMultiTer.SetEnabled(false)

	sfilePath := ecp.sourcFilePath.Text()
	sfilePath = strings.Trim(sfilePath, " ")

	if len(sfilePath) == 0 {
		fmt.Println("请输入设备清单路径")
		ecp.exportBtnByMultiTer.SetEnabled(true)
		return nil
	}

	sfileName := ecp.sfileNameEdit.Text()
	sfileName = strings.Trim(sfileName, " ")
	if len(sfileName) == 0 {
		fmt.Println("请输入文件名")
		ecp.exportBtnByMultiTer.SetEnabled(true)
		return nil
	}

	// 检查是否存在统计文件，不存在则创建
	stat := []string{"设备唯一识别号", "CAN数量", "文件索引"}
	statfname := fmt.Sprintf("export/%s_stat.csv", sfileName)
	_, err := createFile(statfname, stat)
	if err != nil {
		ecp.exportBtnByMultiTer.SetEnabled(true)
		log.Fatal(err)
	}

	// 读取统计信息
	fs, _ := os.Open(statfname)
	r := csv.NewReader(fs)
	content, err := r.ReadAll()
	if err != nil {
		ecp.exportBtnByMultiTer.SetEnabled(true)
		log.Fatal(err)
	}

	processedTers := make(map[string]string)
	lines := 0
	fileNum := -1
	for idx, row := range content {
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

			index, err := strconv.Atoi(row[2])
			if err != nil {
				fileNum = index
			}
		}
	}

	// 读取设备清单
	devicefs, _ := os.Open(sfilePath)
	devicer := csv.NewReader(devicefs)
	content, err = devicer.ReadAll()
	if err != nil {
		ecp.exportBtnByMultiTer.SetEnabled(true)
		log.Fatal(err)
	}

	device := []*string{}
	for _, row := range content {
		if len(row) == 1 {
			_, isExist := processedTers[row[0]]
			if !isExist {
				device = append(device, &row[0])
			}
		}
	}

	fmt.Println(len(device))
	fmt.Println()

	exportCsv(device, statfname, fmt.Sprintf("export/%s", sfileName), lines, fileNum)
	ecp.exportBtnByMultiTer.SetEnabled(true)
	return nil
}

func (ecp *exportCanPage) exportCanByTer() error {
	ecp.exportBtnByTer.SetEnabled(false)

	terNo := ecp.ternoEdit.Text()
	terNo = strings.Trim(terNo, " ")

	if len(terNo) == 0 {
		fmt.Println("请输入终端编号")
		ecp.exportBtnByTer.SetEnabled(true)
		return nil
	}

	filePath := ecp.filePathEdit.Text()
	filePath = strings.Trim(filePath, " ")
	if len(filePath) == 0 {
		fmt.Println("请输入文件路径")
		ecp.exportBtnByTer.SetEnabled(true)
		return nil
	}

	maxLines := 10000
	currLines := 0
	endTime := time.Date(2019, 12, 27, 23, 59, 59, 0, time.Local)
	for currLines < maxLines {

		time.Sleep(6 * time.Second)

		var respBytes []byte
		var err error
		for retryIdx := 0; retryIdx < 5; retryIdx++ {
			respBytes, err = dataSource.ListCanHistory(terNo, "2019-01-01 00:00:00", endTime.Format("2006-01-02 15:04:05"), 10000)

			if err != nil {
				fmt.Println(fmt.Sprintf("%s", err))
				fmt.Println(fmt.Sprintf("12秒后重试"))
				time.Sleep(12 * time.Second)
				continue
			}

			break
		}

		if respBytes == nil {
			ecp.exportBtnByTer.SetEnabled(true)
			break
		}

		// Json解析
		pageInfo := model.PageInfoModel{}
		err = json.Unmarshal(respBytes, &pageInfo)
		if err != nil {
			log.Fatal(err)
		}

		size := len(pageInfo.Rows)
		if size == 0 {
			ecp.exportBtnByTer.SetEnabled(true)
			break
		}

		// 形成二维数组
		buffer := [][]string{}
		for j := 0; j < len(pageInfo.Rows); j++ {
			data := []string{terNo,
				pageInfo.Rows[j].Can.GnssTime,
				pageInfo.Rows[j].Can.Longitude,
				pageInfo.Rows[j].Can.Latitude,
				pageInfo.Rows[j].Can.Speed,
				pageInfo.Rows[j].Can.Status,
				pageInfo.Rows[j].Can.Altitude,
				pageInfo.Rows[j].Can.Hdop,
				pageInfo.Rows[j].Can.UsedSatelliteNumber,
				pageInfo.Rows[j].Can.EngineSpeed,
				pageInfo.Rows[j].Can.OilPressure,
				pageInfo.Rows[j].Can.EngineWorkTime,
				pageInfo.Rows[j].Can.TotalFuelConsumption,
				pageInfo.Rows[j].Can.FuelConsumptionPerHour,
				pageInfo.Rows[j].Can.GWSZ}

			buffer = append(buffer, data)
		}

		endTime, err = time.Parse("20060102150405", buffer[len(buffer)-1][1])
		if err != nil {
			ecp.exportBtnByTer.SetEnabled(true)
			log.Fatal(err)
		}
		endTime = endTime.Add(time.Duration(-1) * time.Second)

		// 批量输出信息到csv文件
		txt, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			ecp.exportBtnByTer.SetEnabled(true)
			log.Fatal(err)
		}
		w := csv.NewWriter(txt)
		w.WriteAll(buffer)
		w.Flush()

		txt.Close()

		currLines += len(buffer)

		fmt.Println(fmt.Sprintf("%d %d", currLines, len(buffer)))
		fmt.Println()
	}

	ecp.exportBtnByTer.SetEnabled(true)
	return nil
}

func (ecp *exportCanPage) exportCanByOrg() error {

	ecp.exportBtnByOrg.SetEnabled(false)

	orgID := ecp.orgIdEdit.Text()
	orgID = strings.Trim(orgID, " ")
	if len(orgID) == 0 {
		ecp.exportBtnByOrg.SetEnabled(true)
		return nil
	}

	// 检查是否存在统计文件，不存在则创建
	stat := []string{"设备唯一识别号", "CAN数量", "文件索引"}
	statfname := fmt.Sprintf("export/%s_stat.csv", orgID)
	_, err := createFile(statfname, stat)
	if err != nil {
		ecp.exportBtnByOrg.SetEnabled(true)
		log.Fatal(err)
	}

	// 读取统计信息
	fs, _ := os.Open(statfname)
	r := csv.NewReader(fs)
	content, err := r.ReadAll()
	if err != nil {
		ecp.exportBtnByOrg.SetEnabled(true)
		log.Fatal(err)
	}

	processedTers := []string{}
	lines := 0
	fileNum := -1
	for idx, row := range content {
		if idx == 0 {
			continue
		}

		if len(row) == 3 {
			processedTers = append(processedTers, row[0])
			size, err := strconv.Atoi(row[1])
			if err == nil {
				lines += size
			}

			index, err := strconv.Atoi(row[2])
			if err != nil {
				fileNum = index
			}
		}
	}

	ternos, err := dataSource.ListVehicleByOrg(orgID, processedTers)
	if err != nil {
		ecp.exportBtnByOrg.SetEnabled(true)
		return err
	}
	fmt.Println(len(ternos))
	fmt.Println()

	exportCsv(ternos, statfname, "export/data", lines, fileNum)
	ecp.exportBtnByOrg.SetEnabled(true)
	return nil
}

func (ecp *exportCanPage) exportCan() error {
	ecp.exportBtn.SetEnabled(false)

	vtID := ecp.vehicleIDEdit.Text()
	vtID = strings.Trim(vtID, " ")

	if len(vtID) == 0 {
		fmt.Println("请输入车辆类型ID")
		ecp.exportBtn.SetEnabled(true)
		return nil
	}

	fileName := ecp.fileNameEdit.Text()
	fileName = strings.Trim(fileName, " ")
	if len(fileName) == 0 {
		fmt.Println("请输入文件名")
		ecp.exportBtn.SetEnabled(true)
		return nil
	}

	// 检查是否存在统计文件，不存在则创建
	stat := []string{"设备唯一识别号", "CAN数量", "文件索引"}
	statfname := fmt.Sprintf("export/%s_stat.csv", fileName)
	_, err := createFile(statfname, stat)
	if err != nil {
		ecp.exportBtn.SetEnabled(true)
		log.Fatal(err)
	}

	// 读取统计信息
	fs, _ := os.Open(statfname)
	r := csv.NewReader(fs)
	content, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	processedTers := []string{}
	lines := 0
	fileNum := -1
	for idx, row := range content {
		if idx == 0 {
			continue
		}

		if len(row) == 3 {
			processedTers = append(processedTers, row[0])
			size, err := strconv.Atoi(row[1])
			if err == nil {
				lines += size
			}

			index, err := strconv.Atoi(row[2])
			if err != nil {
				fileNum = index
			}
		}
	}

	ternos, err := dataSource.ListVehicle(vtID, processedTers)
	if err != nil {
		return err
	}
	fmt.Println(len(ternos))
	fmt.Println()

	exportCsv(ternos, statfname, fmt.Sprintf("export/%s", fileName), lines, fileNum)
	ecp.exportBtn.SetEnabled(true)
	return nil
}

func exportLastCan(ternos []*string, statfname string, datafname string, lines int, fileNum int) error {

	title := []string{"终端编号", "软件版本", "硬件版本", "电源电压", "供电电压", "备用电池电压", "经度", "纬度", "定位状态", "上报时间", "接收时间"}

	for i := 0; i < len(ternos); i++ {

		// time.Sleep(6 * time.Second)

		var respBytes []byte
		var err error
		for retryIdx := 0; retryIdx < 5; retryIdx++ {
			fmt.Println(fmt.Sprintf("%s at %d", *ternos[i], i+1))
			t := time.Now()
			respBytes, err = dataSource.LastCan(*ternos[i], "2018-01-01 00:00:00", "2020-12-20 23:59:59", 1)
			elapsed := time.Since(t)
			fmt.Println(" Api request:", elapsed)

			if err != nil {
				fmt.Println(fmt.Sprintf("%s", err))
				fmt.Println(fmt.Sprintf("6秒后重试"))
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
			log.Fatal(err1)
		}
		elapsed := time.Since(t)
		fmt.Println(" json parse:", elapsed)

		size := len(pageInfo.Rows)
		fmt.Println(fmt.Sprintf(" total %d", size))

		if size == 0 {
			// 统计记录
			statTxt, err := os.OpenFile(statfname, os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				log.Fatal(err)
			}
			wstat := csv.NewWriter(statTxt)
			wstat.Write([]string{*ternos[i], strconv.Itoa(0), strconv.Itoa(fileNum)})
			wstat.Flush()
			fmt.Println()
			continue
		}

		// 形成二维数组
		buffer := [][]string{}
		for j := 0; j < len(pageInfo.Rows); j++ {

			gpsDate, _ := time.Parse("20060102150405", pageInfo.Rows[j].GPSDateTime)
			time, _ := time.Parse("20060102150405", pageInfo.Rows[j].TIME)

			//0 GPS已定位；1 GPS未定位；2非差分定位; 3 差分定位; 4正在估算

			gpsStatus := pageInfo.Rows[j].Status
			if gpsStatus == "0" {
				gpsStatus = "GPS已定位"
			}

			if gpsStatus == "1" {
				gpsStatus = "GPS未定位"
			}

			if gpsStatus == "2" {
				gpsStatus = "2非差分定位"
			}

			if gpsStatus == "3" {
				gpsStatus = "差分定位"
			}

			if gpsStatus == "4" {
				gpsStatus = "4正在估算"
			}

			data := []string{*ternos[i],
				pageInfo.Rows[j].SoftwareVersion,
				pageInfo.Rows[j].HardwareVersion,
				pageInfo.Rows[j].PowerSupplyVolt,
				pageInfo.Rows[j].InputPowerVoltage,
				pageInfo.Rows[j].InputBatteryVoltage,
				pageInfo.Rows[j].Longitude,
				pageInfo.Rows[j].Latitude,
				gpsStatus,
				gpsDate.Format("2006/01/02 15:04:05"),
				time.Format("2006/01/02 15:04:05")}

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
			csvName := fmt.Sprintf("%s_%d.csv", datafname, lines/100000)

			if fileNum != (lines / 100000) {
				affectRows, err := createFile(csvName, title)

				if err != nil {
					log.Fatal(err)
				}

				fileNum = lines / 100000
				lines += affectRows
			}

			txt, err := os.OpenFile(csvName, os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				log.Fatal(err)
			}
			w := csv.NewWriter(txt)
			w.WriteAll(buf[k])
			w.Flush()

			fmt.Println(fmt.Sprintf(" total size at index of %d, %d", k, len(buf[k])))
			lines = lines + len(buf[k])

			total += len(buf[k])
			index = append(index, strconv.Itoa(fileNum))
		}

		// 统计记录
		statTxt, err := os.OpenFile(statfname, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}
		wstat := csv.NewWriter(statTxt)
		wstat.Write([]string{*ternos[i], strconv.Itoa(total), strings.Join(index, ",")})
		wstat.Flush()
		statTxt.Close()

		fmt.Println()
	}

	return nil
}

func exportCsv(ternos []*string, statfname string, datafname string, lines int, fileNum int) error {

	title := []string{"设备唯一识别号", "GNSS时间", "经度", "纬度", "GNSS速度", "GNSS状态", "海拔高度",
		"水平精度因子", "正在使用的卫星数量", "发动机转速", "机油压力", "发动机工作时间", "油耗消耗总量",
		"每小时油耗", "行驶总里程"}

	for i := 0; i < len(ternos); i++ {

		// time.Sleep(6 * time.Second)

		var respBytes []byte
		var err error
		for retryIdx := 0; retryIdx < 5; retryIdx++ {
			fmt.Println(fmt.Sprintf("%s at %d", *ternos[i], i+1))
			t := time.Now()
			respBytes, err = dataSource.ListCanHistory(*ternos[i], "2019-01-01 00:00:00", "2019-12-20 23:59:59", 10000)
			elapsed := time.Since(t)
			fmt.Println(" Api request:", elapsed)

			if err != nil {
				fmt.Println(fmt.Sprintf("%s", err))
				fmt.Println(fmt.Sprintf("6秒后重试"))
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
			log.Fatal(err1)
		}
		elapsed := time.Since(t)
		fmt.Println(" json parse:", elapsed)

		size := len(pageInfo.Rows)
		fmt.Println(fmt.Sprintf(" total %d", size))

		if size == 0 {
			// 统计记录
			statTxt, err := os.OpenFile(statfname, os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				log.Fatal(err)
			}
			wstat := csv.NewWriter(statTxt)
			wstat.Write([]string{*ternos[i], strconv.Itoa(0), strconv.Itoa(fileNum)})
			wstat.Flush()
			fmt.Println()
			continue
		}

		// 形成二维数组
		buffer := [][]string{}
		for j := 0; j < len(pageInfo.Rows); j++ {

			gpsDate, _ := time.Parse("20060102150405", pageInfo.Rows[j].Can.GnssTime)

			ewt, err := strconv.ParseFloat(pageInfo.Rows[j].Can.EngineWorkTime, 64)
			if err != nil {
				ewt = 0
			}

			ewt1, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", ewt), 64)

			op, err := strconv.ParseFloat(pageInfo.Rows[j].Can.OilPressure, 64)
			if err != nil {
				op = 0
			}
			op1, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", op/1000), 64)

			gpsStatus := pageInfo.Rows[j].Can.Status
			if gpsStatus == "0" {
				gpsStatus = "2"
			}

			if gpsStatus == "1" {
				gpsStatus = "0"
			}

			if gpsStatus == "2" {
				gpsStatus = "1"
			}

			if gpsStatus == "3" {
				gpsStatus = "2"
			}

			if gpsStatus == "4" {
				gpsStatus = "6"
			}

			al, err := strconv.ParseFloat(pageInfo.Rows[j].Can.Altitude, 64)
			if err != nil {
				al = 0
			}
			al1, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", al), 64)

			gwsz, err := strconv.ParseFloat(pageInfo.Rows[j].Can.GWSZ, 64)
			if err != nil {
				gwsz = 0
			}

			gwsz1, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", gwsz), 64)

			data := []string{*ternos[i],
				gpsDate.Format("2006-01-02 15:04:05"),
				pageInfo.Rows[j].Can.Longitude,
				pageInfo.Rows[j].Can.Latitude,
				pageInfo.Rows[j].Can.Speed,
				gpsStatus,
				strconv.FormatFloat(al1, 'f', 0, 64),
				pageInfo.Rows[j].Can.Hdop,
				pageInfo.Rows[j].Can.UsedSatelliteNumber,
				pageInfo.Rows[j].Can.EngineSpeed,
				strconv.FormatFloat(op1, 'f', 2, 64),
				strconv.FormatFloat(ewt1, 'f', 1, 64),
				pageInfo.Rows[j].Can.TotalFuelConsumption,
				pageInfo.Rows[j].Can.FuelConsumptionPerHour,
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
			csvName := fmt.Sprintf("%s_%d.csv", datafname, lines/100000)

			if fileNum != (lines / 100000) {
				affectRows, err := createFile(csvName, title)

				if err != nil {
					log.Fatal(err)
				}

				fileNum = lines / 100000
				lines += affectRows
			}

			txt, err := os.OpenFile(csvName, os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				log.Fatal(err)
			}
			w := csv.NewWriter(txt)
			w.WriteAll(buf[k])
			w.Flush()

			fmt.Println(fmt.Sprintf(" total size at index of %d, %d", k, len(buf[k])))
			lines = lines + len(buf[k])

			total += len(buf[k])
			index = append(index, strconv.Itoa(fileNum))
		}

		// 统计记录
		statTxt, err := os.OpenFile(statfname, os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}
		wstat := csv.NewWriter(statTxt)
		wstat.Write([]string{*ternos[i], strconv.Itoa(total), strings.Join(index, ",")})
		wstat.Flush()
		statTxt.Close()

		fmt.Println()
	}

	return nil
}

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

		w := csv.NewWriter(file)
		w.Write(title)
		w.Flush()

		return 1, nil
	}
	return 0, err
}
