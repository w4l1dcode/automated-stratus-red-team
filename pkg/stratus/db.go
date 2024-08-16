package stratus

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// InitDB initializes the SQLite database and creates the table if it does not exist.
func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS used_tactics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tactic TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// MarkTacticAsUsed inserts a tactic into the database.
func MarkTacticAsUsed(db *sql.DB, tactic string) error {
	insertSQL := `INSERT INTO used_tactics (tactic) VALUES (?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer func(statement *sql.Stmt) {
		err := statement.Close()
		if err != nil {
			logrus.Fatalf("Error closing database: %v\n", err)
		}
	}(statement) // Ensure the statement is closed after use

	_, err = statement.Exec(tactic)
	if err != nil {
		return err
	}
	return nil
}

// IsTacticUsed checks if a tactic has already been used.
func IsTacticUsed(db *sql.DB, tactic string) (bool, error) {
	querySQL := `SELECT COUNT(*) FROM used_tactics WHERE tactic = ?`
	var count int
	err := db.QueryRow(querySQL, tactic).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ResetTactics deletes all records from the used_tactics table.
func ResetTactics(db *sql.DB) error {
	resetSQL := `DELETE FROM used_tactics`
	_, err := db.Exec(resetSQL)
	if err != nil {
		return err
	}
	return nil
}
