package itm

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

const dnsRecordBasePath = "v2/config/authdns.json/record"

// DNSRecordOpts specifies settings used to create a new Citrix ITM DNSRecord.
type DNSRecordOpts struct {
	DNSZoneId     int    `json:"dnsZoneId"`
	SubdomainName string `json:"subdomainName"`
	OMAppId       string `json:"response"`
	RecordType    string `json:"recordType"`
	TTL           int    `json:"ttl"`
}

// NewDNSRecordOpts creates and returns a new DNSRecord struct.
func NewDNSRecordOpts(zoneId int, subdomain string, omAppId int, recordtype string, ttl int) DNSRecordOpts {
	result := DNSRecordOpts{
		DNSZoneId:     zoneId,
		SubdomainName: subdomain,
		OMAppId:       "{\"appId\":" + strconv.Itoa(omAppId) + "}",
		RecordType:    recordtype,
		TTL:           ttl,
	}
	return result
}

// DNSRecord species settings of an existing Citrix ITM DNS Record.
type DNSRecord struct {
	Id            int    `json:"id"`
	DNSZoneId     int    `json:"dnsZoneId"`
	SubdomainName string `json:"subdomainName"`
	OMAppId       string `json:"response"`
	RecordType    string `json:"recordType"`
	TTL           int    `json:"ttl"`
}

type dnsRecordListTestFunc func(*DNSRecord) bool

type dnsRecordService interface {
	Create(*DNSRecordOpts) (*DNSRecord, error)
	Update(int, *DNSRecordOpts) (*DNSRecord, error)
	Get(int) (*DNSRecord, error)
	Delete(int) error
}

type dnsRecordServiceImpl struct {
	client *Client
}

// Create a DNSRecord
func (s *dnsRecordServiceImpl) Create(opts *DNSRecordOpts) (*DNSRecord, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.post(dnsRecordBasePath, jsonOpts, nil)
	if err != nil {
		log.Printf("Error issuing post request from DNSRecordsServiceImpl.Create: %v", err)
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode,
		}
	}
	var result DNSRecord
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Update a DNSRecord
func (s *dnsRecordServiceImpl) Update(id int, opts *DNSRecordOpts) (*DNSRecord, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.put(getDNSRecordPath(id), jsonOpts, nil)
	if err != nil {
		log.Printf("Error issuing put request from DNSRecordsServiceImpl.Update: %v", err)
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode,
		}
	}
	var result DNSRecord
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Get the information about a DNS Record using DNS Record ID
func (s *dnsRecordServiceImpl) Get(id int) (*DNSRecord, error) {
	var result DNSRecord
	resp, err := s.client.get(getDNSRecordPath(id))
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode}
	}
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Delete a DNS Record using DNS Record ID
func (s *dnsRecordServiceImpl) Delete(id int) error {
	fmt.Println(getDNSRecordPath(id))
	resp, err := s.client.delete(getDNSRecordPath(id))
	if 204 != resp.StatusCode {
		return &UnexpectedHTTPStatusError{
			Expected: 204,
			Got:      resp.StatusCode,
		}
	}
	return err
}

// Get DNS Record APIs URL
func getDNSRecordPath(id int) string {
	return fmt.Sprintf("%s/%d", dnsRecordBasePath, id)
}
