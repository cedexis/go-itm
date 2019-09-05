package itm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type fakeRoundTripResponse struct {
	resp *http.Response
	err  error
}

type fakeRoundTripper struct {
	responses []fakeRoundTripResponse
}

func newFakeRoundTripper(r []fakeRoundTripResponse) *fakeRoundTripper {
	return &fakeRoundTripper{
		responses: r,
	}
}

func (r *fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Remove the first element responses slice
	var x fakeRoundTripResponse
	x, r.responses = r.responses[0], r.responses[1:]
	return x.resp, x.err
}

type someError struct {
	errorString string
}

func (e *someError) Error() string {
	return e.errorString
}

func newFakeHTTPClient(transport *fakeRoundTripper) *http.Client {
	return &http.Client{
		Transport: transport,
	}
}

func TestErrorIssuingPostOnCreate(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		newFakeRoundTripper(
			[]fakeRoundTripResponse{
				fakeRoundTripResponse{
					resp: nil,
					err: &someError{
						errorString: "foo",
					},
				},
			}))
	Client, _ := NewClient(HTTPClient(fakeClient))
	opts := NewDNSAppOpts()
	opts.SetName("foo")
	opts.SetDescription("foo description")
	opts.SetFallbackCname("fallback.foo.com")
	opts.SetAppData("foo app data")
	app, err := Client.DNSApps.Create(&opts, false)
	if app != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Post https://portal.cedexis.com/api/v2/config/applications/dns.json?publish=false: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingPutOnUpdateName(t *testing.T) {
	baselineConfig := map[string]interface{}{
		"appData":       "some app data",
		"description":   "some description",
		"fallbackCname": "some fallback CNAME",
		"name":          "some name",
	}
	jsonConfig, _ := json.Marshal(baselineConfig)
	getResponse := http.Response{
		StatusCode: 200,
		Body: &fakeReaderCloser{
			reader: strings.NewReader(string(jsonConfig)),
		},
	}
	fakeClient := newFakeHTTPClient(
		newFakeRoundTripper(
			[]fakeRoundTripResponse{
				fakeRoundTripResponse{
					resp: &getResponse,
					err:  nil,
				},
				fakeRoundTripResponse{
					resp: nil,
					err: &someError{
						errorString: "foo",
					},
				},
			}))
	testClient, _ := NewClient(HTTPClient(fakeClient))
	app, err := testClient.DNSApps.UpdateName(123, "updated name")
	if app != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Put https://portal.cedexis.com/api/v2/config/applications/dns.json/123?publish=true: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\n     Got: %s", expectedError, err.Error())
	}
}

func TestNewDnsAppOpts(t *testing.T) {
	var testData = []struct {
		name          string
		description   string
		fallbackCname string
		appData       string
	}{
		{
			"Foo",
			"Foo Description",
			"Foo fallback CNAME",
			"Foo app data",
		},
		{
			"Foo",
			"Foo Description",
			"Foo fallback CNAME",
			`
	Foo app data
	With spaces

	`,
		},
	}
	for _, curr := range testData {
		opts := NewDNSAppOpts()
		opts.SetName(curr.name)
		opts.SetDescription(curr.description)
		opts.SetFallbackCname(curr.fallbackCname)
		opts.SetAppData(curr.appData)
		var optsMap map[string]interface{}
		json.Unmarshal(opts.toJSON(), &optsMap)
		if "V1_JS" != optsMap["type"] {
			t.Error(unexpectedValueString("app type", "V1_JS", optsMap["type"]))
		}
		if "dns" != optsMap["protocol"] {
			t.Error(unexpectedValueString("protocol", "dns", optsMap["protocol"]))
		}
		if curr.name != optsMap["name"] {
			t.Error(unexpectedValueString("name", curr.name, optsMap["name"]))
		}
		if curr.description != optsMap["description"] {
			t.Error(unexpectedValueString("description", curr.description, optsMap["description"]))
		}
		if curr.fallbackCname != optsMap["fallbackCname"] {
			t.Error(unexpectedValueString("fallback CNAME", curr.fallbackCname, optsMap["fallbackCname"]))
		}
		trimmedAppData := strings.TrimSpace(curr.appData)
		if trimmedAppData != optsMap["appData"] {
			t.Error(unexpectedValueString("app data", trimmedAppData, optsMap["appData"]))
		}
	}
}

