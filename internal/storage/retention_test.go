package storage

import (
	"testing"
	"time"

	"github.com/axllent/mailpit/config"
)

// TestDefaultMaxAgePrunes checks the fork default of 30 days removes older messages
func TestDefaultMaxAgePrunes(t *testing.T) {
	setup("")
	defer Close()

	config.UseMessageDates = true
	config.MaxAgeInHours = 30 * 24 // the result of the 30d default
	defer func() {
		config.UseMessageDates = false
		config.MaxAgeInHours = 0
	}()

	store := func(subject string, age time.Duration) {
		b := []byte("Date: " + time.Now().Add(-age).Format(time.RFC1123Z) + "\r\nFrom: sender@example.com\r\nTo: recipient@example.com\r\nSubject: " + subject + "\r\n\r\nHello\r\n")
		if _, err := Store("", &b, nil); err != nil {
			t.Fatal(err)
		}
	}

	store("31 days old", 31*24*time.Hour)
	store("29 days old", 29*24*time.Hour)

	pruneMessages()

	list, err := List("", 0, 0, 10)
	if err != nil {
		t.Fatal(err)
	}

	assertEqual(t, len(list), 1, "message older than 30 days not pruned")
	assertEqual(t, list[0].Subject, "29 days old", "wrong message pruned")

	// the raw message (incl. attachments) and tags are removed too
	var data int
	if err := db.QueryRow(`SELECT COUNT(*) FROM ` + tenant("mailbox_data")).Scan(&data); err != nil {
		t.Fatal(err)
	}
	assertEqual(t, data, 1, "raw message data not pruned")
}
