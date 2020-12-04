package dataSource

import (
	"ConfigurationTools/configurationManager"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

func GetLastestCan(terminalNo string) ([]byte, error) {

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
	req, err := http.NewRequest("POST", configurationManager.AppSetting("getLastestInfos"), reader)
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
	req, err := http.NewRequest("POST", configurationManager.AppSetting("getHistoryInfo2"), reader)
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

func GetAddressByCoordinate(lat float64, lng float64) ([]byte, error) {

	var _lat = strconv.FormatFloat(lat, 'f', -1, 64)
	var _lng = strconv.FormatFloat(lng, 'f', -1, 64)

	//param := []string{}
	//param = append(param, fmt.Sprintf("location=%s,%s", _lat, _lng), "coordtype=wgs84ll", "output=json", "extensions_road=true", "extensions_town=true", "latest_admin=1", "ak=Q5tK3scFKg5zpi6iblnNuy5h7p8u3f6Y")

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

	//url := fmt.Sprintf("http://api.map.baidu.com/geocoder/v2/?%s", strings.Join(param, "&"))
	url := fmt.Sprintf("%s?latlng=%s,%s", configurationManager.AppSetting("ListAddressByCoordinate"), _lat, _lng)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}
