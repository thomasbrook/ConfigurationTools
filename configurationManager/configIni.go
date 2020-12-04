package configurationManager

import (
	"ConfigurationTools/mynotify"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

var iniData map[string]interface{}

type IniKeyType string

const (
	DATABASE         IniKeyType = "database"
	DATABASE_CONNSTR IniKeyType = "database_connstr"
	Model            IniKeyType = "mode"
)

func (key IniKeyType) String() string {
	switch key {
	case DATABASE:
		return "database"
	case DATABASE_CONNSTR:
		return "database_connstr"
	case Model:
		return "mode"
	default:
		return "unknown"
	}
}

// 加载初始化键值对
func initIni() {
	inif, err := os.Open("./.ini")
	if err != nil {
		return
	}

	defer inif.Close()

	iniData = make(map[string]interface{})
	r := bufio.NewReader(inif)

	for {
		// 分行读取文件，返回单行，不包括行尾字节（\r\n or \n）
		data, _, err := r.ReadLine()

		if err == io.EOF {
			break
		}

		if err != nil {
			mynotify.Error("读取初始化文件失败：" + err.Error())
			break
		}

		_data := string(data)
		if strings.Trim(string(_data), " ") == "" {
			continue
		}

		kv := strings.SplitN(_data, "=", 2)
		if len(kv) != 2 {
			continue
		}

		iniData[strings.ToLower(strings.Trim(kv[0], " "))] = strings.Trim(kv[1], " ")
	}
}

// 根据键，获取值
func GetIni(key IniKeyType) (string, error) {

	if val, isExist := iniData[key.String()]; isExist {
		if v, isOk := val.(string); isOk {
			return v, nil
		}
		if v, isOk := val.(float64); isOk {
			return strconv.FormatFloat(v, 'f', 0, 64), nil
		}
	}

	return "", errors.New("数据不存在")
}

// 添加新键值对
func SetIni(key IniKeyType, val string) error {

	val = strings.Trim(val, " ")

	if iniData == nil {
		iniData = make(map[string]interface{})
	}

	iniData[key.String()] = val

	inif, err := os.OpenFile("./.ini", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer inif.Close()

	w := bufio.NewWriter(inif)

	var keys []string
	for k := range iniData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		lineStr := fmt.Sprintf("%s=%s", k, iniData[k])
		fmt.Fprintln(w, lineStr)
	}
	return w.Flush()
}

// 数据库链接字符串
func GetDatabaseConnStr() (string, error) {
	//return DBConfig[1].ConnectionString, nil

	if val, isExist := iniData[DATABASE_CONNSTR.String()]; isExist {

		if connstr, isOk := val.(string); isOk {
			return connstr, nil
		}

		if connstr, isOk := val.(float64); isOk {
			return strconv.FormatFloat(connstr, 'f', 0, 64), nil
		}
	}

	return "", errors.New("数据不存在")
}

func GetMode() (string, error) {
	if val, isExist := iniData[Model.String()]; isExist {

		if modeStr, isOk := val.(string); isOk {
			return modeStr, nil
		}

		if modeStr, isOk := val.(float64); isOk {
			return strconv.FormatFloat(modeStr, 'f', 0, 64), nil
		}
	}

	return "", errors.New("数据不存在")
}

func init() {
	initIni()
}
