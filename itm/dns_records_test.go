package itm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestErrorIssuingPostOnCreateDNSRecord(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	createOps := NewDNSRecordOpts(400, "sub.foo.domain.name", 401, "foo_type", 3333)
	record, err := Client.DNSRecord.Create(&createOps)
	if record != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Post https://itm.cloud.com:443/api/v2/config/authdns.json/record: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingPutOnUpdateDNSRecord(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	updateOpts := NewDNSRecordOpts(400, "sub.foo.domain.name", 401, "foo_type", 3333)
	record, err := Client.DNSRecord.Update(123, &updateOpts)
	if record != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Put https://itm.cloud.com:443/api/v2/config/authdns.json/record/123: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingGetDNSRecord(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	record, err := Client.DNSRecord.Get(123)
	if record != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Get https://itm.cloud.com:443/api/v2/config/authdns.json/record/123: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestNewDNSRecordOpts(t *testing.T) {
	var testData = []struct {
		zoneId        int
		subdomainName string
		omID          int
		recordType    string
		ttl           int
	}{
		{
			400,
			"subdomain",
			401,
			"foo_type",
			3333,
		},
	}
	for _, curr := range testData {
		opts := NewDNSRecordOpts(curr.zoneId, curr.subdomainName, curr.omID, curr.recordType, curr.ttl)
		if err := testValues("zone id", curr.zoneId, opts.DNSZoneId); err != nil {
			t.Error(unexpectedValueString("zone id", curr.zoneId, opts.DNSZoneId))
		}
		if err := testValues("subdomain", curr.subdomainName, opts.SubdomainName); err != nil {
			t.Error(unexpectedValueString("subdomain", curr.subdomainName, opts.SubdomainName))
		}
		if err := testValues("openmix app id", "{\"appId\":"+strconv.Itoa(curr.omID)+"}", opts.OMAppId); err != nil {
			t.Error(unexpectedValueString("openmix app id", "{\"appId\":"+strconv.Itoa(curr.omID)+"}", opts.OMAppId))
		}
		if err := testValues("record type", curr.recordType, opts.RecordType); err != nil {
			t.Error(unexpectedValueString("record type", curr.recordType, opts.RecordType))
		}
		if err := testValues("ttl", curr.ttl, opts.TTL); err != nil {
			t.Error(unexpectedValueString("ttl", curr.ttl, opts.TTL))
		}
	}
}

func TestDNSRecordCreate(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/authdns.json/record", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"dnsZoneId":     float64(400),
			"subdomainName": "sub.foo.domain.name",
			"response":      "{\"appId\":401}",
			"recordType":    "foo_type",
			"ttl":           float64(3333),
		}
		responseBodyObj := DNSRecord{
			Id:            123,
			DNSZoneId:     400,
			SubdomainName: "sub.foo.domain.name",
			OMAppId:       "{\"appId\":401}",
			RecordType:    "foo_type",
			TTL:           3333,
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
	createOps := NewDNSRecordOpts(400, "sub.foo.domain.name", 401, "foo_type", 3333)
	record, err := client.DNSRecord.Create(&createOps)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, record.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("dns zone id", 400, record.DNSZoneId); err != nil {
		t.Error(err)
	}
	if err := testValues("sub domain", "sub.foo.domain.name", record.SubdomainName); err != nil {
		t.Error(err)
	}
	if err := testValues("openmix application id", "{\"appId\":401}", record.OMAppId); err != nil {
		t.Error(err)
	}
	if err := testValues("record type", "foo_type", record.RecordType); err != nil {
		t.Error(err)
	}
	if err := testValues("ttl", 3333, record.TTL); err != nil {
		t.Error(err)
	}
}

func TestDNSRecordUpdate(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/authdns.json/record/123", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"dnsZoneId":     float64(400),
			"subdomainName": "sub.updated_foo.domain.name",
			"response":      "{\"appId\":401}",
			"recordType":    "foo_type",
			"ttl":           float64(3333),
		}
		responseBodyObj := DNSRecord{
			Id:            123,
			DNSZoneId:     400,
			SubdomainName: "sub.updated_foo.domain.name",
			OMAppId:       "{\"appId\":401}",
			RecordType:    "foo_type",
			TTL:           3333,
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
	createOps := NewDNSRecordOpts(400, "sub.updated_foo.domain.name", 401, "foo_type", 3333)
	record, err := client.DNSRecord.Update(123, &createOps)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, record.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("dns zone id", 400, record.DNSZoneId); err != nil {
		t.Error(err)
	}
	if err := testValues("sub domain", "sub.updated_foo.domain.name", record.SubdomainName); err != nil {
		t.Error(err)
	}
	if err := testValues("openmix application id", "{\"appId\":401}", record.OMAppId); err != nil {
		t.Error(err)
	}
	if err := testValues("record type", "foo_type", record.RecordType); err != nil {
		t.Error(err)
	}
	if err := testValues("ttl", 3333, record.TTL); err != nil {
		t.Error(err)
	}
}

func TestDNSRecordGet(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/authdns.json/record/123", func(w http.ResponseWriter, r *http.Request) {
		responseBodyObj := DNSRecord{
			Id:            123,
			DNSZoneId:     400,
			SubdomainName: "sub.foo.domain.name",
			OMAppId:       "{\"appId\":401}",
			RecordType:    "foo_type",
			TTL:           3333,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})
	record, err := client.DNSRecord.Get(123)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, record.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("dns zone id", 400, record.DNSZoneId); err != nil {
		t.Error(err)
	}
	if err := testValues("sub domain", "sub.foo.domain.name", record.SubdomainName); err != nil {
		t.Error(err)
	}
	if err := testValues("openmix application id", "{\"appId\":401}", record.OMAppId); err != nil {
		t.Error(err)
	}
	if err := testValues("record type", "foo_type", record.RecordType); err != nil {
		t.Error(err)
	}
	if err := testValues("ttl", 3333, record.TTL); err != nil {
		t.Error(err)
	}
}

func TestDNSRecordDelete(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/authdns.json/record/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	err := client.DNSRecord.Delete(123)
	if err != nil {
		t.Error(err)
	}
}
