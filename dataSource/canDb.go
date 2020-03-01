package dataSource

import (
	"ConfigurationTools/model"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

func OpenDB() (db *sql.DB, err error) {
	db, err = sql.Open("mysql", model.QZ_Mysql_driver)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func BatchInsertCan(groupInfoId string, can []*model.CanConfig) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT biz_outfield_group_detail 
								SET id=?,outfield_id=?,unit=?,outfield_sn=?,group_info_id=?,
								chinesename=?,formula=?,data_type=?,field_name=?,decimals=?,
								isAlarm=?,data_map=?,is_analysable=?`)
	if err != nil {
		return err
	}

	for i := 0; i < len(can); i++ {
		uid,_ := uuid.NewV4()

		res, err := stmt.Exec(strings.ReplaceAll(uid.String(), "-", ""),
			can[i].OutfieldId,
			can[i].Unit,
			can[i].OutfieldSn,
			groupInfoId,
			can[i].Chinesename,
			can[i].Formula,
			can[i].DataType,
			can[i].FieldName,
			can[i].Decimals,
			can[i].IsAlarm,
			can[i].DataMap,
			can[i].IsAnalysable)
		if err != nil {
			return err
		}

		_, err1 := res.LastInsertId()
		if err1 != nil {
			return err
		}
	}

	return nil
}

func ListVehicleType(searckKey string) (vtype []*model.VehicleType, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlBuffer := []string{}
	sqlBuffer = append(sqlBuffer, `SELECT 
				    vt.type_id,
				    vt.type_name,
				    org.org_name,
 						    ifnull(g.id,''),
				    ifnull(g.remark,''),
				    COUNT(gd.id)
				FROM
				    biz_vehicle_type vt
				        INNER JOIN
				    sys_organization org ON vt.org_id = org.id AND vt.in_use = 1
				        LEFT JOIN
				    biz_outfield_group_info g ON vt.type_id = g.agreement_id
				        LEFT JOIN
				    biz_outfield_group_detail gd ON g.id = gd.group_info_id
				`)

	if strings.TrimSpace(searckKey) != "" {
		sqlBuffer = append(sqlBuffer, "WHERE vt.type_name LIKE ? OR org.org_name LIKE ?")
	}

	sqlBuffer = append(sqlBuffer, ` GROUP BY vt.type_id , vt.type_name , org.org_name , g.id , g.remark
				ORDER BY vt.type_id , g.remark ASC`)

	stmt, err := db.Prepare(strings.Join(sqlBuffer, " "))
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	if strings.TrimSpace(searckKey) != "" {
		rows, err = stmt.Query("%"+searckKey+"%", "%"+searckKey+"%")
	} else {
		rows, err = stmt.Query()
	}

	if err != nil {
		return nil, err
	}

	vtype = []*model.VehicleType{}

	var typeId string
	var typeName string
	var orgName string
	var id string
	var remark int
	var count int

	vt := new(model.VehicleType)
	for rows.Next() {

		rows.Scan(&typeId, &typeName, &orgName, &id, &remark, &count)

		if strings.TrimSpace(vt.TypeId) == "" {
			vt.TypeId = typeId
			vt.TypeName = typeName
			vt.OrgName = orgName
			if remark == 0 {
				vt.GeneralId = id
				vt.GeneralInfoCount = count
			}

			if remark == 1 {
				vt.CanId = id
				vt.CanCount = count
			}

			if remark == 2 {
				vt.CmdId = id
				vt.CmdCount = count
			}
		} else if strings.TrimSpace(vt.TypeId) == typeId {

			if remark == 0 {
				vt.GeneralId = id
				vt.GeneralInfoCount = count
			}

			if remark == 1 {
				vt.CanId = id
				vt.CanCount = count
			}

			if remark == 2 {
				vt.CmdId = id
				vt.CmdCount = count
			}
		} else {
			vtype = append(vtype, vt)

			vt = new(model.VehicleType)
			vt.TypeId = typeId
			vt.TypeName = typeName
			vt.OrgName = orgName

			if remark == 0 {
				vt.GeneralId = id
				vt.GeneralInfoCount = count
			}

			if remark == 1 {
				vt.CanId = id
				vt.CanCount = count
			}

			if remark == 2 {
				vt.CmdId = id
				vt.CmdCount = count
			}
		}
	}

	if strings.TrimSpace(vt.TypeId) != "" {
		vtype = append(vtype, vt)
	}

	return vtype, nil
}

func ListCan(vtid string, remark int) (cans []*model.CanConfig, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT 
							    id,
							    outfield_id,
							    ifnull(unit,''),
							    ifnull(outfield_sn,''),
							    group_info_id,
							    ifnull(chinesename,''),
							    ifnull(formula,''),
							    ifnull(data_type,''),
							    ifnull(field_name,''),
							    ifnull(decimals,''),
							    ifnull(isAlarm,''),
							    ifnull(data_map,''),
							    is_analysable
							FROM
							    biz_outfield_group_detail
							WHERE
							    group_info_id IN (SELECT 
							            id
							        FROM
							            biz_outfield_group_info
							        WHERE
							            agreement_id = ? AND remark = ?)`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(vtid, remark)
	if err != nil {
		return nil, err
	}

	cans = []*model.CanConfig{}
	for rows.Next() {
		can := &model.CanConfig{}
		rows.Scan(&can.Id, &can.OutfieldId, &can.Unit, &can.OutfieldSn, &can.GroupInfoId, &can.Chinesename, &can.Formula,
			&can.DataType, &can.FieldName, &can.Decimals, &can.IsAlarm, &can.DataMap, &can.IsAnalysable)

		cans = append(cans, can)
	}

	return cans, nil
}

func DeleteCan(id string) (int64, error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM biz_outfield_group_detail WHERE id =?`)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func UpdateCan(can *model.CanConfig) (int64, error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE biz_outfield_group_detail
							SET	unit = ?,
								outfield_sn = ?,
								chinesename = ?,
								formula = ?,
								data_type = ?,
								field_name = ?,
								decimals = ?,
								isAlarm = ?,
								data_map = ?,
								is_analysable = ?
							WHERE id = ? `)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(can.Unit, can.OutfieldSn, can.Chinesename, can.Formula,
		can.DataType, can.FieldName, can.Decimals, can.IsAlarm, can.DataMap, can.IsAnalysable, can.Id)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
