package model

type BaiduAddress struct {
	Status  int
	Result  AddressModel
	Message string
}

type AddressModel struct {
	AddressComponent AddressComponentModel

	// 当前位置结合POI的语义化结果描述
	Sematic_Description string
}

type AddressComponentModel struct {
	Province string
	City     string
	District string
	Town     string
}
