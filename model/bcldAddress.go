package model

type BcldAddress struct {
	Status int
	Data   []AddressData
}

type AddressData struct {
	Key   string
	Value string
}
