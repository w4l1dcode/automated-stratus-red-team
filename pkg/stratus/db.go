package stratus

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB(filepath string) *sql.DB {
	log.Println("TEST INIT DB")
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS used_tactics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tactic TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func MarkTacticAsUsed(db *sql.DB, tactic string) {
	insertSQL := `INSERT INTO used_tactics (tactic) VALUES (?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(tactic)
	if err != nil {
		log.Fatal(err)
	}
}

func IsTacticUsed(db *sql.DB, tactic string) bool {
	querySQL := `SELECT COUNT(*) FROM used_tactics WHERE tactic = ?`
	var count int
	err := db.QueryRow(querySQL, tactic).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func ResetTactics(db *sql.DB) {
	resetSQL := `DELETE FROM used_tactics`
	_, err := db.Exec(resetSQL)
	if err != nil {
		log.Fatal(err)
	}
}
