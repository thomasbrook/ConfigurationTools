package model

type CanInfoModel struct {
	//终端编号
	did string
	//GPS时间
	GnssTime string `json:"GPSDateTime"`
	// 经度
	Longitude string `json:"GPSLon"`
	// 纬度
	Latitude string `json:"GPSLat"`
	// GNSS速度
	Speed string `json:"gpsSpeed"`
	// GNSS状态
	Status string `json:"GPSLocationStatus"`
	// 海拔高度
	Altitude string `json:"altitude"`
	// 水平精度因子
	Hdop string `json:"horizontalDilutionOfPrecision"`
	//使用的卫星数量
	UsedSatelliteNumber string `json:"numOfUsedSatellites"`
	// 发动机转速
	EngineSpeed string `json:"engineSpeed"`
	// 机油压力
	OilPressure string `json:"relativeOilPressure"`
	// 发动机工作时间
	EngineWorkTime string `json:"engineWorkTime"`
	// 总油耗量
	TotalFuelConsumption string `json:"totalFuelConsumption"`
	// 每小时油耗
	FuelConsumptionPerHour string `json:"instanFuel"`
	// 行驶总里程
	GWSZ string `json:"totalMileage"`
}

type PageInfoModel struct {
	Rows []*AnonymousModel `json:"rows"`
}

type AnonymousModel struct {
	Can *CanInfoModel `json:"can"`
}
