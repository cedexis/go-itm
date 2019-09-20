package itm

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
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
	serverURL, _ := url.Parse(server.URL)
	client, _ = NewClient(BaseURL(serverURL))
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
		baseURL         string
		expectedBaseURL string
	}{
		{
			nil,
			"http://foo.com/api",
			"http://foo.com/api/",
		},
		{
			nil,
			"http://foo.com/api/",
			"http://foo.com/api/",
		},
		{
			nil,
			"",
			"https://itm.cloud.com:443/api/",
		},
	}
	for _, current := range testData {
		var c *Client
		if current.baseURL == "" {
			c, _ = NewClient(HTTPClient(current.httpClient))
		} else {
			baseURL, _ := url.Parse(current.baseURL)
			c, _ = NewClient(HTTPClient(current.httpClient), BaseURL(baseURL))
		}
		if current.expectedBaseURL != c.BaseURL.String() {
			t.Error(unexpectedValueString("Base URL", current, c.BaseURL))
		}
	}
}

type fakeReaderCloser struct {
	readCount int
	readError error
}

func (r fakeReaderCloser) Close() error {
	return nil
}

func (r fakeReaderCloser) Read(p []byte) (n int, err error) {
	return r.readCount, r.readError
}

func TestIOErrorOnReadAllDuringGet(t *testing.T) {
	response := http.Response{
		Body: fakeReaderCloser{
			readCount: 0,
			readError: &someError{
				errorString: "foo read error",
			},
		},
	}
	// A fake client whose response raises an error to be caught and handled
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: &response,
			err:  nil,
		})
	testClient, _ := NewClient(HTTPClient(fakeClient))
	resp, err := testClient.get("foo/bar")
	expectedError := "foo read error"
	if resp != nil {
		t.Error("Expected nil response")
	}
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s\nGot: %s", expectedError, err.Error())
	}
}

type echoRequestHeadersTransport struct{}

func (r echoRequestHeadersTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	jsonString, _ := json.Marshal(req.Header)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(jsonString)),
	}, nil
}

func TestUserAgentStringOnRequest(t *testing.T) {
	// Create a new client and make a GET request. Ensure that the expected
	// User-Agent header is sent.
	testConfigs := []struct {
		userAgent      string
		expectedHeader string
	}{
		{"foo", "foo"},
		{"", defaultUserAgentString},
	}
	for _, config := range testConfigs {
		fakeHTTPClient := &http.Client{
			Transport: echoRequestHeadersTransport{},
		}
		var client *Client
		if config.userAgent == "" {
			client, _ = NewClient(HTTPClient(fakeHTTPClient))
		} else {
			client, _ = NewClient(HTTPClient(fakeHTTPClient), UserAgentString(config.userAgent))
		}
		resp, err := client.get("foo/bar")
		if err != nil {
			t.Fatal(err)
		}
		var anyJSON map[string]interface{}
		json.Unmarshal(resp.Body, &anyJSON)
		if anyJSON["User-Agent"] == nil {
			t.Error("User-Agent header not sent")
		}
		userAgentHeaders := anyJSON["User-Agent"].([]interface{})
		if len(userAgentHeaders) != 1 {
			log.Printf("Expected only one User-Agent header; got %d", len(userAgentHeaders))
		}
		if userAgentHeaders[0] != config.expectedHeader {
			t.Errorf("Unexpected User-Agent string header; wanted %s; got %s", config.expectedHeader, userAgentHeaders[0])
		}
	}
}
