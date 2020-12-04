package model

type CanDetail struct {
	Index        int
	Id           string
	OutfieldId   string
	Unit         string
	OutfieldSn   float64
	GroupInfoId  string
	GroupName    string
	OrgName      string
	Chinesename  string
	Formula      string
	DataType     string
	FieldName    string
	Decimals     string
	DataMap      string
	IsAlarm      string
	IsAnalysable int
	IsDelete     int
	Checked      bool
	Note         string
}

// 用于Table数据填充
type CanDetailTableAdapter struct {
	Index        int
	Id           string
	Key          string
	Unit         string
	Sort         float64
	GroupId      string
	Alias        string
	Formula      string
	DataType     string
	FieldName    string
	Prec         string
	DataScope    string
	IsAlarm      string
	IsAnalysable string
	IsDelete     string
	Note         string
	Checked      bool
}

type CanDetailWithGroup struct {
	CanDetail

	GroupName string
	Sort      int
	Remark    string
}
