package model

type Prop struct {
	Id            string `xml:"id,attr"`
	En            string `xml:"en"`
	Cn            string `xml:"cn"`
	Unit          string `xml:"unit"`
	Base64        bool   `xml:"base64"`
	ConvertResult bool   `xml:"convertResult"`
}

type CanXml struct {
	Props []Prop `xml:"prop"`
}
