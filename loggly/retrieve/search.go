package retrieve

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Search handles process for search API.
type Search struct {
	client *Client

	SearchOption
}

func newSearch(client *Client) *Search {
	return &Search{
		client: client,
		SearchOption: SearchOption{
			From:  "-24h",
			End:   "now",
			Order: SearchOrderDesc,
		},
	}
}

// SearchOrder specifies the order of results from search API.
type SearchOrder string

const (
	// SearchOrderAsc will sort results ascending.
	SearchOrderAsc = SearchOrder("asc")
	// SearchOrderDesc will sort results descending.
	SearchOrderDesc = SearchOrder("desc")
)

func (s SearchOrder) String() string { return string(s) }

// SearchOrders are all of the SearchOrder.
func SearchOrders() []SearchOrder {
	return []SearchOrder{
		SearchOrderAsc,
		SearchOrderDesc,
	}
}

// SearchOption contains query properties for search API.
type SearchOption struct {
	// Start time for the search. Defaults to “-24h”. (See valid time parameters.)
	From string `query:"from"`
	// End time for the search. Defaults to “now”. (See valid time parameters.)
	End string `query:"until"`
	// Order of results returned, either “asc” or “desc”. Defaults to “desc”.
	Order SearchOrder `query:"order"`
}

type SearchResponse struct {
	RSID struct {
		Status      string  `json:""status"`
		DateFrom    int     `json:""date_from"`
		ElapsedTime float64 `json:""elapsed_time"`
		DateTo      int     `json:""date_to"`
		ID          string  `json:""id"`
	} `json:"rsid"`
}

type SearchError struct {
	StatusCode int
	Body       []byte
}

func (e *SearchError) Error() string {
	return fmt.Sprintf("invalid response (code: %d, body: %s)", e.StatusCode, e.Body)
}
func (s *Search) From(from string) *Search {
	o := s.SearchOption
	o.From = from
	return &Search{
		client:       s.client,
		SearchOption: o,
	}
}
func (s *Search) End(end string) *Search {
	o := s.SearchOption
	o.End = end
	return &Search{
		client:       s.client,
		SearchOption: o,
	}
}
func (s *Search) Order(order SearchOrder) *Search {
	o := s.SearchOption
	o.Order = order
	return &Search{
		client:       s.client,
		SearchOption: o,
	}
}

// Do a request and get results from search API.
// /apiv2/search?q=*&from=-2h&until=now&size=10
func (s *Search) Do(size int, query string) (*SearchResponse, error) {
	res, err := s.client.call("GET", "/apiv2/search", url.Values{
		"q":     []string{query},
		"from":  []string{s.SearchOption.From},
		"until": []string{s.SearchOption.End},
		"order": []string{s.SearchOption.Order.String()},
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
		return nil, &SearchError{
			StatusCode: res.StatusCode,
			Body:       body,
		}
	}

	// Unmarshal responses
	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
