package itm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
)

const dnsAppsBasePath = "v2/config/applications/dns.json"

// DNSAppOpts specifies settings used to create a new Citrix ITM Openmix Application
type DNSAppOpts struct {
	Name          string                   `json:"name"`
	AppData       string                   `json:"appData"`
	Description   string                   `json:"description"`
	FallbackCname string                   `json:"fallbackCname"`
	Platforms     []map[string]interface{} `json:"platforms"`
	Type          string                   `json:"type"`
	Protocol      string                   `json:"protocol"`
	AvlThreshold  int                      `json:"availabilityThreshold"`
}

// NewDNSAppOpts creates and returns a new DNSAppOpts struct. Any leading or
// trailing whitespace in appData is stripped in the resulting object.
func NewDNSAppOpts(name string, appData string, description string, fallback string, platforms []map[string]interface{}, omapptype string, protocol string, threshold int) DNSAppOpts {
	result := DNSAppOpts{
		Name:          name,
		AppData:       strings.TrimSpace(appData),
		Description:   description,
		FallbackCname: fallback,
		Platforms:     platforms,
		Type:          omapptype,
		Protocol:      protocol,
		AvlThreshold:  threshold,
	}
	return result
}

// DNSApp species settings of an existing Citrix Openmix Application
type DNSApp struct {
	Id            int                      `json:"id"`
	Name          string                   `json:"name"`
	AppData       string                   `json:"appData"`
	AppCname      string                   `json:"cname"`
	Description   string                   `json:"description"`
	FallbackCname string                   `json:"fallbackCname"`
	FallbackTtl   int                      `json:"ttl"`
	Platforms     []map[string]interface{} `json:"platforms"`
	AvlThreshold  int                      `json:"availabilityThreshold"`
	Version       int                      `json:"version"`
	Enabled       bool                     `json:"enabled"`
}

type dnsAppsListTestFunc func(*DNSApp) bool

type dnsAppsService interface {
	Create(*DNSAppOpts, bool) (*DNSApp, error)
	Update(int, *DNSAppOpts, bool) (*DNSApp, error)
	Get(int) (*DNSApp, error)
	Delete(int) error
	List(opts ...dnsAppsListTestFunc) ([]DNSApp, error)
}

type dnsAppsServiceImpl struct {
	client *Client
}

// Create a Openmix Application
func (s *dnsAppsServiceImpl) Create(opts *DNSAppOpts, publish bool) (*DNSApp, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	publishVal := "false"
	if publish {
		publishVal = "true"
	}
	qs := &url.Values{
		"publish": []string{
			publishVal,
		},
	}
	resp, err := s.client.post(dnsAppsBasePath, jsonOpts, qs)
	if err != nil {
		log.Printf("Error issuing post request from DNSAppsServiceImpl.Create: %v", err)
		return nil, err
	}
	if 201 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 201,
			Got:      resp.StatusCode,
		}
	}
	var result DNSApp
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Update a Openmix Application
func (s *dnsAppsServiceImpl) Update(id int, opts *DNSAppOpts, publish bool) (*DNSApp, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	publishVal := "false"
	if publish {
		publishVal = "true"
	}
	qs := &url.Values{
		"publish": []string{
			publishVal,
		},
	}
	resp, err := s.client.put(getDNSAppPath(id), jsonOpts, qs)
	if err != nil {
		log.Printf("Error issuing put request from DNSAppsServiceImpl.Update: %v", err)
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode,
		}
	}
	var result DNSApp
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Getting details of an Openmix Application using Openmix Application ID
func (s *dnsAppsServiceImpl) Get(id int) (*DNSApp, error) {
	var result DNSApp
	resp, err := s.client.get(getDNSAppPath(id))
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

// Delete an Openmix Application using Openmix Application ID
func (s *dnsAppsServiceImpl) Delete(id int) error {
	resp, err := s.client.delete(getDNSAppPath(id))
	if 204 != resp.StatusCode {
		return &UnexpectedHTTPStatusError{
			Expected: 204,
			Got:      resp.StatusCode,
		}
	}
	return err
}

// Get list of Openmix Application
func (s *dnsAppsServiceImpl) List(tests ...dnsAppsListTestFunc) ([]DNSApp, error) {
	resp, err := s.client.get(dnsAppsBasePath)
	if err != nil {
		return nil, err
	}
	var all []DNSApp
	var result []DNSApp
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

// Get Openmix Application APIs URL
func getDNSAppPath(id int) string {
	return fmt.Sprintf("%s/%d", dnsAppsBasePath, id)
}
