package model

type DataType int

const (
	NullValue    DataType = 0
	DateValue    DataType = 1
	EnumValue    DataType = 2
	NumericValue DataType = 3
	OtherValue   DataType = 4
	EnumText     DataType = 5
)
