package configurationManager

import (
	"ConfigurationTools/mynotify"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var jsonData map[string]interface{} = make(map[string]interface{})

func initJSON() {
	bytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		mynotify.Error("读取配置文件失败：" + err.Error())
		log.Fatal(err)
		return
	}

	configStr := string(bytes[:])
	reg := regexp.MustCompile(`/\*.*\*/`)

	configStr = reg.ReplaceAllString(configStr, "")
	bytes = []byte(configStr)
	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		mynotify.Error("解析配置文件失败：" + err.Error())
		log.Fatal(err)
		return
	}
}

// 数据库链接字符串
type dbConfig struct {
	Name    string
	Code    int
	ConnStr string
	Mode    string
}

var DBConfig []dbConfig = []dbConfig{}

func initDB() {
	dbs := jsonData["connectionStrings"].([]interface{})
	for i := 0; i < len(dbs); i++ {
		item := dbs[i].(map[string]interface{})

		temp := dbConfig{}
		if connName, isOk := item["name"]; isOk {
			if val, isOk := connName.(string); isOk {
				temp.Name = val
			}
			if val, isOk := connName.(float64); isOk {
				temp.Name = strconv.FormatFloat(val, 'f', 0, 64)
			}
		}

		if code, isoK := item["code"]; isoK {
			if val, isOk := code.(string); isOk {
				i, err := strconv.Atoi(val)
				if err == nil {
					temp.Code = i
				} else {
					temp.Code = -1
				}
			}

			if val, isOk := code.(float64); isOk {
				istr := strconv.FormatFloat(val, 'f', 0, 64)
				i, err := strconv.Atoi(istr)
				if err == nil {
					temp.Code = i
				}
			}
		}

		if mode, isoK := item["mode"]; isoK {
			if val, isOk := mode.(string); isOk {
				temp.Mode = val
			}
			if val, isOk := mode.(float64); isOk {
				temp.Mode = strconv.FormatFloat(val, 'f', 0, 64)
			}
		}

		if connStr, isoK := item["connectionString"]; isoK {
			if val, isOk := connStr.(string); isOk {
				temp.ConnStr = val
			}
			if val, isOk := connStr.(float64); isOk {
				temp.ConnStr = strconv.FormatFloat(val, 'f', 0, 64)
			}
		}

		DBConfig = append(DBConfig, temp)
	}
}

// can xml链接
type canConfig struct {
	Name    string
	TestUrl string
	ProUrl  string
}

var CANConfig canConfig = canConfig{}

func initCanXml() {
	cfg := jsonData["canConfigXml"].(map[string]interface{})

	temp := canConfig{}
	if name, isOk := cfg["name"]; isOk {
		if val, isOk := name.(string); isOk {
			temp.Name = val
		}
		if val, isOk := name.(float64); isOk {
			temp.Name = strconv.FormatFloat(val, 'f', 0, 64)
		}
	}

	if testUrl, isoK := cfg["test"]; isoK {
		if val, isOk := testUrl.(string); isOk {
			temp.TestUrl = val
		}
		if val, isOk := testUrl.(float64); isOk {
			temp.TestUrl = strconv.FormatFloat(val, 'f', 0, 64)
		}
	}

	if proUrl, isoK := cfg["production"]; isoK {
		if val, isOk := proUrl.(string); isOk {
			temp.ProUrl = val
		}
		if val, isOk := proUrl.(float64); isOk {
			temp.ProUrl = strconv.FormatFloat(val, 'f', 0, 64)
		}
	}

	CANConfig = temp
}

// can分组编码
type canGroup struct {
	GroupName string
	Code      int
}

var CanGroup []canGroup = []canGroup{}
var CanGroupMap map[int]string = make(map[int]string)

func initCanGroup() {
	group := jsonData["canGroup"].([]interface{})
	for i := 0; i < len(group); i++ {
		item := group[i].(map[string]interface{})

		temp := canGroup{}

		if name, isOk := item["Name"]; isOk {
			if val, isOk := name.(string); isOk {
				temp.GroupName = val
			}
			if val, isOk := name.(float64); isOk {
				temp.GroupName = strconv.FormatFloat(val, 'f', 0, 64)
			}
		}

		if code, isoK := item["Code"]; isoK {
			if val, isOk := code.(string); isOk {
				i, err := strconv.Atoi(val)
				if err == nil {
					temp.Code = i
				} else {
					temp.Code = -1
				}
			}

			if val, isOk := code.(float64); isOk {
				istr := strconv.FormatFloat(val, 'f', 0, 64)
				i, err := strconv.Atoi(istr)
				if err == nil {
					temp.Code = i
				}
			}
		}

		if temp.Code != -1 {
			CanGroup = append(CanGroup, temp)
			CanGroupMap[temp.Code] = temp.GroupName
		}
	}
}

type appSetting struct {
	Key       string
	TestValue string
	ProValue  string
}

var appSettings map[string]string = make(map[string]string)

func AppSetting(key string) string {
	mode, err := GetMode()
	if err != nil {
		return ""
	}
	return appSettings[fmt.Sprintf("%s_%s", key, strings.ToLower(mode))]
}

func initAppSetting() {
	setting := jsonData["appSettings"].([]interface{})
	for i := 0; i < len(setting); i++ {
		item := setting[i].(map[string]interface{})
		temp := appSetting{}

		if name, isOk := item["Key"]; isOk {
			if val, isOk := name.(string); isOk {
				temp.Key = val
			}
			if val, isOk := name.(float64); isOk {
				temp.Key = strconv.FormatFloat(val, 'f', 0, 64)
			}
		}

		if strings.Trim(temp.Key, " ") == "" {
			continue
		}

		if devValue, isoK := item["test"]; isoK {
			if val, isOk := devValue.(string); isOk {
				temp.TestValue = val
			}

			if val, isOk := devValue.(float64); isOk {
				temp.TestValue = strconv.FormatFloat(val, 'f', 0, 64)
			}

			appSettings[fmt.Sprintf("%s_%s", temp.Key, "test")] = temp.TestValue
		}

		if proValue, isoK := item["production"]; isoK {
			if val, isOk := proValue.(string); isOk {
				temp.ProValue = val
			}

			if val, isOk := proValue.(float64); isOk {
				temp.ProValue = strconv.FormatFloat(val, 'f', 0, 64)
			}

			appSettings[fmt.Sprintf("%s_%s", temp.Key, "production")] = temp.ProValue
		}
	}
}

func init() {
	initJSON()

	initDB()
	initCanXml()
	initCanGroup()
	initAppSetting()
}
