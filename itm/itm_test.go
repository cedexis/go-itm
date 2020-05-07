package itm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type testServerInfo struct {
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
}

func (info *testServerInfo) closeServer() {
	info.server.Close()
}

type testServerInfoOpts struct {
	clientOpts []ClientOpt
}

func newTestServerInfo(opts *testServerInfoOpts) testServerInfo {
	result := testServerInfo{
		mux: http.NewServeMux(),
	}
	result.server = httptest.NewServer(result.mux)
	serverURL, _ := url.Parse(result.server.URL)
	clientArgs := []ClientOpt{
		BaseURL(serverURL),
	}
	if opts != nil && opts.clientOpts != nil {
		for _, opt := range opts.clientOpts {
			clientArgs = append(clientArgs, opt)
		}
	}
	result.client, _ = NewClient(clientArgs...)
	return result
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
			"https://portal.cedexis.com/api/",
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
	reader    io.Reader
	readCount int
	readError error
}

func (r *fakeReaderCloser) Close() error {
	return nil
}

func (r *fakeReaderCloser) Read(p []byte) (n int, err error) {
	if r.reader != nil {
		return r.reader.Read(p)
	}
	return r.readCount, r.readError
}

func TestIOErrorOnReadAllDuringGet(t *testing.T) {
	response := http.Response{
		Body: &fakeReaderCloser{
			readCount: 0,
			readError: &someError{
				errorString: "foo read error",
			},
		},
	}
	// A fake client whose response raises an error to be caught and handled
	fakeClient := newFakeHTTPClient(
		newFakeRoundTripper([]fakeRoundTripResponse{
			fakeRoundTripResponse{
				resp: &response,
				err:  nil,
			},
		}))
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

func TestUserAgentStringOnRequest(t *testing.T) {
	// Create a new client and make a GET request.
	// Ensure that the expected User-Agent header is sent.
	testConfigs := []struct {
		userAgent               string
		expectedUserAgentHeader string
	}{
		{"foo", "foo"},
		{"", defaultUserAgentString},
	}
	for _, config := range testConfigs {
		serverInfo := newTestServerInfo(
			&testServerInfoOpts{
				clientOpts: []ClientOpt{
					UserAgentString(config.userAgent),
				},
			},
		)
		defer serverInfo.closeServer()

		handler := func(w http.ResponseWriter, req *http.Request) {
			if req.Header["User-Agent"] == nil {
				t.Error("User-Agent header not sent")
			} else {
				if 1 != len(req.Header["User-Agent"]) {
					t.Errorf("Expected only one User-Agent header; got %d\n", len(req.Header["User-Agent"]))
				}
				if req.Header["User-Agent"][0] != config.expectedUserAgentHeader {
					t.Errorf(unexpectedValueString("User-Agent header",
						config.expectedUserAgentHeader,
						req.Header["User-Agent"][0]))
				}
			}
		}
		serverInfo.mux.HandleFunc("/foo/bar", handler)
		serverInfo.client.get("foo/bar")
	}
}
