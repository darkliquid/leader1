package database

import (
	"database/sql"
	"github.com/darkliquid/leader1/config"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var db *sql.DB
var logger *log.Logger
var cfg *config.DbSettings

func init() {
	logger = log.New(os.Stdout, "[database] ", log.LstdFlags)
}

func Config(dbCfg *config.DbSettings) {
	cfg = dbCfg
}

func DB() (*sql.DB, error) {
	var err error

	// No DB? Set it up!
	if db == nil {
		db, err = openDB()
	} else if err = db.Ping(); err != nil {
		db.Close()
		logger.Printf("MySQL (2) error: %s", err.Error())
		db, err = openDB()
		return db, err
	}
	return db, err
}

func openDB() (*sql.DB, error) {
	logger.Printf("MySQL: setting up new connection")

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		logger.Printf("MySQL (1) error: %s", err.Error())
		return db, err
	}

	// Set connection limits
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	if err = db.Ping(); err != nil {
		db.Close()
		logger.Printf("MySQL (2) error: %s", err.Error())
		return db, err
	}

	logger.Printf("MySQL: connected")

	return db, err
}
