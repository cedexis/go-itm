package itm

import (
	"encoding/json"
	"fmt"
	"log"
)

const dnsZoneBasePath = "v2/config/authdns.json"

// DNSZoneOpts specifies settings used to create a new Citrix ITM DNS Zone
type DNSZoneOpts struct {
	IsPrimary   bool   `json:"isPrimary"`
	DomainName  string `json:"domainName"`
	Description string `json:"description"`
}

// NewDNSZoneOpts creates and returns a new DNSZone struct.
func NewDNSZoneOpts(name string, description string) DNSZoneOpts {
	result := DNSZoneOpts{
		IsPrimary:   true,
		DomainName:  name,
		Description: description,
	}
	return result
}

// DNSZoneApp species settings of an existing Citrix ITM DNS Zone
type DNSZone struct {
	Id          int                      `json:"id"`
	IsPrimary   bool                     `json:isPrimary`
	DomainName  string                   `json:"domainName"`
	Description string                   `json:"description"`
	Records     []map[string]interface{} `json:"records"`
}

type dnsZoneListTestFunc func(*DNSZone) bool

type dnsZoneService interface {
	Create(*DNSZoneOpts) (*DNSZone, error)
	Update(int, *DNSZoneOpts) (*DNSZone, error)
	Get(int) (*DNSZone, error)
	Delete(int) error
	List(opts ...dnsZoneListTestFunc) ([]DNSZone, error)
}

type dnsZoneServiceImpl struct {
	client *Client
}

// Create a DNSZone
func (s *dnsZoneServiceImpl) Create(opts *DNSZoneOpts) (*DNSZone, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.post(dnsZoneBasePath, jsonOpts, nil)
	if err != nil {
		log.Printf("Error issuing post request from DNSZonesServiceImpl.Create: %v", err)
		return nil, err
	}
	if 200 != resp.StatusCode && 201 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode,
		}
	}
	var result DNSZone
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Update a DNSZone
func (s *dnsZoneServiceImpl) Update(id int, opts *DNSZoneOpts) (*DNSZone, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.put(getDNSZonePath(id), jsonOpts, nil)
	if err != nil {
		log.Printf("Error issuing put request from DNSZonesServiceImpl.Update: %v", err)
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode,
		}
	}
	var result DNSZone
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Get the information about a DNS Zone using DNS Zone ID
func (s *dnsZoneServiceImpl) Get(id int) (*DNSZone, error) {
	var result DNSZone
	resp, err := s.client.get(getDNSZonePath(id))
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

// Delete a DNS Zone using DNS Zone ID
func (s *dnsZoneServiceImpl) Delete(id int) error {
	resp, err := s.client.delete(getDNSZonePath(id))
	if 204 != resp.StatusCode {
		return &UnexpectedHTTPStatusError{
			Expected: 204,
			Got:      resp.StatusCode,
		}
	}
	return err
}

// Gives the list of existing DNS Zones
func (s *dnsZoneServiceImpl) List(tests ...dnsZoneListTestFunc) ([]DNSZone, error) {
	resp, err := s.client.get(dnsZoneBasePath)
	if err != nil {
		return nil, err
	}
	var all []DNSZone
	var result []DNSZone
	json.Unmarshal(resp.Body, &all)
	for _, current := range all {
		stillOk := true
		for _, currentTest := range tests {
			stillOk = currentTest(&current)
			if !stillOk {
				break
			}
		}
		if stillOk {
			result = append(result, current)
		}
	}
	return result, nil
}

// Get DNS Zone APIs URL
func getDNSZonePath(id int) string {
	return fmt.Sprintf("%s/%d", dnsZoneBasePath, id)
}
