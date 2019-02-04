package itm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

type fakeRoundTripper struct {
	resp *http.Response
	err  error
}

func (r fakeRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r.resp, r.err
}

type someError struct {
	errorString string
}

func (e *someError) Error() string {
	return e.errorString
}

func newFakeHTTPClient(transport fakeRoundTripper) *http.Client {
	return &http.Client{
		Transport: transport,
	}
}

func TestErrorIssuingPostOnCreate(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	createOps := NewDNSAppOpts("foo", "foo description", "fallback.foo.com", "foo app data")
	app, err := Client.DNSApps.Create(&createOps, false)
	if app != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Post https://portal.cedexis.com/api/v2/config/applications/dns.json?publish=false: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingPutOnUpdate(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	testClient, _ := NewClient(HTTPClient(fakeClient))
	updateOpts := NewDNSAppOpts("foo", "foo description", "foo fallback", "foo appData")
	app, err := testClient.DNSApps.Update(123, &updateOpts, true)
	if app != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Put https://portal.cedexis.com/api/v2/config/applications/dns.json/123?publish=true: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
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
	}
	for _, curr := range testData {
		opts := NewDNSAppOpts(curr.name, curr.description, curr.fallbackCname, curr.appData)
		if err := testValues("app type", "V1_JS", opts.Type); err != nil {
			t.Error(unexpectedValueString("app type", "V1_JS", opts.Type))
		}
		if err := testValues("protocol", "dns", opts.Protocol); err != nil {
			t.Error(unexpectedValueString("protocol", "dns", opts.Protocol))
		}
		if err := testValues("name", curr.name, opts.Name); err != nil {
			t.Error(unexpectedValueString("name", curr.name, opts.Name))
		}
		if err := testValues("description", curr.description, opts.Description); err != nil {
			t.Error(unexpectedValueString("description", curr.description, opts.Description))
		}
		if err := testValues("fallback CNAME", curr.fallbackCname, opts.FallbackCname); err != nil {
			t.Error(unexpectedValueString("fallback CNAME", curr.description, opts.FallbackCname))
		}
		if err := testValues("app data", curr.appData, opts.AppData); err != nil {
			t.Error(unexpectedValueString("app data", curr.appData, opts.AppData))
		}
	}
}

func TestDnsAppCreate(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/applications/dns.json", func(w http.ResponseWriter, r *http.Request) {
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
			Id:            123,
			Version:       1,
			Name:          "foo",
			Description:   "foo description",
			FallbackCname: "fallback.foo.com",
			FallbackTtl:   20,
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
	})
	createOps := NewDNSAppOpts("foo", "foo description", "fallback.foo.com", "foo app data")
	app, err := client.DNSApps.Create(&createOps, false)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, app.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("version", 1, app.Version); err != nil {
		t.Error(err)
	}
	if err := testValues("name", "foo", app.Name); err != nil {
		t.Error(err)
	}
	if err := testValues("description", "foo description", app.Description); err != nil {
		t.Error(err)
	}
	if err := testValues("fallback CNAME", "fallback.foo.com", app.FallbackCname); err != nil {
		t.Error(err)
	}
	if err := testValues("fallback TTL", 20, app.FallbackTtl); err != nil {
		t.Error(err)
	}
	if err := testValues("app data", "foo app data", app.AppData); err != nil {
		t.Error(err)
	}
	if err := testValues("app CNAME", "foo app cname", app.AppCname); err != nil {
		t.Error(err)
	}
}

func TestDnsAppUpdate(t *testing.T) {
	t.Skip("TODO")
}

func TestDnsAppGet(t *testing.T) {
	t.Skip("TODO")
}

func TestDnsAppDelete(t *testing.T) {
	t.Skip("TODO")
}

func TestDnsAppList(t *testing.T) {
	t.Skip("TODO")
}
