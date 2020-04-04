package db

import (
	"database/sql"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

func GetStoredEvents() error {
	db, err := sql.Open("sqlite3", config.TadpolesDatabaseFile)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	return nil
}


func StoreEvents(events []api.Event) error {
	db, err := sql.Open("sqlite3", config.TadpolesDatabaseFile)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	for _, event := range events {
		log.Debugln(event)
		for _, attachment := range event.Attachments {
			log.Debugln(attachment)
		}
	}

	return nil
}
