# wegottickets-events-crawler

This is a simple initial implementation of a web crawler of events from www.wegottickets.com. It will save the found events to a local file.
More details, please take a look the in the **Running section**.

## Usage

### Building

```bash
go build
```

Or for cross platform:

Using [goxc](https://github.com/laher/goxc).

```bash
goxc
```

### Running

```bash
Usage of ./wegottickets-events-crawler:
  -limit int
        The limit of pages the crawler will fetch.
  -out string
        The output path the crawler will save the events. (default "./events.json")
  -url string
        The wegottickets search result URL. (default "http://www.wegottickets.com/searchresults/region/0/latest")
```

**Some valid URLs**

 * http://www.wegottickets.com/searchresults/region/0/latest
 * http://www.wegottickets.com/searchresults/region/0/all

### Testing

```bash
./scripts/test.sh
```

## Tasks (summary)

* Fetch all pages from: http://www.wegotickets.com/searchresults/all
* Fetch the detail data for each event. Fields:
  * the artists playing
  * the city
  * the name of the venue
  * the date
  * the price

**Ideas for Improvements**

* Option to run "light crawler" - Fetching data only from the search result pages
* Keep a hash with all fetched data. Then is possible to avoid fetching a previous "visited" event.
* Respect robots.txt
* Show progress (remaining pages to fetch)
* Filter the search to only include music events
