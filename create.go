package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func CreatDB(install bool) (*sql.DB, error) { // функция создает файл базы данных и создает таблицу scheduler в файле

	if install {
		_, err := os.Create(DBFile) // создаём файл scheduler.db в пути dbfile
		if err != nil {
			log.Fatal(err)
		}

		createTable := `CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date INTEGER,
		title TEXT NOT NULL,
		comment TEXT,
		repeat VARCHAR(128)
		);`
		indexTable := `CREATE INDEX DATE ON scheduler (date);`

		_, err = db.Exec(createTable)
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.Exec(indexTable)
		if err != nil {
			log.Fatal(err)
		}
	}
	db, err := sql.Open("sqlite", "./scheduler.db")

	if err != nil {
		log.Fatal(err)
	}
	return db, err
}
