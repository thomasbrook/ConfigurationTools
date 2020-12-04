package dataSource

import (
	"ConfigurationTools/model"
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type VehicleTypeEntity struct {
	TypeId                string
	TypeName              string
	OrgId                 string
	OrgName               string
	IsIntelligent         int
	IsFilterMissingColumn int

	CanGroup []*CanGroupEntity
}

type CanGroupEntity struct {
	Id   string
	Name string
	Code int
	Sort float64
}

// GetVehicleType 查询单个车辆类型信息，包括基本信息及各个分组字段数量统计
func (vt *VehicleTypeEntity) GetVehicleType() (*VehicleTypeEntity, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	vtype := &VehicleTypeEntity{}
	err = db.QueryRow(`SELECT 
										vt.type_id,
										vt.type_name,
										org.org_name,
										org.id,
									 	ifnull(vt.is_intelligent, 0),
										ifnull(vt.is_filter_missing_column, 0)
									FROM
										biz_vehicle_type vt
											INNER JOIN
										sys_organization org ON vt.org_id = org.id
									WHERE
										vt.type_id = ?`, vt.TypeId).Scan(&vtype.TypeId, &vtype.TypeName, &vtype.OrgName, &vtype.OrgId, &vtype.IsIntelligent, &vtype.IsFilterMissingColumn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	if vtype.TypeId != "" {
		rows, err := db.Query(`SELECT 
							id, group_name, group_sn, remark
						FROM
							biz_outfield_group_info
						WHERE
							agreement_id = ?
						ORDER BY group_sn ASC`, vt.TypeId)
		if err != nil {
			return vtype, err
		}
		defer rows.Close()

		for rows.Next() {
			g := &CanGroupEntity{}
			rows.Scan(&g.Id, &g.Name, &g.Sort, &g.Code)
			vtype.CanGroup = append(vtype.CanGroup, g)
		}

		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}

	return vtype, nil
}

// ListCan 查询某个分组下的字段列表
func (vt *VehicleTypeEntity) ListCan(groupId string) (cans []*model.CanDetail, err error) {
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
							    is_analysable,
							    is_delete
							FROM
							    biz_outfield_group_detail
							WHERE
							    group_info_id = ?
							ORDER BY outfield_sn ASC`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(groupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cans = []*model.CanDetail{}
	idx := 0
	for rows.Next() {
		idx = idx + 1
		can := &model.CanDetail{Index: idx}
		rows.Scan(&can.Id, &can.OutfieldId, &can.Unit, &can.OutfieldSn, &can.GroupInfoId, &can.Chinesename, &can.Formula,
			&can.DataType, &can.FieldName, &can.Decimals, &can.IsAlarm, &can.DataMap, &can.IsAnalysable, &can.IsDelete)

		cans = append(cans, can)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cans, nil
}

// DeleteCan 删除某个字段
func (vt *VehicleTypeEntity) DeleteCan(id string) (int64, error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM biz_outfield_group_detail WHERE id =?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// UpdateCan 更新指定字段配置信息
func (vt *VehicleTypeEntity) UpdateCan(can *model.CanDetail) (int64, error) {
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
	defer stmt.Close()

	res, err := stmt.Exec(can.Unit, can.OutfieldSn, can.Chinesename, can.Formula,
		can.DataType, can.FieldName, can.Decimals, can.IsAlarm, can.DataMap, can.IsAnalysable, can.Id)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// BatchUpdateCanDetail 批量更新字段配置
func (vt *VehicleTypeEntity) BatchUpdateCanDetail(can []*model.CanDetail) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		db.Close()
		return err
	}

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare("UPDATE biz_outfield_group_detail SET unit = ?, outfield_id = ?, outfield_sn = ?, chinesename = ?, formula = ?, data_type = ?, decimals = ?, isAlarm = ?, data_map = ?, is_analysable = ? WHERE id =?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, v := range can {
		_, err1 := stmt.Exec(v.Unit, v.OutfieldId, v.OutfieldSn, v.Chinesename, v.Formula, v.DataType, v.Decimals, v.IsAlarm, v.DataMap, v.IsAnalysable, v.Id)
		if err1 != nil {
			err = err1
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	//stmt.Close()
	return nil
}

// BatchInsertCan 为某个分组信息，批量添加字段。（从大数据导入使用）
func (vt *VehicleTypeEntity) InsertCanFromBigdata(groupInfoId string, can []*model.CanDetail) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`INSERT biz_outfield_group_detail 
								SET id=?,outfield_id=?,unit=?,outfield_sn=?,group_info_id=?,
								chinesename=?,formula=?,data_type=?,field_name=?,decimals=?,
								isAlarm=?,data_map=?,is_analysable=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var sort int
	err = tx.QueryRow(`SELECT IFNULL(MAX(outfield_sn), 1) FROM biz_outfield_group_detail where group_info_id = ?`, groupInfoId).Scan(&sort)
	if err != nil {
		return err
	}

	for i := 0; i < len(can); i++ {
		uid, _ := uuid.NewV4()

		_, err1 := stmt.Exec(strings.ReplaceAll(uid.String(), "-", ""),
			can[i].OutfieldId,
			can[i].Unit,
			sort+1+i,
			groupInfoId,
			can[i].Chinesename,
			can[i].Formula,
			can[i].DataType,
			can[i].FieldName,
			can[i].Decimals,
			can[i].IsAlarm,
			can[i].DataMap,
			can[i].IsAnalysable)
		if err1 != nil {
			err = err1
			return err1
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ToggleDeleteStatus isDelete 0、未删除；1、已删除
func (vt *VehicleTypeEntity) DeleteCanField(id []string, isSoftDelete bool) (int, error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	whereSql := []string{}
	if len(id) > 0 {
		whereSql = append(whereSql, fmt.Sprintf(" id IN ('%s')", strings.Join(id, "','")))
	}

	linkOper := ""
	if len(whereSql) > 0 {
		linkOper = " AND "
	}

	var sql string
	if isSoftDelete {
		sql = fmt.Sprintf(`UPDATE biz_outfield_group_detail	SET is_delete = 1 WHERE 1=1 %s %s`, linkOper, strings.Join(whereSql, " AND "))
	} else {
		sql = fmt.Sprintf(`DELETE FROM biz_outfield_group_detail WHERE 1=1 %s %s`, linkOper, strings.Join(whereSql, " AND "))
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		return 0, err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowCount), nil
}

func (vt *VehicleTypeEntity) CancelDelete(id []string) (int, error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	whereSql := []string{}
	if len(id) > 0 {
		whereSql = append(whereSql, fmt.Sprintf(" id IN ('%s')", strings.Join(id, "','")))
	}

	linkOper := ""
	if len(whereSql) > 0 {
		linkOper = " AND "
	}

	sql := fmt.Sprintf(`UPDATE biz_outfield_group_detail SET is_delete = 0 WHERE 1=1 %s %s`, linkOper, strings.Join(whereSql, " AND "))

	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		return 0, err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowCount), nil
}

// CopyPasteCan 将某个分组下的字段，复制到当前分组下。 （从已存在车系导入）
// mode：1、如果字段已存在，更新该字段；2、如果字段已存在，跳过该字段
func (vt *VehicleTypeEntity) SyncCanFromExisting(groupId string, can []model.CanDetail, mode int) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		db.Close()
		return err
	}

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	// 查询目标分组内，已存在的字段
	var build strings.Builder
	for _, v := range can {
		build.WriteString(fmt.Sprintf("%s','", v.OutfieldId))
	}
	sql := fmt.Sprintf(`SELECT outfield_id FROM biz_outfield_group_detail WHERE group_info_id = ? AND outfield_id IN ('%s')`, build.String())
	rows, err := tx.Query(sql, groupId)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 将已存在的字段存储与哈希内
	existKeyMap := make(map[string]bool)
	for rows.Next() {
		var key string
		rows.Scan(&key)
		existKeyMap[key] = true
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	var sort int
	err = tx.QueryRow(`SELECT IFNULL(MAX(outfield_sn), 1) FROM biz_outfield_group_detail where group_info_id = ?`, groupId).Scan(&sort)
	if err != nil {
		return err
	}

	for idx, v := range can {

		_, isExist := existKeyMap[v.OutfieldId]
		if isExist {

			// 如果字段已存在，跳过该字段
			if mode == 2 {
				continue
			}

			// 如果字段已存在，更新该字段
			if mode == 1 {
				vt := new(model.CanDetail)
				err1 := tx.QueryRow(`SELECT 
										outfield_id,
							    		ifnull(unit,''),
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
										id = ?`, v.Id).Scan(&vt.OutfieldId, &vt.Unit, &vt.Chinesename, &vt.Formula, &vt.DataType, &vt.FieldName, &vt.Decimals, &vt.IsAlarm, &vt.DataMap, &vt.IsAnalysable)

				if err1 != nil {
					err = err1
					return err
				}

				_, err1 = tx.Exec(`UPDATE biz_outfield_group_detail 
											 SET unit  = ?,
												 chinesename  = ?,
												 formula  = ?,
												 data_type  = ?,
												 field_name  = ?,
												 decimals  = ?,
												 isAlarm  = ?,
												 data_map  = ?,
												 is_analysable  = ?
											WHERE group_info_id = ? AND outfield_id = ?`,
					vt.Unit, vt.Chinesename, vt.Formula, vt.DataType, vt.FieldName, vt.Decimals, vt.IsAlarm, vt.DataMap, vt.IsAnalysable, groupId, v.OutfieldId)

				if err1 != nil {
					err = err1
					return err
				}
			}
			continue
		}

		_, err1 := tx.Exec(`INSERT INTO biz_outfield_group_detail (id, outfield_id, unit, outfield_sn, group_info_id, chinesename, formula, data_type, field_name, decimals, isAlarm, data_map, is_analysable, is_delete) 
					   				   SELECT md5(uuid()), outfield_id, unit, ?, ?, chinesename, formula, data_type, field_name, decimals, isAlarm, data_map, is_analysable, 0 FROM biz_outfield_group_detail where id=?`,
			sort+1+idx, groupId, v.Id)

		if err1 != nil {
			err = err1
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ListAllCan 获取车型下的所有can
func (vt *VehicleTypeEntity) ListAllCan(vehicleTypeId string) ([]*model.CanDetailWithGroup, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT 
								IFNULL(g.group_name, ''),
								g.group_sn,
								IFNULL(g.remark, ''),
								IFNULL(d.outfield_id, ''),
								IFNULL(d.unit, ''),
								IFNULL(d.outfield_sn, ''),
								d.group_info_id,
								IFNULL(d.chinesename, ''),
								IFNULL(d.formula, ''),
								IFNULL(d.data_type, ''),
								IFNULL(d.field_name, ''),
								IFNULL(d.decimals, ''),
								IFNULL(d.isAlarm, ''),
								IFNULL(d.data_map, ''),
								d.is_analysable,
								d.is_delete
							FROM
								biz_outfield_group_detail d
									INNER JOIN
								biz_outfield_group_info g ON g.id = d.group_info_id
									AND g.agreement_id = ?
							ORDER BY g.group_sn ASC , d.outfield_sn ASC `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(vehicleTypeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cans := []*model.CanDetailWithGroup{}
	for rows.Next() {
		can := &model.CanDetailWithGroup{}
		rows.Scan(&can.GroupName, &can.Sort, &can.Remark, &can.OutfieldId, &can.Unit, &can.OutfieldSn, &can.GroupInfoId, &can.Chinesename, &can.Formula,
			&can.DataType, &can.FieldName, &can.Decimals, &can.IsAlarm, &can.DataMap, &can.IsAnalysable, &can.IsDelete)

		cans = append(cans, can)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cans, nil
}

// SynchCanFromCan 从csv文件导入CAN
func (vt *VehicleTypeEntity) SyncCanFromCsv(groupIds []string, can []*model.CanDetailWithGroup, mode int) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		db.Close()
		return err
	}

	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	// 查询已存在的 Can 信息
	var canKey []string
	for _, v := range can {
		canKey = append(canKey, v.OutfieldId)
	}

	sql := fmt.Sprintf(`SELECT 
									g.id, ifnull(g.remark,''), ifnull(gd.outfield_id,'')
								FROM
									biz_outfield_group_detail gd
										RIGHT JOIN
									biz_outfield_group_info g ON g.id = gd.group_info_id
								WHERE
									gd.group_info_id IN ('%s')
										AND gd.outfield_id IN ('%s')`, strings.Join(groupIds, "','"), strings.Join(canKey, "','"))
	rows, err := tx.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	existCanMap := make(map[string]bool)
	for rows.Next() {
		var key string
		var groupId string
		var groupCode string
		rows.Scan(&groupId, &groupCode, &key)

		if strings.Trim(key, " ") != "" {
			_, isExist := existCanMap[fmt.Sprintf("%s_%s", groupId, key)]
			if !isExist {
				existCanMap[fmt.Sprintf("%s_%s", groupId, key)] = true
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	for _, v := range can {
		// 如果字段编码不正确，则跳过
		if strings.Trim(v.Remark, " ") == "" {
			continue
		}

		// 如果分组ID不存在，则跳过
		if strings.Trim(v.GroupInfoId, " ") == "" {
			continue
		}

		// 判断CAN字段是否已存在
		_, isExist := existCanMap[fmt.Sprintf("%s_%s", v.GroupInfoId, v.OutfieldId)]
		if isExist {

			// 如果字段已存在，跳过该字段
			if mode == 2 {
				continue
			}

			// 如果字段已存在，更新该字段
			if mode == 1 {
				_, err = tx.Exec(`UPDATE biz_outfield_group_detail 
											 SET unit  = ?,
												 chinesename  = ?,
												 formula  = ?,
												 data_type  = ?,
												 field_name  = ?,
												 decimals  = ?,
												 isAlarm  = ?,
												 data_map  = ?,
												 is_analysable  = ?
											WHERE group_info_id = ? AND outfield_id = ?`,
					v.Unit, v.Chinesename, v.Formula, v.DataType, v.FieldName, v.Decimals, v.IsAlarm, v.DataMap, v.IsAnalysable, v.GroupInfoId, v.OutfieldId)

				if err != nil {
					return err
				}
			}

			continue
		}

		uid, _ := uuid.NewV4()
		_, err = tx.Exec(`INSERT INTO biz_outfield_group_detail 
									(id,outfield_id,unit,outfield_sn,group_info_id,chinesename,formula,data_type,field_name,decimals,isAlarm,data_map,is_analysable)
									VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
			strings.ReplaceAll(uid.String(), "-", ""), v.OutfieldId, v.Unit, v.OutfieldSn, v.GroupInfoId, v.Chinesename, v.Formula, v.DataType, v.FieldName, v.Decimals, v.IsAlarm, v.DataMap, v.IsAnalysable)

		if err != nil {
			return err
		}

	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (vt *VehicleTypeEntity) Update() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmtUpdateVt, err := db.Prepare(`UPDATE biz_vehicle_type SET type_name=?,is_intelligent=?,is_filter_missing_column=? WHERE type_id = ?`)
	if err != nil {
		return err
	}
	defer stmtUpdateVt.Close()

	_, err = stmtUpdateVt.Exec(vt.TypeName, vt.IsIntelligent, vt.IsFilterMissingColumn, vt.TypeId)
	if err != nil {
		return err
	}

	stmtUpdateG, err := db.Prepare(`UPDATE biz_outfield_group_info
											SET
												group_name = ?,
												group_sn = ?,
												remark = ?
											WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmtUpdateG.Close()

	stmtInsertG, err := db.Prepare(`INSERT biz_outfield_group_info SET id = ?,group_name= ?,group_sn= ?,remark= ?,agreement_id= ?`)
	if err != nil {
		return err
	}
	defer stmtInsertG.Close()

	for i := 0; i < len(vt.CanGroup); i++ {
		if vt.CanGroup[i].Id != "" {
			_, err = stmtUpdateG.Exec(vt.CanGroup[i].Name, vt.CanGroup[i].Sort, vt.CanGroup[i].Code, vt.CanGroup[i].Id)
			if err != nil {
				return err
			}
			continue
		}

		uid, _ := uuid.NewV4()
		_, err = stmtInsertG.Exec(strings.ReplaceAll(uid.String(), "-", ""), vt.CanGroup[i].Name, vt.CanGroup[i].Sort, vt.CanGroup[i].Code, vt.TypeId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vt *VehicleTypeEntity) Add() error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmtInsertVt, err := db.Prepare(`INSERT biz_vehicle_type SET type_id=?,type_name=?,org_id=?,is_intelligent=?,is_filter_missing_column=?,mark = ?,in_use = ?`)
	if err != nil {
		return err
	}
	defer stmtInsertVt.Close()

	uid, _ := uuid.NewV4()

	vtid := strings.ReplaceAll(uid.String(), "-", "")
	_, err = stmtInsertVt.Exec(vtid, vt.TypeName, vt.OrgId, vt.IsIntelligent, vt.IsFilterMissingColumn, 1, 1)
	if err != nil {
		return err
	}

	stmtInsertG, err := db.Prepare(`INSERT biz_outfield_group_info SET id = ?,group_name= ?,group_sn= ?,remark= ?,agreement_id= ?`)
	if err != nil {
		return err
	}
	defer stmtInsertG.Close()

	for i := 0; i < len(vt.CanGroup); i++ {
		uid, _ := uuid.NewV4()

		_, err = stmtInsertG.Exec(strings.ReplaceAll(uid.String(), "-", ""), vt.CanGroup[i].Name, vt.CanGroup[i].Sort, vt.CanGroup[i].Code, vtid)
		if err != nil {
			return err
		}
	}

	return nil
}
