package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

const (
	defaultURL      = "http://www.wegottickets.com/searchresults/region/0/latest"
	defaultDataPath = "./events.json"
)

var yellow = color.New(color.FgYellow).SprintFunc()

func crawlEvents(eventsFirstPageURL string, dataPath string, limit int) error {
	data := NewData(dataPath)
	if err := data.Open(); err != nil {
		return err
	}

	defer data.Close()

	page := 1
	eventsURL := eventsFirstPageURL
	for eventsURL != "" && (limit == 0 || page <= limit) {
		color.Cyan("Fetching page %d (%s)", page, eventsURL)
		doc, err := FetchDocumentFromURL(eventsURL)
		if err != nil {
			return err
		}

		events := FetchEventsFromSearchResult(doc)
		color.Yellow("Found %d events:", len(events))

		// crawl all event details concurrently
		ch := make(chan *Event, len(events))
		for _, event := range events {
			go crawlEventDetail(event, ch)
		}

		done := 0
		for event := range ch {
			done++
			if done == len(events) {
				close(ch)
			}

			if err := data.SaveEvent(event); err != nil {
				return err
			}
		}

		eventsURL, _ = FetchURLNextPage(doc)
		page++
	}
	return nil
}

func crawlEventDetail(event *Event, ch chan *Event) {
	fmt.Printf("  %s Event %d: %s\n", yellow("*"), event.ID, event.Title)
	// TODO: Fetch event detail data
	ch <- event
}

func main() {
	eventsURL := flag.String("url", defaultURL, "The wegottickets search result URL.")
	dataPath := flag.String("out", defaultDataPath, "The output path the crawler will save the events.")
	limit := flag.Int("limit", 0, "The limit of pages the crawler will fetch.")
	flag.Parse()

	if err := crawlEvents(*eventsURL, *dataPath, *limit); err != nil {
		color.Red("Error on crawling:\n%s", err)
		os.Exit(1)
	}
}
