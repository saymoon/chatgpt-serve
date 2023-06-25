package main

import (
	"database/sql"
    "log"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks(
		id TEXT PRIMARY KEY,
		instance_id TEXT,
		conversation_id TEXT,
		model TEXT,
		prompt TEXT,
		response TEXT,
		status TEXT,
		error_message TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens(
		token TEXT PRIMARY KEY,
		is_admin BOOLEAN
	);`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}


