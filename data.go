package main

import (
	"bufio"
	"encoding/json"
	"os"
	"time"
)

// Event contains the information about the event
type Event struct {
	Artist string
	City   string
	Venue  string
	Time   time.Time
	Price  float64
	Title  string
	Link   string
	ID     int64
}

// Data handles the events' persistence to a local file
type Data struct {
	path     string
	file     *os.File
	writer   *bufio.Writer
	enconder *json.Encoder
}

// NewData creates a new Data
func NewData(path string) *Data {
	return &Data{path: path}
}

// Open prepares the file to persist later the events
// It will automatically create a new file in case it not exists yet
func (d *Data) Open() error {
	// detect if file exists
	var _, err = os.Stat(d.path)

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		// create file if not exists
		var file, err = os.Create(d.path)
		if err != nil {
			return err
		}
		d.file = file
	}

	if d.file == nil {
		file, err := os.OpenFile(d.path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			return err
		}
		d.file = file
	}

	// create a buffer for writing
	d.writer = bufio.NewWriter(d.file)

	// create JSON econder
	d.enconder = json.NewEncoder(d.writer)

	return nil
}

// Close the file buffer
func (d *Data) Close() {
	d.writer.Flush()
	d.file.Close()
}

// SaveEvent append a event as JSON to the output file
func (d *Data) SaveEvent(event *Event) error {
	return d.enconder.Encode(event)
}
