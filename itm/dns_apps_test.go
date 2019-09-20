package itm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
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

func TestErrorIssuingPostOnCreateDNSApps(t *testing.T) {
	var platformList []map[string]interface{}
	platformInstance := make(map[string]interface{})
	platformInstance["id"] = float64(0)
	platformInstance["cname"] = "foo.com"
	platformList = append(platformList, platformInstance)
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	createOps := NewDNSAppOpts("foo", "foo app data", "foo description", "fallback.foo.com", platformList, "RT_HTTP_PERFORMANCE", "dns", 80)
	app, err := Client.DNSApps.Create(&createOps, false)
	if app != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Post https://itm.cloud.com:443/api/v2/config/applications/dns.json?publish=false: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingPutOnUpdateDNSApps(t *testing.T) {
	var platformList []map[string]interface{}
	platformInstance := make(map[string]interface{})
	platformInstance["id"] = 0
	platformInstance["cname"] = "foo.com"
	platformList = append(platformList, platformInstance)
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	testClient, _ := NewClient(HTTPClient(fakeClient))
	updateOpts := NewDNSAppOpts("foo", "foo appData", "foo description", "foo fallback", platformList, "RT_HTTP_PERFORMANCE", "dns", 80)
	app, err := testClient.DNSApps.Update(123, &updateOpts, true)
	if app != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Put https://itm.cloud.com:443/api/v2/config/applications/dns.json/123?publish=true: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingGetDNSApps(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	testClient, _ := NewClient(HTTPClient(fakeClient))
	app, err := testClient.DNSApps.Get(123)
	if app != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Get https://itm.cloud.com:443/api/v2/config/applications/dns.json/123: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestNewDnsAppOpts(t *testing.T) {
	var platformList []map[string]interface{}
	platformInstance := make(map[string]interface{})
	platformInstance["id"] = float64(0)
	platformInstance["cname"] = "foo.com"
	platformList = append(platformList, platformInstance)
	var testData = []struct {
		name          string
		appData       string
		description   string
		fallbackCname string
		platform      []map[string]interface{}
		omapptype     string
		protocol      string
		threshold     int
	}{
		{
			"Foo",
			"Foo Description",
			"Foo fallback CNAME",
			"Foo app data",
			platformList,
			"RT_HTTP_PERFORMANCE",
			"dns",
			80,
		},
		{
			"Foo",
			"Foo app data With spaces",
			"Foo Description",
			"Foo fallback CNAME",
			[]map[string]interface{}{},
			"V1_JS",
			"dns",
			80,
		},
	}
	for _, curr := range testData {
		opts := NewDNSAppOpts(curr.name, curr.appData, curr.description, curr.fallbackCname, curr.platform, curr.omapptype, curr.protocol, curr.threshold)
		if err := testValues("name", curr.name, opts.Name); err != nil {
			t.Error(unexpectedValueString("name", curr.name, opts.Name))
		}
		trimmed := strings.TrimSpace(curr.appData)
		if err := testValues("app data", trimmed, opts.AppData); err != nil {
			t.Error(unexpectedValueString("app data", trimmed, opts.AppData))
		}
		if err := testValues("description", curr.description, opts.Description); err != nil {
			t.Error(unexpectedValueString("description", curr.description, opts.Description))
		}
		if err := testValues("fallback CNAME", curr.fallbackCname, opts.FallbackCname); err != nil {
			t.Error(unexpectedValueString("fallback CNAME", curr.description, opts.FallbackCname))
		}
		if err := reflect.DeepEqual(curr.platform, opts.Platforms); !err {
			t.Error(unexpectedValueString("platforms", curr.platform, opts.Platforms))
		}
		if err := testValues("app type", curr.omapptype, opts.Type); err != nil {
			t.Error(unexpectedValueString("app type", curr.omapptype, opts.Type))
		}
		if err := testValues("protocol", curr.protocol, opts.Protocol); err != nil {
			t.Error(unexpectedValueString("protocol", curr.protocol, opts.Protocol))
		}
		if err := testValues("threshold", curr.threshold, opts.AvlThreshold); err != nil {
			t.Error(unexpectedValueString("fallback CNAME", curr.threshold, opts.AvlThreshold))
		}
	}
}

