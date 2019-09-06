package itm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
)

const dnsAppsBasePath = "v2/config/applications/dns.json"

// DNSAppOpts accumulates options used to create and update a Citrix ITM DNS app
type DNSAppOpts struct {
	data map[string]interface{}
}

// NewDNSAppOpts creates a DNSAppOpts struct populated with sensible defaults
func NewDNSAppOpts() DNSAppOpts {
	result := DNSAppOpts{
		data: map[string]interface{}{},
	}
	result.data["protocol"] = "dns"
	result.data["type"] = "V1_JS"
	return result
}

func (opts DNSAppOpts) toJSON() []byte {
	result, _ := json.Marshal(opts.data)
	return result
}

// SetAppData sets JavaScript code of the app
func (opts DNSAppOpts) SetAppData(value string) {
	opts.data["appData"] = strings.TrimSpace(value)
}

// SetDescription sets the description of the app
func (opts DNSAppOpts) SetDescription(value string) {
	opts.data["description"] = value
}

// SetFallbackCname sets the fallback CNAME of the app
func (opts DNSAppOpts) SetFallbackCname(value string) {
	opts.data["fallbackCname"] = value
}

// SetTTL sets the default TTL of the app
func (opts DNSAppOpts) SetTTL(value int) {
	opts.data["ttl"] = value
}

// SetName sets the name of the app
func (opts DNSAppOpts) SetName(value string) {
	opts.data["name"] = value
}

// DNSApp species settings of an existing Citrix DNS app
type DNSApp struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Enabled       bool   `json:"enabled"`
	FallbackCname string `json:"fallbackCname"`
	TTL           int    `json:"ttl"`
	AppData       string `json:"appData"`
	AppCname      string `json:"cname"`
	Version       int    `json:"version"`
}

type dnsAppsListTestFunc func(*DNSApp) bool

type dnsAppsService interface {
	Create(*DNSAppOpts, bool) (*DNSApp, error)
	Get(int) (*DNSApp, error)
	Delete(int) error
	List(opts ...dnsAppsListTestFunc) ([]DNSApp, error)
	UpdateName(int, string) (*DNSApp, error)
	UpdateDescription(int, string) (*DNSApp, error)
	UpdateAppData(int, string) (*DNSApp, error)
	UpdateFallbackCname(int, string) (*DNSApp, error)
	UpdateTTL(int, int) (*DNSApp, error)
}

type dnsAppsServiceImpl struct {
	client *Client
}

// Create a DNS app
func (s *dnsAppsServiceImpl) Create(opts *DNSAppOpts, publish bool) (*DNSApp, error) {
	publishVal := "false"
	if publish {
		publishVal = "true"
	}
	qs := &url.Values{
		"publish": []string{
			publishVal,
		},
	}
	resp, err := s.client.post(dnsAppsBasePath, opts.toJSON(), qs)
	if err != nil {
		log.Printf("[WARN] Error issuing post request from DNSAppsServiceImpl.Create: %v", err)
		return nil, err
	}
	if 201 != resp.StatusCode {
		return nil, fmt.Errorf("Unexpected HTTP response. HTTP status code: %d. Body: %s", resp.StatusCode, resp.Body)
	}
	var result DNSApp
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Update a DNS app
func (s *dnsAppsServiceImpl) update(id int, opts *DNSAppOpts, publish bool) (*DNSApp, error) {
	publishVal := "false"
	if publish {
		publishVal = "true"
	}
	qs := &url.Values{
		"publish": []string{
			publishVal,
		},
	}
	resp, err := s.client.put(getDNSAppPath(id), opts.toJSON(), qs)
	if err != nil {
		log.Printf("[WARN] Error issuing put request from DNSAppsServiceImpl.Update: %v", err)
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("Unexpected HTTP response. HTTP status code: %d. Body: %s", resp.StatusCode, resp.Body)
	}
	var result DNSApp
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

func (s *dnsAppsServiceImpl) Get(id int) (*DNSApp, error) {
	var result DNSApp
	resp, err := s.client.get(getDNSAppPath(id))
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("Unexpected HTTP response. HTTP status code: %d. Body: %s", resp.StatusCode, resp.Body)
	}
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

func (s *dnsAppsServiceImpl) Delete(id int) error {
	resp, err := s.client.delete(getDNSAppPath(id))
	if 204 != resp.StatusCode {
		return fmt.Errorf("Unexpected HTTP response. HTTP status code: %d. Error: %s", resp.StatusCode, err)
	}
	return nil
}

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

func (s *dnsAppsServiceImpl) UpdateName(id int, name string) (*DNSApp, error) {
	opts, err := s.makeUpdateDefaults(id)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] UpdateName baseline config: %#v", opts)
	opts.SetName(name)
	return s.update(id, opts, true)
}

func (s *dnsAppsServiceImpl) UpdateDescription(id int, description string) (*DNSApp, error) {
	opts, err := s.makeUpdateDefaults(id)
	if err != nil {
		return nil, err
	}
	opts.SetDescription(description)
	return s.update(id, opts, true)
}

func (s *dnsAppsServiceImpl) UpdateAppData(id int, appData string) (*DNSApp, error) {
	opts, err := s.makeUpdateDefaults(id)
	if err != nil {
		return nil, err
	}
	opts.SetAppData(appData)
	return s.update(id, opts, true)
}

func (s *dnsAppsServiceImpl) UpdateFallbackCname(id int, cname string) (*DNSApp, error) {
	opts, err := s.makeUpdateDefaults(id)
	if err != nil {
		return nil, err
	}
	opts.SetFallbackCname(cname)
	return s.update(id, opts, true)
}

func (s *dnsAppsServiceImpl) UpdateTTL(id int, ttl int) (*DNSApp, error) {
	opts, err := s.makeUpdateDefaults(id)
	if err != nil {
		return nil, err
	}
	opts.SetTTL(ttl)
	return s.update(id, opts, true)
}

func getDNSAppPath(id int) string {
	return fmt.Sprintf("%s/%d", dnsAppsBasePath, id)
}

func (s *dnsAppsServiceImpl) makeUpdateDefaults(id int) (*DNSAppOpts, error) {
	app, err := s.Get(id)
	if err != nil {
		return nil, fmt.Errorf("Error querying for app (id %d): %v", id, err)
	}
	result := NewDNSAppOpts()
	result.SetName(app.Name)
	result.SetDescription(app.Description)
	result.SetFallbackCname(app.FallbackCname)
	result.SetTTL(app.TTL)
	result.SetAppData(app.AppData)
	return &result, nil
}
