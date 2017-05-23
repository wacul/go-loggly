package retrieve

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/wacul/transport/expbackoff"
)

// Client for loggly retrieve API.
// See https://www.loggly.com/docs/api-retrieving-data/
type Client struct {
	client   *http.Client
	host     string
	userName string
	password string
}

// New client will be generated
func New(account, userName, password string) *Client {
	return &Client{
		client: &http.Client{
			Transport: &expbackoff.Transport{
				Min:    5 * time.Second,
				Max:    20 * time.Second,
				Factor: 1.5,
				RetryFunc: func(res *http.Response, err error) bool {
					if err != nil {
						return false
					}
					return res.StatusCode/100 == 5
				},
			},
		},
		host:     fmt.Sprintf("%s.loggly.com", account),
		userName: userName,
		password: password,
	}
}

// Search endpoint.
func (c *Client) Search() *Search {
	return newSearch(c)
}

// Events endpoint.
func (c *Client) Events() *Events {
	return newEvents(c)
}

// TODO: implement Fields()
// See https://www.loggly.com/docs/api-retrieving-data/#facet .

func (c *Client) call(method, path string, query url.Values, body io.Reader) (*http.Response, error) {
	req, err := c.createRequest(method, path, query, body)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c *Client) createRequest(method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	u := url.URL{
		Scheme: "http",
		Host:   c.host,
		Path:   path,
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.userName, c.password)
	return req, err
}
