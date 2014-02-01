package database

import (
	"database/sql"
	"fmt"
	"github.com/darkliquid/leader1/config"
	"github.com/fluffle/golog/logging"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DB() (*sql.DB, error) {
	var err error

	// No DB? Set it up!
	if db == nil {
		db, err = sql.Open("mysql", config.Config.Db.DSN)
		if err != nil {
			logging.Error(fmt.Sprintf("MySQL (1) error: %s", err.Error()))
			return db, err
		}
		if err = db.Ping(); err != nil {
			db.Close()
			logging.Error(fmt.Sprintf("MySQL (2) error: %s", err.Error()))
			return db, err
		}
		logging.Info("MySQL: connected")
	} else if err = db.Ping(); err != nil {
		db.Close()
		logging.Error(fmt.Sprintf("MySQL (2) error: %s", err.Error()))
		return db, err
	}
	return db, err
}
