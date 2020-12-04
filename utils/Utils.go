package utils

import (
	"ConfigurationTools/configurationManager"
	"strings"
	"syscall"
)

type KeyValuePair struct {
	Code int
	Name string
}

type KeyValuePair2 struct {
	Key   string
	Value string
}

func KnownConfigUrl() []*KeyValuePair2 {
	var model = []*KeyValuePair2{}

	if strings.TrimSpace(configurationManager.CANConfig.TestUrl) != "" {
		model = append(model, &KeyValuePair2{"测试环境CAN字段定义", configurationManager.CANConfig.TestUrl})
	}

	if strings.TrimSpace(configurationManager.CANConfig.ProUrl) != "" {
		model = append(model, &KeyValuePair2{"生产环境CAN字段定义", configurationManager.CANConfig.ProUrl})
	}

	return model
}

func KnownDataType() []*KeyValuePair {
	return []*KeyValuePair{
		{1, "日期时间"},
		{2, "数字枚举"},
		{3, "数据"},
		{5, "文本枚举"},
		{6, "文本多枚举"},
		{7, "多字段组合枚举"},
		{4, "其他"},
	}
}

func ToDataType(code int) string {
	switch code {
	case 1:
		return "日期时间"
	case 2:
		return "数字枚举"
	case 3:
		return "数据"
	case 4:
		return "其他"
	case 5:
		return "文本枚举"
	case 6:
		return "文本多枚举"
	case 7:
		return "多字段组合枚举"
	default:
		return ""
	}
}

func KnownAlarm() []*KeyValuePair {
	return []*KeyValuePair{
		{1, "报警项"},
		{0, "非报警项"},
	}
}

func ToAlarm(code int) string {
	switch code {
	case 1: //报警项
		return "√"
	case 0: //非报警项
		return ""
	default:
		return ""
	}
}

func KnownAnaly() []*KeyValuePair {
	return []*KeyValuePair{
		{1, "可分析项"},
		{2, "默认分析项"},
		{0, "不可分析项"},
	}
}

func ToAlayly(code int) string {
	switch code {
	case 1: //可分析项
		return "√"
	case 2: //默认分析项
		return "√√"
	case 0: //不可分析项
		return ""
	default:
		return ""
	}
}

func GetSystemMetrics(nIndex int) int {
	ret, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(nIndex))
	return int(ret)
}
