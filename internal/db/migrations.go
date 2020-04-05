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
	err = createEventTable(db)
	if err != nil {
		return err
	}

	// ---------------
	//ATTACHMENT TABLE
	err = createAttachmentTable(db)
	if err != nil {
		return err
	}

	return nil
}

func createEventTable(db *sql.DB) error {
	sqlCreateEventTable := `
CREATE TABLE IF NOT EXISTS event (
    id               INTEGER PRIMARY KEY,
    event_time       TEXT,
    create_time      TEXT,
    comment          TEXT,
    child_name       TEXT,
    time_zone        TEXT,
    event_key        TEXT NOT NULL UNIQUE,
    location_display TEXT
);
`
	_, err := db.Exec(sqlCreateEventTable)

	return err
}

func createAttachmentTable(db *sql.DB) error {
	sqlCreateAttachmentTable := `
CREATE TABLE IF NOT EXISTS attachment (
    id INTEGER PRIMARY KEY,
    event_id INTEGER,
    attachment_key TEXT NOT NULL,
    mime_type TEXT
);
CREATE INDEX IF NOT EXISTS idx_attachment_event_id
ON attachment(event_id);
`
	_, err := db.Exec(sqlCreateAttachmentTable)

	return err
}
