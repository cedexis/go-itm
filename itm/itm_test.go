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
	resp, err := testClient.get("foo path")
	expectedError := "foo read error"
	if resp != nil {
		t.Error("Expected nil response")
	}
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s\nGot: %s", expectedError, err.Error())
	}
}
