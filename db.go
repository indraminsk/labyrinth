package main

import (
	"database/sql"
	"fmt"
	"greenjade/config"

	_ "github.com/lib/pq"
)

func openDB(cfg config.DSNType) (db *sql.DB) {
	var (
		err error

		dsn string
	)

	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.DBName)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("error [connect to db]:", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
