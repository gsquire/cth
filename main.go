package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"

	"cloud.google.com/go/pubsub"
	"github.com/gsquire/cth/decode"
	_ "github.com/mattn/go-sqlite3"
)

func dbInsert(db *sql.DB, scan *decode.Scan, scanResponse string) error {
	const (
		insertStmt = "INSERT INTO scans (ip, port, service, last_scan, response) VALUES (?, ?, ?, ?, ?)"
	)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(scan.Ip, scan.Port, scan.Service, scan.Timestamp, scanResponse)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func updateRecords(messages <-chan *pubsub.Message) {
	db, err := sql.Open("sqlite3", "./scan-results.sqlite")
	if err != nil {
		log.Fatalf("error opening database: %s\n", err)
	}
	defer db.Close()

	// We have a few message handling scenarios here:
	// - Failure to parse: Acknowledge and continue.
	// - Failure to insert: Negative Acknowledge and try again.
	// - Success: Acknowledge and continue.
	for msg := range messages {
		result, err := decode.DecodeMessage(msg.Data)
		if err != nil {
			log.Printf("dropping message due to invalid format: %s\n", err)
			msg.Ack()
			continue
		}

		err = dbInsert(db, result.ParsedScan, result.Response)
		if err != nil {
			log.Printf("retrying insert of failed message due to error: %s\n", err)
			msg.Nack()
			continue
		}

		msg.Ack()
	}
}

func main() {
	var processingRate = flag.Int("r", 1000, "how many messages we want in flight at any time")
	flag.Parse()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "test-project")
	if err != nil {
		log.Fatalf("could not create a new client: %s\n", err)
	}
	defer client.Close()

	// Supply a buffer so we can apply some back pressure when processing.
	messages := make(chan *pubsub.Message, *processingRate)
	go updateRecords(messages)

	cancelCtx, cancel := context.WithCancel(ctx)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		log.Println("received interrupt, cleaning up and exiting")
		cancel()
	}()

	sub := client.Subscription("scan-sub")
	err = sub.Receive(cancelCtx, func(_ context.Context, msg *pubsub.Message) {
		messages <- msg
	})
	if err != nil {
		log.Fatalf("subscription `Receive` encountered an error: %s\n", err)
	}

	// Now that we are done sending values, close the channel so we exit our receiver goroutine.
	close(messages)
}
