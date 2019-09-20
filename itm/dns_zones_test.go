package itm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestErrorIssuingPostOnCreateDNSZone(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	createOps := NewDNSZoneOpts("foo.domain.name", "foo description")
	zone, err := Client.DNSZone.Create(&createOps)
	if zone != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Post https://itm.cloud.com:443/api/v2/config/authdns.json: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingPutOnUpdateDNSZone(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	updateOpts := NewDNSZoneOpts("foo.domain.name", "foo description")
	zone, err := Client.DNSZone.Update(123, &updateOpts)
	if zone != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Put https://itm.cloud.com:443/api/v2/config/authdns.json/123: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingGetDNSZone(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	testClient, _ := NewClient(HTTPClient(fakeClient))
	zone, err := testClient.DNSZone.Get(123)
	if zone != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Get https://itm.cloud.com:443/api/v2/config/authdns.json/123: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestNewDNSZoneOpts(t *testing.T) {
	var testData = []struct {
		domainName  string
		description string
	}{
		{
			"foo.domain.name",
			"foo description",
		},
	}
	for _, curr := range testData {
		opts := NewDNSZoneOpts(curr.domainName, curr.description)
		if err := testValues("is zone primary", true, opts.IsPrimary); err != nil {
			t.Error(unexpectedValueString("is zone primary", true, opts.IsPrimary))
		}
		if err := testValues("domain name", curr.domainName, opts.DomainName); err != nil {
			t.Error(unexpectedValueString("domain name", curr.domainName, opts.DomainName))
		}
		if err := testValues("description", curr.description, opts.Description); err != nil {
			t.Error(unexpectedValueString("description", curr.domainName, opts.Description))
		}
	}
}

func TestDNSZoneCreate(t *testing.T) {
	teardown := setup()
	defer teardown()
	var records []map[string]interface{}
	mux.HandleFunc("/v2/config/authdns.json", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"isPrimary":   true,
			"domainName":  "foo.domain.name",
			"description": "foo description",
		}

		responseBodyObj := DNSZone{
			Id:          123,
			IsPrimary:   true,
			DomainName:  "foo.domain.name",
			Description: "foo description",
			Records:     records,
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
	createOps := NewDNSZoneOpts("foo.domain.name", "foo description")
	zone, err := client.DNSZone.Create(&createOps)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, zone.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("is primary", true, zone.IsPrimary); err != nil {
		t.Error(err)
	}
	if err := testValues("domain name", "foo.domain.name", zone.DomainName); err != nil {
		t.Error(err)
	}
	if err := testValues("description", "foo description", zone.Description); err != nil {
		t.Error(err)
	}
	for index, record := range zone.Records {
		if !reflect.DeepEqual(record, records[index]) {
			t.Error(unexpectedValueString("Records in Zone", record, records[index]))
		}
	}
}

func TestDNSZoneUpdate(t *testing.T) {
	teardown := setup()
	defer teardown()
	var records []map[string]interface{}
	record := map[string]interface{}{
		"id":            float64(1234),
		"dnsZoneId":     float64(123),
		"ttl":           float64(3600),
		"subdomainName": "xyz",
		"recordType":    "A",
		"response":      "{\"addresses\":[\"2.3.2.1\"]}",
		"quickEdit":     false}
	records = append(records, record)
	mux.HandleFunc("/v2/config/authdns.json/123", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"isPrimary":   true,
			"domainName":  "foo.updated_domain.name",
			"description": "foo description",
		}

		responseBodyObj := DNSZone{
			Id:          123,
			IsPrimary:   true,
			DomainName:  "foo.updated_domain.name",
			Description: "foo description",
			Records:     records,
		}
		err := json.NewDecoder(r.Body).Decode(&parsedBody)
		if err != nil {
			t.Fatalf("JSON decoding error: %v", err)
		}
		if !reflect.DeepEqual(expectedRequestData, parsedBody) {
			t.Error(unexpectedValueString("Request body", expectedRequestData, parsedBody))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})
	createOps := NewDNSZoneOpts("foo.updated_domain.name", "foo description")
	zone, err := client.DNSZone.Update(123, &createOps)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, zone.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("is primary", true, zone.IsPrimary); err != nil {
		t.Error(err)
	}
	if err := testValues("domain name", "foo.updated_domain.name", zone.DomainName); err != nil {
		t.Error(err)
	}
	if err := testValues("description", "foo description", zone.Description); err != nil {
		t.Error(err)
	}
	for index, record := range zone.Records {
		if !reflect.DeepEqual(record, records[index]) {
			t.Error(unexpectedValueString("Records in Zone", record, records[index]))
		}
	}
}

