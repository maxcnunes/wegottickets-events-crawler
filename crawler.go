package main

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

const (
	baseEventsURL        = "http://www.wegottickets.com/event/"
	selectorEvents       = ".content.block-group.chatterbox-margin"
	selectorPrice        = ".searchResultsPrice strong"
	selectorVenueDetails = ".venue-details h4"
	selectorNextPage     = ".pagination_link_text.nextlink"
	selectorPrevPage     = ".pagination_link_text.prevlink"
	timeParseFormat      = "Mon 2 Jan, 2006, 3:04pm"
)

var (
	regexpPrice = regexp.MustCompile(`\d+\.\d+`)
	regexpTime  = regexp.MustCompile(`^(.*\d{1,2})(st|nd|rd|th)( .*)$`)
)

// FetchDocumentFromURL loads the DOM from a URL
func FetchDocumentFromURL(url string) (*goquery.Document, error) {
	return goquery.NewDocument(url)
}

// FetchDocumentFromFile loads the DOM from a local file
func FetchDocumentFromFile(path string) (*goquery.Document, error) {
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	return goquery.NewDocumentFromReader(data)
}

// FetchEventsFromSearchResult get the events from the search result page
func FetchEventsFromSearchResult(doc *goquery.Document) []*Event {
	sel := doc.Find(selectorEvents)
	size := len(sel.Nodes)
	if size == 0 {
		return []*Event{}
	}
	events := make([]*Event, size)
	for i := range sel.Nodes {
		el := sel.Eq(i)
		title := el.Find("h2 a").Text()
		url, _ := el.Find("h2 a").Attr("href")
		id, err := parseID(url)
		if err != nil {
			color.Red("Error parsing id for event. %s", err)
		}

		price, err := parsePrice(el.Find(selectorPrice).Text())
		if err != nil {
			color.Red("Error parsing price for event %d. %s", id, err)
		}

		venue, time := "", time.Time{}
		details := el.Find(selectorVenueDetails)
		numDetails := len(details.Nodes)
		if numDetails >= 1 {
			venue = details.Eq(0).Text()
		}
		if numDetails >= 2 {
			time, err = parseTime(details.Eq(1).Text())
			if err != nil {
				color.Red("Error parsing time for event %d. %s", id, err)
			}
		}

		events[i] = &Event{
			ID:     id,
			Artist: title, // TODO: Confirm the title is really the artist info
			Title:  title,
			Link:   url,
			Price:  price,
			Venue:  venue,
			Time:   time,
		}
	}
	return events
}

// FetchURLNextPage get the URL to the next page from the search result page
func FetchURLNextPage(doc *goquery.Document) (string, bool) {
	return doc.Find(selectorNextPage).Attr("href")
}

// FetchURLPrevPage get the URL to the previous page from the search result page
func FetchURLPrevPage(doc *goquery.Document) (string, bool) {
	return doc.Find(selectorPrevPage).Attr("href")
}

func parsePrice(rawPrice string) (float64, error) {
	matches := regexpPrice.FindAllString(rawPrice, -1)
	if len(matches) == 0 {
		return 0, errors.New("Can not find price value from element in the page (" + rawPrice + ")")
	}
	p, _ := strconv.ParseFloat(matches[0], 64)
	return p, nil
}

func parseTime(rawTime string) (time.Time, error) {
	// remove day "complement" (e.g. th,rd,st)
	// otherwise can cause problems parsing the time
	rt := regexpTime.ReplaceAllString(rawTime, "${1}${3}")
	t, err := time.Parse(timeParseFormat, rt)
	if err != nil {
		return time.Time{}, errors.New("Can not find time value from element in the page (" + rawTime + ")")
	}
	return t, nil
}

func parseID(url string) (int64, error) {
	id := strings.Replace(url, baseEventsURL, "", -1)
	val, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, errors.New("Can not find id value from element in the page (" + baseEventsURL + ")")
	}
	return val, nil
}
