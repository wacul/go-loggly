package retrieve

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
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
		client:   &http.Client{},
		host:     fmt.Sprintf("%s.loggly.com", account),
		userName: userName,
		password: password,
	}
}

// Search endpoint.
func (c *Client) Search() *Search {
	return &Search{client: c}
}

// Events endpoint.
func (c *Client) Events() *Events {
	return &Events{client: c}
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
	if query == nil {
		u.RawQuery = query.Encode()
	}
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.userName, c.password)
	return req, err
}
