package db

import (
	"database/sql"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func RunMigrations() error {
	db, err := sql.Open("sqlite3", config.TadpolesDatabaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer utils.CloseWithLog(db)

	// ----------
	//EVENT TABLE
	sqlCreateEventTable := `
CREATE TABLE IF NOT EXISTS event (
    id INTEGER PRIMARY KEY,
    event_time TEXT,
    comment TEXT,
    child_name TEXT,
    time_zone TEXT,
    event_key TEXT,
    member TEXT
);
`
	_, err = db.Exec(sqlCreateEventTable)
	if err != nil {
		return err
	}

	// ----------
	//EVENT_ATTACHMENT TABLE
	sqlCreateAttachmentTable := `
CREATE TABLE IF NOT EXISTS event_attachment (
    id INTEGER PRIMARY KEY,
    event_id INTEGER,
    attachment_key TEXT,
    mime_type TEXT
);
`
	_, err = db.Exec(sqlCreateAttachmentTable)
	if err != nil {
		return err
	}

	return nil
}
