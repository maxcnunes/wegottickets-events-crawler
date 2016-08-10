package main

import "testing"

func TestFetchDocumentFromFile(t *testing.T) {
	t.Log("Fetch DOM data from a source html file")
	doc, err := FetchDocumentFromFile("./fixtures/page-01.html")
	if err != nil {
		t.Error("Error fetching document", err)
		return
	}

	if len(doc.Find(".content.block-group.chatterbox-margin").Nodes) <= 0 {
		t.Error("Expected DOM include the \"event\" elements")
	}
}

func TestFetchEventsFromSearchResult(t *testing.T) {
	t.Log("Fetch a DOM data from a source html file")
	doc, err := FetchDocumentFromFile("./fixtures/page-01.html")
	if err != nil {
		t.Error("Error fetching document", err)
		return
	}

	numResultsPerPage := 10
	events := FetchEventsFromSearchResult(doc)
	found := len(events)
	if found != numResultsPerPage {
		t.Errorf(
			"Expected to fetch all events from the page (expected: %d != actual: %d)",
			numResultsPerPage,
			found,
		)
		return
	}

	// Check all "required" fields
	for _, event := range events {
		if event.Title == "" {
			t.Error("Expected to all found events include the title value")
			return
		}

		if event.Link == "" {
			t.Error("Expected to all found events include the link to the detail page")
			return
		}

		if event.Price < 0 {
			t.Errorf("Expected to all found events include the price value (got %f)", event.Price)
			return
		}

		if event.Venue == "" {
			t.Error("Expected to all found events include the vanue value")
			return
		}

		if event.Time.IsZero() {
			t.Errorf("Expected to all found events include the time value (got %#v)", event.Time)
			return
		}
	}
}

func TestFetchURLNextPageInTheFirstPage(t *testing.T) {
	t.Log("Should fetch URL to the next page")
	doc, err := FetchDocumentFromFile("./fixtures/page-01.html")
	if err != nil {
		t.Error("Error fetching document", err)
		return
	}

	url, exists := FetchURLNextPage(doc)
	if !exists {
		t.Error("Expected to find the URL to the next search result")
		return
	}

	if url != "http://www.wegottickets.com/searchresults/page/2/latest#paginate" {
		t.Errorf("Expected to find the right URL for the next page (got %s)", url)
		return
	}
}

func TestFetchURLNextPageInTheLastPage(t *testing.T) {
	t.Log("Should not fetch URL to the next page")
	doc, err := FetchDocumentFromFile("./fixtures/page-last.html")
	if err != nil {
		t.Error("Error fetching document", err)
		return
	}

	_, exists := FetchURLNextPage(doc)
	if exists {
		t.Error("Expected to not find the URL to the next search result")
		return
	}
}

func TestFetchURLPrevPageInTheFirstPage(t *testing.T) {
	t.Log("Should not fetch URL to the prev page")
	doc, err := FetchDocumentFromFile("./fixtures/page-01.html")
	if err != nil {
		t.Error("Error fetching document", err)
		return
	}

	_, exists := FetchURLPrevPage(doc)
	if exists {
		t.Error("Expected to not find the URL to the prev search result")
		return
	}
}

func TestFetchURLPrevPageInTheLastPage(t *testing.T) {
	t.Log("Should fetch URL to the prev page")
	doc, err := FetchDocumentFromFile("./fixtures/page-last.html")
	if err != nil {
		t.Error("Error fetching document", err)
		return
	}

	url, exists := FetchURLPrevPage(doc)
	if !exists {
		t.Error("Expected to find the URL to the prev search result")
		return
	}

	if url != "http://www.wegottickets.com/searchresults/page/64/latest#paginate" {
		t.Errorf("Expected to find the right URL for the prev page (got %s)", url)
		return
	}
}
