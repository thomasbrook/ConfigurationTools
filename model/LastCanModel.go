package model

type LastCanModel struct {
	//终端编号
	did string

	TIME string `json:"TIME"`
	//GPS时间
	GPSDateTime string `json:"GPSDateTime"`
	// 经度
	Longitude string `json:"GPSLon"`
	// 纬度
	Latitude string `json:"GPSLat"`
	// GNSS状态
	Status string `json:"GPSLocationStatus"`
	// 供电电压
	InputPowerVoltage string `json:"inputPowerVoltage"`
	// 备用电池电压
	InputBatteryVoltage string `json:"inputBatteryVoltage"`
	// 电源电压
	PowerSupplyVolt string `json:"powerSupplyVolt"`
	// 硬件版本
	HardwareVersion string `json:"hardwareVersion"`
	// 软件版本
	SoftwareVersion string `json:"softwareVersion"`
}

type LastRowModel struct {
	Rows []*LastCanModel `json:"list"`
}

type AnonymousCanModel struct {
	Can *LastCanModel `json:"can"`
}
