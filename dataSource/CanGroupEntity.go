package dataSource

import (
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type GroupCanDetail struct {
	Id           string
	OutfieldId   string
	Alias        string
	FieldName    string
	Formula      string
	DataType     int
	Unit         string
	DataScope    string
	Decimals     int
	Sort         float64
	IsAlarm      bool
	IsAnalysable bool
	IsDelete     bool

	GroupInfoId string
}

func (gdc *GroupCanDetail) Add() (int64, error) {
	if strings.TrimSpace(gdc.OutfieldId) == "" {
		return 0, errors.New("key不能为空")
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO biz_outfield_group_detail (
			id, 
			outfield_id, 
			unit, 
			outfield_sn, 
			group_info_id, 
			chinesename, 
			formula, 
			data_type, 
			field_name, 
			decimals, 
			isAlarm, 
			data_map, 
			is_analysable, 
			is_delete)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	uid, _ := uuid.NewV4()
	id := strings.ReplaceAll(uid.String(), "-", "")
	_, err = stmt.Exec(id, gdc.OutfieldId, gdc.Unit, gdc.Sort, gdc.GroupInfoId, gdc.Alias, gdc.Formula, gdc.DataType,
		gdc.FieldName, gdc.Decimals, gdc.IsAlarm, gdc.DataScope, gdc.IsAnalysable, gdc.IsDelete)

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (gdc *GroupCanDetail) Update() (int64, error) {
	if strings.TrimSpace(gdc.Id) == "" {
		return 0, errors.New("主键不能为空")
	}

	if strings.TrimSpace(gdc.OutfieldId) == "" {
		return 0, errors.New("CAN key不能为空")
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE biz_outfield_group_detail
						SET
							outfield_id = ?,
							unit = ?,
							outfield_sn = ?,
							group_info_id = ?,
							chinesename = ?,
							formula = ?,
							data_type = ?,
							field_name = ?,
							decimals = ?,
							isAlarm = ?,
							data_map = ?,
							is_analysable = ?,
							is_delete = ?
						WHERE id = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(gdc.OutfieldId, gdc.Unit, gdc.Sort, gdc.GroupInfoId, gdc.Alias, gdc.Formula, gdc.DataType,
		gdc.FieldName, gdc.Decimals, gdc.IsAlarm, gdc.DataScope, gdc.IsAnalysable, gdc.IsDelete, gdc.Id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (gdc *GroupCanDetail) Get() (*GroupCanDetail, error) {
	if strings.TrimSpace(gdc.Id) == "" {
		return nil, errors.New("主键不能为空")
	}

	db, err := OpenDB()
	if err != nil {
		return nil, err
	}

	sql := `SELECT id,
				outfield_id,
				unit,
				outfield_sn,
				group_info_id,
				chinesename,
				formula,
				data_type,
				field_name,
				decimals,
				isAlarm,
				data_map,
				is_analysable,
				is_delete
			FROM biz_outfield_group_detail WHERE id = ?`

	err = db.QueryRow(sql, gdc.Id).Scan(gdc.Id, gdc.OutfieldId, gdc.Unit, gdc.Sort, gdc.GroupInfoId, gdc.Alias,
		gdc.Formula, gdc.DataType, gdc.FieldName, gdc.Decimals, gdc.IsAlarm, gdc.DataScope, gdc.IsAnalysable, gdc.IsDelete)
	if err != nil {
		return nil, err
	}

	return gdc, nil
}

func (gdc *GroupCanDetail) ToggleDelete(isDelete bool) (int64, error) {
	if strings.TrimSpace(gdc.Id) == "" {
		return 0, errors.New("主键不能为空")
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE biz_outfield_group_detail
						SET	is_delete = ?
						WHERE id = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	delFlag := 0
	if isDelete {
		delFlag = 1
	}

	res, err := stmt.Exec(delFlag, gdc.Id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}
