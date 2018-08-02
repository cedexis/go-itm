package itm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

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
		opts := NewDnsAppOpts(curr.name, curr.description, curr.fallbackCname, curr.appData)
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
	testData := []struct {
		dnsAppOpts            dnsAppOpts
		expectedRequestData   map[string]interface{}
		responseStatus        int
		responseBodyObj       DnsApp
		responseContentType   string
		expectedId            int
		expectedVersion       int
		expectedName          string
		expectedDescription   string
		expectedFallbackCname string
		expectedFallbackTtl   int
		expectedAppData       string
		expectedAppCname      string
	}{
		{
			NewDnsAppOpts("foo", "foo description", "fallback.foo.com", "foo app data"),
			map[string]interface{}{
				"name":          "foo",
				"description":   "foo description",
				"fallbackCname": "fallback.foo.com",
				"appData":       "foo app data",
				"type":          "V1_JS",
				"protocol":      "dns",
			},
			http.StatusCreated,
			DnsApp{
				Id:            123,
				Version:       1,
				Name:          "foo",
				Description:   "foo description",
				FallbackCname: "fallback.foo.com",
				FallbackTtl:   20,
				AppData:       "foo app data",
				AppCname:      "foo app cname",
			},
			"application/json",
			123,
			1,
			"foo",
			"foo description",
			"fallback.foo.com",
			20,
			"foo app data",
			"foo app cname",
		},
	}
	for _, curr := range testData {
		mux.HandleFunc("/v2/config/applications/dns.json", func(w http.ResponseWriter, r *http.Request) {
			var parsedBody map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&parsedBody)
			if err != nil {
				t.Fatalf("JSON decoding error: %v", err)
			}
			if !reflect.DeepEqual(curr.expectedRequestData, parsedBody) {
				t.Error(unexpectedValueString("Request body", curr.expectedRequestData, parsedBody))
			}
			w.Header().Set("Content-Type", curr.responseContentType)
			w.WriteHeader(curr.responseStatus)
			responseBody, _ := json.Marshal(curr.responseBodyObj)
			fmt.Fprint(w, string(responseBody))
		})
		app, err := client.DnsApps.Create(&curr.dnsAppOpts, false)
		if err != nil {
			t.Error(err)
		}
		if err := testValues("id", curr.expectedId, app.Id); err != nil {
			t.Error(err)
		}
		if err := testValues("version", curr.expectedVersion, app.Version); err != nil {
			t.Error(err)
		}
		if err := testValues("name", curr.expectedName, app.Name); err != nil {
			t.Error(err)
		}
		if err := testValues("description", curr.expectedDescription, app.Description); err != nil {
			t.Error(err)
		}
		if err := testValues("fallback CNAME", curr.expectedFallbackCname, app.FallbackCname); err != nil {
			t.Error(err)
		}
		if err := testValues("fallback TTL", curr.expectedFallbackTtl, app.FallbackTtl); err != nil {
			t.Error(err)
		}
		if err := testValues("app data", curr.expectedAppData, app.AppData); err != nil {
			t.Error(err)
		}
		if err := testValues("app CNAME", curr.expectedAppCname, app.AppCname); err != nil {
			t.Error(err)
		}
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
