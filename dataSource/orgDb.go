package dataSource

import (
	"ConfigurationTools/model"
	"database/sql"
	"strings"

	"github.com/satori/go.uuid"
)

func ListSubOrg(parentId string) ([]*model.Org, error) {
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
	if strings.TrimSpace(parentId) == "" {
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

	var rows *sql.Rows
	if strings.TrimSpace(parentId) == "" {
		rows, err = stmt.Query()
	} else {
		rows, err = stmt.Query(parentId)
	}

	if err != nil {
		return nil, err
	}

	orgs := []*model.Org{}
	for rows.Next() {
		org := &model.Org{}
		rows.Scan(&org.Id, &org.OrgName, &org.OrgType)

		orgs = append(orgs, org)
	}

	return orgs, nil
}

func EditVehicleType(vtid string, vtname string, isCreateGeneral bool, isCreateCan bool, isCreateCmd bool) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE biz_vehicle_type SET type_name=? WHERE type_id = ?`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(vtname, vtid)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(`INSERT biz_outfield_group_info SET id = ?,group_name= ?,group_sn= ?,remark= ?,agreement_id= ?`)
	if err != nil {
		return err
	}

	if isCreateGeneral {
		uid,_:= uuid.NewV4()
		_, err = stmt.Exec(strings.ReplaceAll(uid.String(), "-", ""), "常规信息", 1, 0, vtid)
		if err != nil {
			return err
		}
	}

	if isCreateCan {
		uid,_ := uuid.NewV4()
		_, err = stmt.Exec(strings.ReplaceAll(uid.String(), "-", ""), "CAN信息", 2, 1, vtid)
		if err != nil {
			return err
		}
	}

	if isCreateCmd {
		uid,_  := uuid.NewV4()

		_, err = stmt.Exec(strings.ReplaceAll(uid.String(), "-", ""), "指令下发", 3, 2, vtid)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddVehicleType(orgId string, vtname string, isCreateGeneral bool, isCreateCan bool, isCreateCmd bool) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT biz_vehicle_type SET type_id=?,type_name=?,org_id=?,mark = ?,in_use = ?`)
	if err != nil {
		return err
	}

	uid,_ := uuid.NewV4()

	typeId := strings.ReplaceAll(uid.String(), "-", "")
	_, err = stmt.Exec(typeId, vtname, orgId, 1, 1)
	if err != nil {
		return err
	}

	stmt, err = db.Prepare(`INSERT biz_outfield_group_info SET id = ?,group_name= ?,group_sn= ?,remark= ?,agreement_id= ?`)
	if err != nil {
		return err
	}

	if isCreateGeneral {
		uid,_ := uuid.NewV4()

		_, err = stmt.Exec(strings.ReplaceAll(uid.String(), "-", ""), "常规信息", 1, 0, typeId)
		if err != nil {
			return err
		}
	}

	if isCreateCan {
		uid,_:= uuid.NewV4()

		_, err = stmt.Exec(strings.ReplaceAll(uid.String(), "-", ""), "CAN信息", 2, 1, typeId)
		if err != nil {
			return err
		}
	}

	if isCreateCmd {
		uid,_ = uuid.NewV4()
		_, err = stmt.Exec(strings.ReplaceAll(uid.String(), "-", ""), "指令下发", 3, 2, typeId)
		if err != nil {
			return err
		}
	}

	return nil
}
