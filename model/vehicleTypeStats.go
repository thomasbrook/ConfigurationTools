package model

type VehicleTypeStats struct {
	TypeId   string
	TypeName string
	OrgId    string
	OrgName  string

	Group []*GroupStats
}

type GroupStats struct {
	Name  string
	Code  int
	Id    string
	Count int
}