func TestDNSZoneGet(t *testing.T) {
	teardown := setup()
	defer teardown()
	var records []map[string]interface{}
	record := map[string]interface{}{
		"id":            float64(1234),
		"dnsZoneId":     float64(123),
		"ttl":           float64(3600),
		"subdomainName": "xyz",
		"recordType":    "A",
		"response":      "{\"addresses\":[\"2.3.2.1\"]}",
		"quickEdit":     false}
	records = append(records, record)
	mux.HandleFunc("/v2/config/authdns.json/123", func(w http.ResponseWriter, r *http.Request) {
		responseBodyObj := DNSZone{
			Id:          123,
			IsPrimary:   true,
			DomainName:  "foo.updated_domain.name",
			Description: "foo description",
			Records:     records,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})
	zone, err := client.DNSZone.Get(123)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("is primary", true, zone.IsPrimary); err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, zone.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("domain name", "foo.updated_domain.name", zone.DomainName); err != nil {
		t.Error(err)
	}
	if err := testValues("description", "foo description", zone.Description); err != nil {
		t.Error(err)
	}
	for index, record := range zone.Records {
		if !reflect.DeepEqual(record, records[index]) {
			t.Error(unexpectedValueString("Records in Zone", record, records[index]))
		}
	}
}

func TestDNSZoneDelete(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/authdns.json/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	err := client.DNSZone.Delete(123)
	if err != nil {
		t.Error(err)
	}
}

func TestDNSZoneList(t *testing.T) {
	teardown := setup()
	defer teardown()
	var domain1_records []map[string]interface{}
	var domain2_records []map[string]interface{}
	var domains []DNSZone
	record1 := map[string]interface{}{
		"id":            float64(1234),
		"dnsZoneId":     float64(123),
		"ttl":           float64(3600),
		"subdomainName": "xyz",
		"recordType":    "A",
		"response":      "{\"addresses\":[\"2.3.2.1\"]}",
		"quickEdit":     false}
	record2 := map[string]interface{}{
		"id":            float64(1235),
		"dnsZoneId":     float64(123),
		"ttl":           float64(3600),
		"subdomainName": "pqr",
		"recordType":    "A",
		"response":      "{\"addresses\":[\"2.2.2.2\"]}",
		"quickEdit":     false}
	record3 := map[string]interface{}{
		"id":            float64(1236),
		"dnsZoneId":     float64(456),
		"ttl":           float64(3600),
		"subdomainName": "xyz",
		"recordType":    "A",
		"response":      "{\"addresses\":[\"2.2.2.2\"]}",
		"quickEdit":     false}

	domain1_records = append(domain1_records, record1, record2)
	domain2_records = append(domain2_records, record3)

	domain1_zone := DNSZone{
		Id:          123,
		IsPrimary:   true,
		DomainName:  "foo.domain1.name",
		Description: "foo1 description",
		Records:     domain1_records,
	}

	domain2_zone := DNSZone{
		Id:          456,
		IsPrimary:   true,
		DomainName:  "foo.domain2.name",
		Description: "foo2 description",
		Records:     domain2_records,
	}

	domains = append(domains, domain1_zone, domain2_zone)
	mux.HandleFunc("/v2/config/authdns.json/", func(w http.ResponseWriter, r *http.Request) {
		responseBodyObj := domains
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})

	zonelist, err := client.DNSZone.List()
	if err != nil {
		t.Error(err)
	}
	for index, zone := range zonelist {
		if !reflect.DeepEqual(zone, domains[index]) {
			t.Error(unexpectedValueString("Zone parameter", zone, domains[index]))
		}
	}
}
