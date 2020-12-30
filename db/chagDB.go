package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var chagDb *sql.DB

func InitChagDB(dbName string) {
	var err error
	chagDb, err = sql.Open("sqlite3", "simple.sqlite")
	if err != nil {
		log.Panic(err)
	}

}

//TODO Use string for the id cause int doesn't make sense
func CheckUserIsRegistered(id string) bool {
	rows, _ := chagDb.Query(`Select id_telegram from Users where id_telegram in (?);`, id)
	defer rows.Close()
	return rows.Next() //returns true if there's an entry for the specific user id

}
