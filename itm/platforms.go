package itm

import (
	"encoding/json"
	"fmt"
	"log"
)

const platformBasePath = "v2/config/platforms.json"

// PlatformOpts specifies settings used to create a new Citrix ITM Platform
type PlatformOpts struct {
	PlatformAlias    string                 `json:"name"`
	PlatformDispName string                 `json:"displayName"`
	Category         map[string]interface{} `json:"category"`
	RadarOpts        map[string]interface{} `json:"radarConfig"`
	PublicProviderID int                    `json:"publicProviderArchetypeId"`
	Description      string                 `json:"intendedUse"`
}

// NewPlatformOpts creates and returns a new Platform struct.
func NewPlatformOpts(alias string, category map[string]interface{}, radar map[string]interface{}, publicProviderID int, description string) PlatformOpts {
	result := PlatformOpts{
		PlatformAlias:    alias,
		PlatformDispName: alias,
		Category:         category,
		RadarOpts:        radar,
		PublicProviderID: publicProviderID,
		Description:      description,
	}
	return result
}

// Platform species settings of an existing Citrix ITM Platform
type Platform struct {
	Id               int                    `json:"id"`
	PlatformAlias    string                 `json:"name"`
	PlatformDispName string                 `json:"displayName"`
	Category         map[string]interface{} `json:"category"`
	RadarOpts        map[string]interface{} `json:"radarConfig"`
	PublicProviderID int                    `json:"publicProviderArchetypeId"`
	Description      string                 `json:"intendedUse"`
}

type platformListTestFunc func(*Platform) bool

type platformService interface {
	Create(*PlatformOpts) (*Platform, error)
	Update(int, *PlatformOpts) (*Platform, error)
	Get(int) (*Platform, error)
	Delete(int) error
	List(opts ...platformListTestFunc) ([]Platform, error)
}

type platformServiceImpl struct {
	client *Client
}

// Create a Platform
func (s *platformServiceImpl) Create(opts *PlatformOpts) (*Platform, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.post(platformBasePath, jsonOpts, nil)
	if err != nil {
		log.Printf("Error issuing post request from PlatformsServiceImpl.Create: %v", err)
		return nil, err
	}
	if 201 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 201,
			Got:      resp.StatusCode,
		}
	}
	var result Platform
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Update a Platform
func (s *platformServiceImpl) Update(id int, opts *PlatformOpts) (*Platform, error) {
	jsonOpts, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.put(getPlatformPath(id), jsonOpts, nil)
	if err != nil {
		log.Printf("Error issuing put request from PlatformsServiceImpl.Update: %v", err)
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, &UnexpectedHTTPStatusError{
			Expected: 200,
			Got:      resp.StatusCode,
		}
	}
	var result Platform
	json.Unmarshal(resp.Body, &result)
	return &result, nil
}

// Get the information about Platfrom using Platform ID
func (s *platformServiceImpl) Get(id int) (*Platform, error) {
	var result Platform
	resp, err := s.client.get(getPlatformPath(id))
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

// Delete a Platform using Platform ID
func (s *platformServiceImpl) Delete(id int) error {
	resp, err := s.client.delete(getPlatformPath(id))
	if 204 != resp.StatusCode {
		return &UnexpectedHTTPStatusError{
			Expected: 204,
			Got:      resp.StatusCode,
		}
	}
	return err
}

// Gives the list of existing Platform
func (s *platformServiceImpl) List(tests ...platformListTestFunc) ([]Platform, error) {
	resp, err := s.client.get(platformBasePath)
	if err != nil {
		return nil, err
	}
	var all []Platform
	var result []Platform
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

// Get Platform APIs URL
func getPlatformPath(id int) string {
	return fmt.Sprintf("%s/%d", platformBasePath, id)
}
