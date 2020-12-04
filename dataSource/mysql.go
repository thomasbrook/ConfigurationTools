package dataSource

import (
	"ConfigurationTools/configurationManager"
	"database/sql"
)

func OpenDB() (db *sql.DB, err error) {
	connStr, err := configurationManager.GetDatabaseConnStr()
	if err != nil {
		return nil, err
	}
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Ping() bool {
	db, err := OpenDB()
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return false
	}

	return true
}
