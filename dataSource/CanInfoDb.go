package dataSource

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// ListVehicle 加载车辆类型下的终端列表
func ListVehicle(vehicleTypeID string, processedTers []string) ([]*string, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	whereSQL := []string{}
	if len(processedTers) > 0 {
		whereSQL = append(whereSQL, fmt.Sprintf("ter.terminal_no NOT IN ('%s')", strings.Join(processedTers, "','")))
	}

	lineOper := ""
	if len(whereSQL) > 0 {
		lineOper = " AND "
	}

	sql := fmt.Sprintf(`SELECT 
							ter.terminal_no
						FROM
							biz_vehicle_info v
								INNER JOIN
							biz_terminal_info ter ON v.id = ter.vehicle_id
						WHERE
							v.vehicle_type_id = ? %s %s
						ORDER BY ter.terminal_no ASC `, lineOper, strings.Join(whereSQL, " AND "))
	stmt, err := db.Prepare(sql)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(vehicleTypeID)
	if err != nil {
		return nil, err
	}

	ters := []*string{}
	for rows.Next() {
		var ter *string
		rows.Scan(&ter)

		ters = append(ters, ter)
	}
	return ters, nil
}

// ListVehicleByOrg 查询机构下的终端列表
func ListVehicleByOrg(orgID string, processedTers []string) ([]*string, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	whereSQL := []string{}
	if len(processedTers) > 0 {
		whereSQL = append(whereSQL, fmt.Sprintf("ter.terminal_no NOT IN ('%s')", strings.Join(processedTers, "','")))
	}

	lineOper := ""
	if len(whereSQL) > 0 {
		lineOper = " AND "
	}

	sql := fmt.Sprintf(`SELECT 
							ter.terminal_no
						FROM
							biz_terminal_info ter
						WHERE
							ter.org_id = ? %s %s
						ORDER BY ter.terminal_no ASC `, lineOper, strings.Join(whereSQL, " AND "))
	stmt, err := db.Prepare(sql)

	fmt.Println(sql)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(orgID)
	if err != nil {
		return nil, err
	}

	ters := []*string{}
	for rows.Next() {
		var ter *string
		rows.Scan(&ter)

		ters = append(ters, ter)
	}
	return ters, nil
}

func LastCan(terminalNo string, startDate string, endDate string, pageSize int) ([]byte, error) {

	// 哈希参数
	param := make(map[string]interface{})
	param["did"] = terminalNo
	param["showColumns"] = "3006,softwareVersion,3007,hardwareVersion,inputPowerVoltage,inputBatteryVoltage,powerSupplyVolt,GPSLon,GPSLat,GPSLocationStatus,TIME,GPSDateTime"
	param["cascade"] = true

	bytesData, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bytesData)
	//http://172.16.1.245:30000
	req, err := http.NewRequest("POST", "http://172.16.1.245:8086/basicInfo/getLastestInfos", reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json;charset=UTF-8")

	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*60) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 60)) //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 60,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

// ListCanHistory 获取终端CAN历史信息
func ListCanHistory(terminalNo string, startDate string, endDate string, pageSize int) ([]byte, error) {

	// 哈希参数
	param := make(map[string]interface{})
	param["did"] = terminalNo
	param["startTime"] = startDate
	param["endTime"] = endDate
	param["reqColumns"] = "GPSLon,GPSLat"
	param["showColumns"] = "GPSDateTime,GPSLon,GPSLat,gpsSpeed,GPSLocationStatus,altitude,horizontalDilutionOfPrecision,numOfUsedSatellites,engineSpeed,relativeOilPressure,engineWorkTime,totalFuelConsumption,instanFuel,totalMileage"
	param["pageNum"] = 1
	param["pageSize"] = pageSize
	// param["filterIfMissing"] = false
	param["reversed"] = true
	param["condition"] = "GPSLat > 0 && GPSLon > 0"

	bytesData, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bytesData)
	//http://172.16.1.245:30000
	req, err := http.NewRequest("POST", "http://172.16.1.245:30000/spring-cloud-server-nrmv/nrmv/getHistoryInfo2", reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", "application/json;charset=UTF-8")

	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*60) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 60)) //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 60,
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}
