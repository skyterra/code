package main

import (
	"encoding/xml"
	"errors"
	"log"
	"net/http"
	"regexp"
)

type (
	// item defines the fields associated with the item tag in the rss document.
	item struct {
		XMLName     xml.Name `xml:"item"`
		PubDate     string   `xml:"pubDate"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
		Link        string   `xml:"link"`
		GUID        string   `xml:"guid"`
		GeoRssPoint string   `xml:"georss:point"`
	}

	// image defines the fields associated with the image tag in the rss document.
	image struct {
		XMLName xml.Name `xml:"image"`
		URL     string   `xml:"url"`
		Title   string   `xml:"title"`
		Link    string   `xml:"link"`
	}

	// channel defines the fields associated with the channel tag in the rss document.
	channel struct {
		XMLName        xml.Name `xml:"channel"`
		Title          string   `xml:"title"`
		Description    string   `xml:"description"`
		Link           string   `xml:"link"`
		PubDate        string   `xml:"pubDate"`
		LastBuildDate  string   `xml:"lastBuildDate"`
		TTL            string   `xml:"ttl"`
		Language       string   `xml:"language"`
		ManagingEditor string   `xml:"managingEditor"`
		WebMaster      string   `xml:"webMaster"`
		Image          image    `xml:"image"`
		Item           []item   `xml:"item"`
	}

	// document defines the fields associated with the rss document.
	document struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)

// rssMatcher implements the Matcher interface.
type rssMatcher struct {
	Feed
}

// Match looks at the document for the specified search term.
func (m *rssMatcher) Match(searchTerm string) ([]Result, error) {
	var results []Result

	log.Printf("Search Feed Type[%s] Site[%s] For Uri[%s]\n", m.Type, m.Name, m.Uri)

	// Retrieve the data to search.
	document, err := m.retrieve()
	if err != nil {
		return nil, err
	}

	for _, channelItem := range document.Channel.Item {
		// Check the title for the search term.
		matched, err := regexp.MatchString(searchTerm, channelItem.Title)
		if err != nil {
			return nil, err
		}

		// If we found a match save the result.
		if matched {
			results = append(results, Result{
				Field:   "Title",
				Content: channelItem.Title,
			})
		}

		// Check the description for the search term.
		matched, err = regexp.MatchString(searchTerm, channelItem.Description)
		if err != nil {
			return nil, err
		}

		// If we found a match save the result.
		if matched {
			results = append(results, Result{
				Field:   "Description",
				Content: channelItem.Description,
			})
		}
	}

	return results, nil
}

// retrieve performs a HTTP Get request for the rss feed and unmarshals the results.
func (m *rssMatcher) retrieve() (*document, error) {
	if m.Uri == "" {
		return nil, errors.New("No rss feed uri provided")
	}

	// Retrieve the rss feed document from the web.
	resp, err := http.Get(m.Uri)
	if err != nil {
		return nil, err
	}

	// Close the response once we return from the function.
	defer resp.Body.Close()

	// TODO: check response status, don't process a non 200 response body.

	// Unmarshal the rss feed document into our struct type.
	var document document
	if err := xml.NewDecoder(resp.Body).Decode(&document); err != nil {
		return nil, err
	}

	return &document, nil
}
