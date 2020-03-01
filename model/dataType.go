package model

type DataType int

const (
	NullValue    DataType = 0
	DateValue    DataType = 1
	EnumValue    DataType = 2
	NumericValue DataType = 3
	OtherValue   DataType = 4
	EnumText     DataType = 5

	//QZ_Mysql_driver = "root:l520m6ti@tcp(192.168.11.102:3306)/agfpiddb?charset=utf8"
	QZ_Mysql_driver = "aguser:z3ZelaJc@tcp(172.16.1.27:3306)/agfpiddb?charset=utf8"

	//CONFIG_URL = "http://192.168.11.8/config.xml"
	CONFIG_URL = "http://172.16.1.245:8000/config.xml"
)
