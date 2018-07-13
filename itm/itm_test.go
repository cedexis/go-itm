package itm

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

func setup() func() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client, _ = NewClient(BaseURL(stringToURL(server.URL)))
	return func() {
		server.Close()
	}
}

func testClientDefaults(t *testing.T, c *Client) {
	if c.BaseURL.String() != defaultBaseURL {
		t.Error(unexpectedValueString("Base URL", c.BaseURL, defaultBaseURL))
	}
	if c.UserAgentString != defaultUserAgentString {
		t.Error(unexpectedValueString("User Agent String", c.UserAgentString, defaultUserAgentString))
	}
}

func TestNewClient(t *testing.T) {
	testData := []struct {
		httpClient      *http.Client
		baseURL         *url.URL
		expectedBaseURL string
	}{
		{
			nil,
			stringToURL("http://foo.com/api"),
			"http://foo.com/api/",
		},
		{
			nil,
			stringToURL("http://foo.com/api/"),
			"http://foo.com/api/",
		},
		{
			nil,
			nil,
			"https://portal.cedexis.com/api/",
		},
	}
	for _, current := range testData {
		c, _ := NewClient(HTTPClient(current.httpClient), BaseURL(current.baseURL))
		if current.expectedBaseURL != c.BaseURL.String() {
			t.Error(unexpectedValueString("Base URL", current, c.BaseURL))
		}
	}
}
