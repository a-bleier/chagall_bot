package db

import (
	"database/sql"
	"fmt"
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

func CheckUserIsRegistered(id string) bool {
	rows, _ := chagDb.Query(`Select id_telegram from Users where id_telegram in (?);`, id)
	defer rows.Close()
	return rows.Next() //returns true if there's an entry for the specific user id

}

func ListAllBirthdays(user_id string) []string {
	rows, _ := chagDb.Query(`SELECT Date, Name, Contact FROM Birthdays 
WHERE UserID IN (SELECT id_internal from Users where id_telegram = ?)`, user_id)
	defer rows.Close()
	var output []string = make([]string, 0)
	for rows.Next() {
		var (
			date    string
			name    string
			contact string
		)
		if err := rows.Scan(&date, &name, &contact); err != nil {
			panic(err)
		}
		line := fmt.Sprintf("%s %s %s", date, name, contact)
		fmt.Println(line)
		output = append(output, line)
	}
	return output
}
