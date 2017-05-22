package retrieve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Events handles process for search API.
type Events struct {
	client *Client

	EventsOption
}

func newEvent(client *Client) Events {
	return Events{
		client: client,
		EventsOption: EventsOption{
			Page:   0,
			Format: "json",
		},
	}
}

// EventsOption contains query properties for search API.
type EventsOption struct {
	// Page of results you’d like returned. Defaults to 0 if not defined. Pages are also zero-indexed, starting at 0 for the first page.
	Page int
	// Format of data you’d like to have your results returned. Options are “raw”, “csv”, or “json”. If “json” is used as a format, it will provide a file to download. Default is “json”.
	Format string
	// Columns to limit your result set to a set list of fields if you'd like.
	Columns []string
}

type EventsResponse struct {
	TotalEvents int `json:"total_events"`
	Page        int `json:"page"`
	Events      []struct {
		Tags       []string `json:"tags"`
		Timestamp  int64    `json:"timestamp"`
		LogMessage string   `json:"logmsg"`
		Event      struct {
			SysLog struct {
				Priority  string    `json:"priority"`
				Timestamp time.Time `json:"timestamp"`
				Host      string    `json:"host"`
				Severity  string    `json:"severity"`
				Facility  string    `json:"facility"`
			} `json:"syslog"`
			JSON map[string]interface{} `json:"json"`
		} `json:"event"`
		LogTypes []string `json:"logtypes"`
		ID       string   `json:"id"`
	} `json:"events"`
}

type EventsError struct {
	StatusCode int
	Body       []byte
}

func (e *EventsError) Error() string {
	return fmt.Sprintf("invalid response (code: %d, body: %s)", e.StatusCode, e.Body)
}

func (s *Events) Page(page int) *Events {
	o := s.EventsOption
	o.Page = page
	return &Events{
		client:       s.client,
		EventsOption: o,
	}
}

func (s *Events) Format(format string) *Events {
	o := s.EventsOption
	o.Format = format
	return &Events{
		client:       s.client,
		EventsOption: o,
	}
}

func (s *Events) Columns(columns []string) *Events {
	o := s.EventsOption
	o.Columns = columns
	return &Events{
		client:       s.client,
		EventsOption: o,
	}
}

// Do a request and get results from search API.
// /apiv2/events?rsid=728480292"
func (s *Events) Do(rsid string) (*EventsResponse, error) {
	res, err := s.client.call("GET", "/apiv2/events", url.Values{
		"rsid":    []string{rsid},
		"page":    []string{strconv.Itoa(s.EventsOption.Page)},
		"format":  []string{s.EventsOption.Format},
		"columns": s.EventsOption.Columns,
	}, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, &EventsError{
			StatusCode: res.StatusCode,
			Body:       body,
		}
	}

	// Unmarshal responses
	var result EventsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
