package model

type VehicleType struct {
	TypeId   string
	TypeName string
	OrgName  string

	GeneralId        string
	GeneralInfoCount int
	HasGeneral       bool

	CanId    string
	CanCount int
	HasCan   bool

	CmdId    string
	CmdCount int
	HasCmd   bool
}
