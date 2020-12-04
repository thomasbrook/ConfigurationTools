package dataSource

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// CommandConfig 指令配置模型
type CommandConfigEntity struct {
	Id             int
	CmdName        string
	Url            string
	Desc           string
	RenderTemplate string
	ParamTemplate  string
	DefaultValue   string
	GroupType      int
	CreateTime     time.Time
}

// ListCommand 获取指令配置列表
func (cc *CommandConfigEntity) ListCommand() ([]*CommandConfigEntity, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var where []string
	if strings.TrimSpace(cc.CmdName) != "" {
		where = append(where, `(commandName LIKE ? OR description LIKE ?)`)
	}

	var linkOper string
	if len(where) > 0 {
		linkOper = " AND "
	}

	_sql := fmt.Sprintf(`SELECT id,
				ifnull(commandName,''),
				ifnull(url,''),
				ifnull(description,''),
				ifnull(renderTemplate,''),
				ifnull(paramTemplate,''),
				groupType,
				defaultValue,
			 	ifnull(createTime,'')
			FROM biz_command_config 
			WHERE isDelete = 0 %s %s
			ORDER BY createTime DESC`, linkOper, strings.Join(where, " AND "))

	var rows *sql.Rows
	if strings.TrimSpace(cc.CmdName) != "" {
		rows, err = db.Query(_sql, "%"+cc.CmdName+"%", "%"+cc.CmdName+"%")
	} else {
		rows, err = db.Query(_sql)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cmds []*CommandConfigEntity
	for rows.Next() {
		cmd := new(CommandConfigEntity)
		var createDate string
		rows.Scan(&cmd.Id, &cmd.CmdName, &cmd.Url, &cmd.Desc, &cmd.RenderTemplate, &cmd.ParamTemplate, &cmd.GroupType, &cmd.DefaultValue, &createDate)
		t, err := time.ParseInLocation("2006-01-02 15:04:05", createDate, time.Local)
		if err == nil {
			cmd.CreateTime = t
		}
		cmds = append(cmds, cmd)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return cmds, nil
}

// GetCommand 获取指令配置详情
func (cc *CommandConfigEntity) GetCommand() (*CommandConfigEntity, error) {
	if cc.Id == 0 {
		return nil, errors.New("ID 不正确")
	}

	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	_sql := `SELECT id,
				commandName,
				url,
				description,
				renderTemplate,
				paramTemplate,
				groupType,
				defaultValue
			FROM biz_command_config
			WHERE id = ?`

	var cmd *CommandConfigEntity
	err = db.QueryRow(_sql, cc.Id).Scan(&cmd.Id, &cmd.CmdName, &cmd.Url, &cmd.Desc, &cmd.RenderTemplate, &cmd.ParamTemplate, &cmd.GroupType, &cmd.DefaultValue)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

// UpdateCommand 更新指令配置
func (cc *CommandConfigEntity) UpdateCommand() (int64, error) {

	if cc.Id == 0 {
		return 0, errors.New("ID 不正确")
	}

	if strings.TrimSpace(cc.CmdName) == "" {
		return 0, errors.New("指令名称不能为空")
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE biz_command_config
			SET
				commandName = ?,
				url = ?,
				description = ?,
				renderTemplate = ?,
				paramTemplate = ?,
				groupType = ?,
				defaultValue = ?,
				updateTime = ?
			WHERE id = ?`)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(cc.CmdName, cc.Url, cc.Desc, cc.RenderTemplate, cc.ParamTemplate, cc.GroupType, cc.DefaultValue, time.Now(), cc.Id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}

	return count, nil
}

// DeleteCommand 删除指令配置
func (cc *CommandConfigEntity) DeleteCommand() (int64, error) {

	if cc.Id == 0 {
		return 0, errors.New("ID 不正确")
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE biz_command_config SET isDelete = 1 WHERE id = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(cc.Id)
	if err != nil {
		return 0, nil
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (cc *CommandConfigEntity) AddCommand() (int64, error) {
	if strings.TrimSpace(cc.CmdName) == "" {
		return 0, errors.New("指令名称不能为空")
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO biz_command_config
									SET commandName = ?,
									url = ?,
									description = ?,
									renderTemplate = ?,
									paramTemplate = ?,
									groupType = ?,
									defaultValue = ?,
									createTime = ?,
									updateTime = ?`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(cc.CmdName, cc.Url, cc.Desc, cc.RenderTemplate, cc.ParamTemplate, cc.GroupType, cc.DefaultValue, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return id, nil
}

type TerminalTypeEntity struct {
	Id       string
	TypeName string
}

func (tt *TerminalTypeEntity) ListTerminalType() ([]*TerminalTypeEntity, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	_sql := `SELECT id, terminal_typename FROM biz_terminal_type ORDER BY terminal_typename ASC`
	rows, err := db.Query(_sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []*TerminalTypeEntity
	for rows.Next() {
		temp := new(TerminalTypeEntity)
		rows.Scan(&temp.Id, &temp.TypeName)

		data = append(data, temp)
	}

	err = rows.Err()
	if err != nil {
		return data, err
	}

	return data, nil
}

func (tt *TerminalTypeEntity) ListCmdConfigRelId() ([]int, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	_sql := `SELECT commandConfigId FROM biz_terminaltype_cmdcfg_rel WHERE terminalTypeId = ?`
	rows, err := db.Query(_sql, tt.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []int
	for rows.Next() {
		var temp int
		rows.Scan(&temp)
		data = append(data, temp)
	}

	err = rows.Err()
	if err != nil {
		return data, err
	}
	return data, nil
}

func (tt *TerminalTypeEntity) RelateCommandConfig(cmdCfgId []int) error {
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

	_sql := `DELETE FROM biz_terminaltype_cmdcfg_rel WHERE terminalTypeId = ?`
	_, err = tx.Exec(_sql, tt.Id)
	if err != nil {
		return err
	}

	for _, val := range cmdCfgId {
		uid, _ := uuid.NewV4()
		_, err = tx.Exec(`INSERT INTO biz_terminaltype_cmdcfg_rel(id,commandConfigId,terminalTypeId,realtionTime)VALUES(?,?,?,?)`,
			strings.ReplaceAll(uid.String(), "-", ""), val, tt.Id, time.Now())
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