func TestDnsAppCreate(t *testing.T) {
	serverInfo := newTestServerInfo(nil)
	defer serverInfo.closeServer()

	handler := func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"name":          "foo",
			"description":   "foo description",
			"fallbackCname": "fallback.foo.com",
			"appData":       "foo app data",
			"type":          "V1_JS",
			"protocol":      "dns",
		}
		responseBodyObj := DNSApp{
			ID:            123,
			Version:       1,
			Name:          "foo",
			Description:   "foo description",
			FallbackCname: "fallback.foo.com",
			TTL:           20,
			AppData:       "foo app data",
			AppCname:      "foo app cname",
		}
		err := json.NewDecoder(r.Body).Decode(&parsedBody)
		if err != nil {
			t.Fatalf("JSON decoding error: %v", err)
		}
		if !reflect.DeepEqual(expectedRequestData, parsedBody) {
			t.Error(unexpectedValueString("Request body", expectedRequestData, parsedBody))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	}

	serverInfo.mux.HandleFunc("/v2/config/applications/dns.json", handler)

	opts := NewDNSAppOpts()
	opts.SetName("foo")
	opts.SetDescription("foo description")
	opts.SetAppData("foo app data")
	opts.SetFallbackCname("fallback.foo.com")

	app, err := serverInfo.client.DNSApps.Create(&opts, false)
	if err != nil {
		t.Error(err)
	}
	if 123 != app.ID {
		t.Error(unexpectedValueString("id", 123, app.ID))
	}
	if 1 != app.Version {
		t.Error(unexpectedValueString("version", 1, app.Version))
	}
	if "foo" != app.Name {
		t.Error(unexpectedValueString("name", "foo", app.Name))
	}
	if "foo description" != app.Description {
		t.Error(unexpectedValueString("description", "foo description", app.Description))
	}
	if "fallback.foo.com" != app.FallbackCname {
		t.Error(unexpectedValueString("fallback CNAME", "fallback.foo.com", app.FallbackCname))
	}
	if 20 != app.TTL {
		t.Error(unexpectedValueString("fallback TTL", 20, app.TTL))
	}
	if "foo app data" != app.AppData {
		t.Error(unexpectedValueString("app data", "foo app data", app.AppData))
	}
	if "foo app cname" != app.AppCname {
		t.Error(unexpectedValueString("app CNAME", "foo app cname", app.AppCname))
	}
}

func TestNewDNSAppOptsToJSON(t *testing.T) {
	opts := NewDNSAppOpts()
	opts.SetName("some name")
	opts.SetDescription("some description")
	opts.SetAppData("some app data")
	opts.SetFallbackCname("some fallback CNAME")
	opts.SetTTL(20)
	opts.SetFallbackCname("some fallback CNAME")
	var optsMap map[string]interface{}
	json.Unmarshal(opts.toJSON(), &optsMap)
	if "some name" != optsMap["name"] {
		t.Error(unexpectedValueString("name", "some name", optsMap["name"]))
	}
	if "some description" != optsMap["description"] {
		t.Error(unexpectedValueString("description", "some description", optsMap["description"]))
	}
	if "some app data" != optsMap["appData"] {
		t.Error(unexpectedValueString("appData", "some app data", optsMap["appData"]))
	}
	if "some fallback CNAME" != optsMap["fallbackCname"] {
		t.Error(unexpectedValueString("fallback CNAME", "some fallback CNAME", optsMap["fallbackCname"]))
	}
	if 20.0 != optsMap["ttl"] {
		t.Error(unexpectedValueString("ttl", 20.0, optsMap["ttl"]))
	}
}
