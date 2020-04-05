package db

import (
	"database/sql"
	"github.com/leocov-dev/tadpoles-backup/config"
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"github.com/leocov-dev/tadpoles-backup/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"time"
)

func RetrieveEvents() (events []*api.Event, err error) {
	db, err := sql.Open("sqlite3", config.TadpolesDatabaseFile)
	if err != nil {
		return nil, err
	}
	defer utils.CloseWithLog(db)

	sqlQuery := `
SELECT e.id AS event_id,
       e.event_time AS event_event_time,
       e.create_time AS event_create_time,
       e.comment AS event_comment,
       e.child_name AS event_child_name,
       e.time_zone AS event_time_zone,
       e.event_key AS event_event_key,
       e.location_display AS event_location_display,
       a.attachment_key AS event_attachment_attachment_key,
       a.mime_type AS event_attachment_mime_type
FROM event AS e
JOIN attachment AS a on e.id = a.event_id
`

	rows, err := db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer utils.CloseWithLog(rows)

	eventMap := make(map[int64]*api.Event)
	for rows.Next() {
		var eventId int64
		event := &api.Event{}
		attachment := &api.Attachment{}

		err = rows.Scan(
			&eventId,
			&event.EventTime,
			&event.CreateTime,
			&event.Comment,
			&event.ChildName,
			&event.TimeZone,
			&event.EventKey,
			&event.LocationDisplay,
			&attachment.AttachmentKey,
			&attachment.MimeType,
		)
		if err != nil {
			return nil, err
		}

		existingEvent, _ := eventMap[eventId]
		if existingEvent == nil {
			eventMap[eventId] = event
			existingEvent = event
		}

		existingEvent.Attachments = append(existingEvent.Attachments, attachment)

	}

	for _, event := range eventMap {
		events = append(events, event)
	}

	return events, nil
}

func GetMaxStoredCacheTimestamp() (maxCachedTimestamp time.Time, err error) {
	db, err := sql.Open("sqlite3", config.TadpolesDatabaseFile)
	if err != nil {
		return maxCachedTimestamp, err
	}
	defer utils.CloseWithLog(db)

	row := db.QueryRow(`SELECT MAX(create_time) FROM event`)
	if row != nil {
		err := row.Scan((*utils.JsonTime)(&maxCachedTimestamp))
		if err != nil {
			return maxCachedTimestamp, err
		}
	}
	return maxCachedTimestamp, nil
}

func StoreEvents(events []*api.Event) error {
	db, err := sql.Open("sqlite3", config.TadpolesDatabaseFile)
	if err != nil {
		return err
	}
	defer utils.CloseWithLog(db)

	for _, event := range events {
		log.Debugf("Event: %s", event)
		sqlAddEvent := `
INSERT INTO event(event_time,
                  create_time,
                  comment,
                  child_name,
                  time_zone,
                  event_key,
                  location_display)
VALUES(?,?,?,?,?,?,?);
`
		result, err := db.Exec(sqlAddEvent,
			event.EventTime,
			event.CreateTime,
			event.Comment,
			event.ChildName,
			event.TimeZone,
			event.EventKey,
			event.LocationDisplay,
		)
		if err != nil {
			return err
		}

		eventId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for _, attachment := range event.Attachments {
			sqlAddAttachment := `
INSERT INTO attachment(event_id,
                       attachment_key,
                       mime_type)
VALUES(?, ?, ?);
		`
			_, err := db.Exec(sqlAddAttachment,
				eventId,
				attachment.AttachmentKey,
				attachment.MimeType,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