func TestDnsAppCreate(t *testing.T) {
	teardown := setup()
	defer teardown()
	var platformList []map[string]interface{}
	platformInstance := make(map[string]interface{})
	platformInstance["id"] = float64(0)
	platformInstance["cname"] = "foo.com"
	platformList = append(platformList, platformInstance)
	var platInterface []interface{}
	mux.HandleFunc("/v2/config/applications/dns.json", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"name":                  "foo",
			"appData":               "foo app data",
			"description":           "foo description",
			"fallbackCname":         "fallback.foo.com",
			"platforms":             append(platInterface, platformInstance),
			"type":                  "RT_HTTP_PERFORMANCE",
			"protocol":              "dns",
			"availabilityThreshold": float64(80),
		}

		responseBodyObj := DNSApp{
			Id:            123,
			Name:          "foo",
			AppData:       "foo app data",
			AppCname:      "foo app cname",
			Description:   "foo description",
			FallbackCname: "fallback.foo.com",
			Platforms:     platformList,
			FallbackTtl:   20,
			AvlThreshold:  80,
			Version:       1,
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

	createOps := NewDNSAppOpts("foo", "foo app data", "foo description", "fallback.foo.com", platformList, "RT_HTTP_PERFORMANCE", "dns", 80)
	app, err := client.DNSApps.Create(&createOps, false)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, app.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("name", "foo", app.Name); err != nil {
		t.Error(err)
	}
	if err := testValues("app data", "foo app data", app.AppData); err != nil {
		t.Error(err)
	}
	if err := testValues("app CNAME", "foo app cname", app.AppCname); err != nil {
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
	if err := testValues("availability threshold", 80, app.AvlThreshold); err != nil {
		t.Error(err)
	}
	if err := testValues("version", 1, app.Version); err != nil {
		t.Error(err)
	}
	if err := testValues("version", 1, app.Version); err != nil {
		t.Error(err)
	}
	for index, platform := range app.Platforms {
		if !reflect.DeepEqual(platform, platformList[index]) {
			t.Error(unexpectedValueString("Zone parameter", platform, platformList[index]))
		}
	}
}

func TestDnsAppUpdate(t *testing.T) {
	teardown := setup()
	defer teardown()
	var platformList []map[string]interface{}
	platformInstance := make(map[string]interface{})
	platformInstance["id"] = float64(0)
	platformInstance["cname"] = "foo.com"
	platformList = append(platformList, platformInstance)
	var platInterface []interface{} // this is because json.NewDecoder.Decode is coverting []map[string]interface{} to []interface{} while json conversion og platformList.

	OMApp := DNSApp{
		Id:            123,
		Name:          "updated_foo",
		AppData:       "updated foo app data",
		AppCname:      "updated foo app cname",
		Description:   "updated foo description",
		FallbackCname: "fallback.foo.com",
		FallbackTtl:   20,
		Platforms:     platformList,
		AvlThreshold:  90,
		Version:       1,
	}
	mux.HandleFunc("/v2/config/applications/dns.json/123", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"name":                  "updated_foo",
			"appData":               "updated foo app data",
			"description":           "updated foo description",
			"fallbackCname":         "fallback.foo.com",
			"platforms":             append(platInterface, platformInstance),
			"type":                  "RT_HTTP_PERFORMANCE",
			"protocol":              "dns",
			"availabilityThreshold": float64(90),
		}

		responseBodyObj := OMApp
		err := json.NewDecoder(r.Body).Decode(&parsedBody)
		if err != nil {
			t.Fatalf("JSON decoding error: %v", err)
		}
		if !reflect.DeepEqual(expectedRequestData, parsedBody) {
			fmt.Println(reflect.TypeOf(parsedBody["platforms"]).String())
			fmt.Println(reflect.TypeOf(expectedRequestData["platforms"]).String())
			t.Error(unexpectedValueString("Request body", expectedRequestData, parsedBody))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})
	createOps := NewDNSAppOpts("updated_foo", "updated foo app data", "updated foo description", "fallback.foo.com", platformList, "RT_HTTP_PERFORMANCE", "dns", 90)
	app, err := client.DNSApps.Update(123, &createOps, false)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, app.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("name", "updated_foo", app.Name); err != nil {
		t.Error(err)
	}
	if err := testValues("app data", "updated foo app data", app.AppData); err != nil {
		t.Error(err)
	}
	if err := testValues("app CNAME", "updated foo app cname", app.AppCname); err != nil {
		t.Error(err)
	}
	if err := testValues("description", "updated foo description", app.Description); err != nil {
		t.Error(err)
	}
	if err := testValues("fallback CNAME", "fallback.foo.com", app.FallbackCname); err != nil {
		t.Error(err)
	}
	if err := testValues("fallback TTL", 20, app.FallbackTtl); err != nil {
		t.Error(err)
	}
	if err := testValues("availability threshold", 90, app.AvlThreshold); err != nil {
		t.Error(err)
	}
	if err := testValues("version", 1, app.Version); err != nil {
		t.Error(err)
	}
	for index, platform := range app.Platforms {
		if !reflect.DeepEqual(platform, platformList[index]) {
			t.Error(unexpectedValueString("DNS App parameter platforms", platform, platformList[index]))
		}
	}
}

func TestDnsAppGet(t *testing.T) {
	teardown := setup()
	defer teardown()
	var platformList []map[string]interface{}
	platformInstance := make(map[string]interface{})
	platformInstance["id"] = float64(0)
	platformInstance["cname"] = "foo.com"
	platformList = append(platformList, platformInstance)
	mux.HandleFunc("/v2/config/applications/dns.json/123", func(w http.ResponseWriter, r *http.Request) {

		responseBodyObj := DNSApp{
			Id:            123,
			Name:          "foo",
			AppData:       "foo app data",
			AppCname:      "foo app cname",
			Description:   "foo description",
			FallbackCname: "fallback.foo.com",
			Platforms:     platformList,
			FallbackTtl:   20,
			AvlThreshold:  80,
			Version:       1,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})

	app, err := client.DNSApps.Get(123)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, app.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("name", "foo", app.Name); err != nil {
		t.Error(err)
	}
	if err := testValues("app data", "foo app data", app.AppData); err != nil {
		t.Error(err)
	}
	if err := testValues("app CNAME", "foo app cname", app.AppCname); err != nil {
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
	if err := testValues("availability threshold", 80, app.AvlThreshold); err != nil {
		t.Error(err)
	}
	if err := testValues("version", 1, app.Version); err != nil {
		t.Error(err)
	}
	for index, platform := range app.Platforms {
		if !reflect.DeepEqual(platform, platformList[index]) {
			t.Error(unexpectedValueString("DNS App parameter platforms", platform, platformList[index]))
		}
	}
}

func TestDnsAppDelete(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/applications/dns.json/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	err := client.DNSApps.Delete(123)
	if err != nil {
		t.Error(err)
	}
}

func TestDnsAppList(t *testing.T) {
	teardown := setup()
	defer teardown()
	var OMApps []DNSApp

	var platformList1 []map[string]interface{}
	platformInstance1 := make(map[string]interface{})
	platformInstance1["id"] = float64(0)
	platformInstance1["cname"] = "foo.com"
	platformInstance2 := make(map[string]interface{})
	platformInstance2["id"] = float64(1)
	platformInstance2["cname"] = "bar.com"
	platformList1 = append(platformList1, platformInstance1, platformInstance2)

	OMApp1 := DNSApp{
		Id:            123,
		Name:          "foo",
		AppData:       "foo app data",
		AppCname:      "foo app cname",
		Description:   "foo description",
		FallbackCname: "fallback.foo.com",
		Platforms:     platformList1,
		FallbackTtl:   20,
		AvlThreshold:  80,
		Version:       1,
	}

	var platformList2 []map[string]interface{}
	platformInstance3 := make(map[string]interface{})
	platformInstance3["id"] = float64(0)
	platformInstance3["cname"] = "foo.com"
	platformList2 = append(platformList2, platformInstance3)

	OMApp2 := DNSApp{
		Id:            456,
		Name:          "bar",
		AppData:       "bar app data",
		AppCname:      "bar app cname",
		Description:   "bar description",
		FallbackCname: "fallback.bar.com",
		Platforms:     platformList2,
		FallbackTtl:   20,
		AvlThreshold:  90,
		Version:       1,
	}

	OMApps = append(OMApps, OMApp1, OMApp2)
	mux.HandleFunc("/v2/config/applications/dns.json/", func(w http.ResponseWriter, r *http.Request) {
		responseBodyObj := OMApps
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})

	omapplist, err := client.DNSApps.List()
	if err != nil {
		t.Error(err)
	}
	for index, omapp := range omapplist {
		if !reflect.DeepEqual(omapp, OMApps[index]) {
			t.Error(unexpectedValueString("DNS App parameters", omapp, OMApps[index]))
		}
	}
}
