package dataSource

import (
	"ConfigurationTools/model"
	"database/sql"
	"fmt"
	"strings"
)

type OrganizationEntity struct {
	Id      string
	OrgName string
	OrgType int

	SearchKey string
}

func (org *OrganizationEntity) ListSubOrg() ([]*OrganizationEntity, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sqlStr := `SELECT 
			    id, org_name, is_manufacturer
			FROM
			    sys_organization
			WHERE
			    in_use = 1 `

	condition := []string{}
	if strings.TrimSpace(org.Id) == "" {
		condition = append(condition, " (parent_id IS NULL OR parent_id = '')")
	} else {
		condition = append(condition, " parent_id = ? ")
	}

	if len(condition) > 0 {
		sqlStr += " AND " + strings.Join(condition, " AND ")
	}

	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	if strings.TrimSpace(org.Id) == "" {
		rows, err = stmt.Query()
	} else {
		rows, err = stmt.Query(org.Id)
	}

	if err != nil {
		return nil, err
	}

	orgs := []*OrganizationEntity{}
	for rows.Next() {
		org := &OrganizationEntity{}
		rows.Scan(&org.Id, &org.OrgName, &org.OrgType)

		orgs = append(orgs, org)
	}
	defer rows.Close()

	return orgs, nil
}

// ListVehicleType 车辆类型列表，包括基本信息及各个分组字段数量统计
func (org *OrganizationEntity) ListVehicleType() (vtype []*model.VehicleTypeStats, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// 1、查询车型列表，按机构名称、车型名称排序
	where := []string{}
	if strings.TrimSpace(org.SearchKey) != "" {
		where = append(where, " (vt.type_name LIKE ? OR org.org_name LIKE ?) ")
	}

	linkOper := ""
	if len(where) > 0 {
		linkOper = " WHERE "
	}
	vtsql := fmt.Sprintf(`SELECT 
					vt.type_id, vt.type_name, org.org_name, org.id
				FROM
					biz_vehicle_type vt
						INNER JOIN
					sys_organization org ON vt.org_id = org.id AND vt.in_use = 1 %s %s
				ORDER BY org.org_name, vt.type_name`, linkOper, strings.Join(where, " AND "))

	var vtRows *sql.Rows
	if strings.TrimSpace(org.SearchKey) != "" {
		vtRows, err = db.Query(vtsql, "%"+org.SearchKey+"%", "%"+org.SearchKey+"%")
	} else {
		vtRows, err = db.Query(vtsql)
	}

	if err != nil {
		return nil, err
	}
	defer vtRows.Close()

	vtype = []*model.VehicleTypeStats{}
	for vtRows.Next() {
		vt := new(model.VehicleTypeStats)
		vtRows.Scan(&vt.TypeId, &vt.TypeName, &vt.OrgName, &vt.OrgId)
		vtype = append(vtype, vt)
	}

	err = vtRows.Err()
	if err != nil {
		return nil, err
	}

	// 如果车型不存在，直接返回
	if len(vtype) == 0 {
		return vtype, nil
	}

	// 2. 查询各个车型的分组统计
	typeId := []string{}
	for _, v := range vtype {
		typeId = append(typeId, v.TypeId)
	}

	whereSql := []string{}
	if len(typeId) > 0 {
		whereSql = append(whereSql, fmt.Sprintf(" agreement_id IN ('%s')", strings.Join(typeId, "','")))
	}

	linkOper = ""
	if len(whereSql) > 0 {
		linkOper = " WHERE "
	}

	groupStatSql := fmt.Sprintf(`SELECT 
						IFNULL(g.group_name, ''),
						IFNULL(g.remark, ''),
						g.agreement_id,
						g.id,
						COUNT(gd.id)
					FROM
						biz_outfield_group_info g
							LEFT JOIN
						biz_outfield_group_detail gd ON g.id = gd.group_info_id %s %s
					GROUP BY g.group_name, g.remark, g.agreement_id, g.id
					ORDER BY g.agreement_id, g.group_sn ASC`, linkOper, strings.Join(whereSql, " AND "))

	rows, err := db.Query(groupStatSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupStat := make(map[string][]*model.GroupStats)

	var _typeId string
	var _groupId string
	var _groupName string
	var _remark int
	var _count int

	for rows.Next() {
		rows.Scan(&_groupName, &_remark, &_typeId, &_groupId, &_count)

		_, isExist := groupStat[_typeId]
		if isExist {
			groupStat[_typeId] = append(groupStat[_typeId], &model.GroupStats{
				Name:  _groupName,
				Code:  _remark,
				Id:    _groupId,
				Count: _count,
			})
		} else {
			groupStat[_typeId] = []*model.GroupStats{
				{
					Name:  _groupName,
					Code:  _remark,
					Id:    _groupId,
					Count: _count,
				},
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// 3.数据聚合
	for _, v := range vtype {
		g, isExist := groupStat[v.TypeId]
		if isExist {
			if v.Group == nil {
				v.Group = []*model.GroupStats{}
			}

			v.Group = g

			continue
		}
		v.Group = []*model.GroupStats{}
	}

	return vtype, nil
}

func (org *OrganizationEntity) ListGroupCanDetail() (cans []*model.CanDetail, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT 
										c.id,
										c.group_info_id,
										g.group_name,
										org.org_name,
										c.outfield_id,
										IFNULL(c.unit, ''),
										IFNULL(c.outfield_sn, ''),
										IFNULL(c.chinesename, ''),
										IFNULL(c.formula, ''),
										IFNULL(c.data_type, ''),
										IFNULL(c.field_name, ''),
										IFNULL(c.decimals, ''),
										IFNULL(c.isAlarm, ''),
										IFNULL(c.data_map, ''),
										c.is_analysable,
										c.is_delete
									FROM
										biz_outfield_group_detail c
											INNER JOIN
										biz_outfield_group_info g ON c.group_info_id = g.id
											INNER JOIN
										biz_vehicle_type vt ON g.agreement_id = vt.type_id
											INNER JOIN
										sys_organization org ON vt.org_id = org.id AND org.in_use = 1
									WHERE
										g.group_name LIKE ?
											AND org.org_name LIKE ?
											AND c.outfield_id LIKE ?
											AND c.chinesename LIKE ?
											AND c.field_name LIKE ?
									ORDER BY org.org_name ASC , g.group_name ASC , c.outfield_sn ASC`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	org.SearchKey = "%" + org.SearchKey + "%"
	rows, err := stmt.Query(org.SearchKey, org.SearchKey, org.SearchKey, org.SearchKey, org.SearchKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cans = []*model.CanDetail{}
	idx := 0
	for rows.Next() {
		idx = idx + 1
		can := &model.CanDetail{Index: idx}
		rows.Scan(&can.Id, &can.GroupInfoId, &can.GroupName, &can.OrgName, &can.OutfieldId, &can.Unit, &can.OutfieldSn, &can.Chinesename, &can.Formula,
			&can.DataType, &can.FieldName, &can.Decimals, &can.IsAlarm, &can.DataMap, &can.IsAnalysable, &can.IsDelete)

		cans = append(cans, can)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cans, nil
}
