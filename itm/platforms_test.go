package itm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestErrorIssuingPostOnCreatePlatform(t *testing.T) {
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	createOps := NewPlatformOpts("foo", category, radar, 12345, "foo description")
	platform, err := Client.Platform.Create(&createOps)
	if platform != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Post https://itm.cloud.com:443/api/v2/config/platforms.json: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingPutOnUpdatePlatform(t *testing.T) {
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	Client, _ := NewClient(HTTPClient(fakeClient))
	updateOpts := NewPlatformOpts("foo", category, radar, 12345, "foo description")
	platform, err := Client.Platform.Update(123, &updateOpts)
	if platform != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Put https://itm.cloud.com:443/api/v2/config/platforms.json/123: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestErrorIssuingGetPlatform(t *testing.T) {
	fakeClient := newFakeHTTPClient(
		fakeRoundTripper{
			resp: nil,
			err: &someError{
				errorString: "foo",
			},
		})
	testClient, _ := NewClient(HTTPClient(fakeClient))
	platform, err := testClient.Platform.Get(123)
	if platform != nil {
		t.Error("Expected nil result")
	}
	expectedError := "Get https://itm.cloud.com:443/api/v2/config/platforms.json/123: foo"
	if expectedError != err.Error() {
		t.Errorf("Unexpected error.\nExpected: %s.\nGot: %s", expectedError, err.Error())
	}
}

func TestNewPlatformOpts(t *testing.T) {
	var testData = []struct {
		alias       string
		category    map[string]interface{}
		radarOpts   map[string]interface{}
		pubprovID   int
		description string
	}{
		{
			"foo",
			map[string]interface{}{"id": 1},
			map[string]interface{}{"usePublicData": true},
			12345,
			"foo description",
		},
	}
	for _, curr := range testData {
		opts := NewPlatformOpts(curr.alias, curr.category, curr.radarOpts, curr.pubprovID, curr.description)
		if err := testValues("platform alias", curr.alias, opts.PlatformAlias); err != nil {
			t.Error(unexpectedValueString("platform alias", curr.alias, opts.PlatformAlias))
		}
		if err := testValues("platform display name", curr.alias, opts.PlatformDispName); err != nil {
			t.Error(unexpectedValueString("platform disp name", curr.alias, opts.PlatformDispName))
		}
		if op := reflect.DeepEqual(curr.category, opts.Category); !op {
			t.Error(unexpectedValueString("platform category", curr.category, opts.Category))
		}
		if op := reflect.DeepEqual(curr.radarOpts, opts.RadarOpts); !op {
			t.Error(unexpectedValueString("platform radar options", curr.radarOpts, opts.RadarOpts))
		}
		if err := testValues("platform public provider id", curr.pubprovID, opts.PublicProviderID); err != nil {
			t.Error(unexpectedValueString("platform public provider id", curr.pubprovID, opts.PublicProviderID))
		}
		if err := testValues("description", curr.description, opts.Description); err != nil {
			t.Error(unexpectedValueString("description", curr.description, opts.Description))
		}
	}
}

func TestPlatformCreate(t *testing.T) {
	teardown := setup()
	defer teardown()
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	mux.HandleFunc("/v2/config/platforms.json", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"name":                      "foo",
			"displayName":               "foo",
			"category":                  category,
			"radarConfig":               radar,
			"publicProviderArchetypeId": float64(12345),
			"intendedUse":               "foo description",
		}
		responseBodyObj := Platform{
			Id:               123,
			PlatformAlias:    "foo",
			PlatformDispName: "foo",
			Category:         category,
			RadarOpts:        radar,
			PublicProviderID: 12345,
			Description:      "foo description",
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
	createOps := NewPlatformOpts("foo", category, radar, 12345, "foo description")
	platform, err := client.Platform.Create(&createOps)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, platform.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("platform alias", "foo", platform.PlatformAlias); err != nil {
		t.Error(err)
	}
	if err := testValues("platform display name", "foo", platform.PlatformDispName); err != nil {
		t.Error(err)
	}
	if op := reflect.DeepEqual(category, platform.Category); !op {
		t.Error(unexpectedValueString("platform category", category, platform.Category))
	}
	if op := reflect.DeepEqual(radar, platform.RadarOpts); !op {
		t.Error(unexpectedValueString("platform radar options", radar, platform.RadarOpts))
	}
	if err := testValues("platform plublic provider id", 12345, platform.PublicProviderID); err != nil {
		t.Error(err)
	}
	if err := testValues("description", "foo description", platform.Description); err != nil {
		t.Error(err)
	}
}

func TestPlatformUpdate(t *testing.T) {
	teardown := setup()
	defer teardown()
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	mux.HandleFunc("/v2/config/platforms.json/123", func(w http.ResponseWriter, r *http.Request) {
		var parsedBody map[string]interface{}
		expectedRequestData := map[string]interface{}{
			"name":                      "updated_foo_name",
			"displayName":               "updated_foo_name",
			"category":                  category,
			"radarConfig":               radar,
			"publicProviderArchetypeId": float64(67890),
			"intendedUse":               "updated foo description",
		}
		responseBodyObj := Platform{
			Id:               123,
			PlatformAlias:    "updated_foo_name",
			PlatformDispName: "updated_foo_name",
			Category:         category,
			RadarOpts:        radar,
			PublicProviderID: 67890,
			Description:      "updated foo description",
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
	createOps := NewPlatformOpts("updated_foo_name", category, radar, 67890, "updated foo description")
	platform, err := client.Platform.Update(123, &createOps)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, platform.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("platform alias", "updated_foo_name", platform.PlatformAlias); err != nil {
		t.Error(err)
	}
	if err := testValues("platform display name", "updated_foo_name", platform.PlatformDispName); err != nil {
		t.Error(err)
	}
	if op := reflect.DeepEqual(category, platform.Category); !op {
		t.Error(unexpectedValueString("platform category", category, platform.Category))
	}
	if op := reflect.DeepEqual(radar, platform.RadarOpts); !op {
		t.Error(unexpectedValueString("platform radar options", radar, platform.RadarOpts))
	}
	if err := testValues("platform plublic provider id", 67890, platform.PublicProviderID); err != nil {
		t.Error(err)
	}
	if err := testValues("description", "updated foo description", platform.Description); err != nil {
		t.Error(err)
	}
}

func TestPlatformGet(t *testing.T) {
	teardown := setup()
	defer teardown()
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	mux.HandleFunc("/v2/config/platforms.json/123", func(w http.ResponseWriter, r *http.Request) {
		responseBodyObj := Platform{
			Id:               123,
			PlatformAlias:    "foo",
			PlatformDispName: "foo",
			Category:         category,
			RadarOpts:        radar,
			PublicProviderID: 12345,
			Description:      "foo description",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})
	platform, err := client.Platform.Get(123)
	if err != nil {
		t.Error(err)
	}
	if err := testValues("id", 123, platform.Id); err != nil {
		t.Error(err)
	}
	if err := testValues("platform alias", "foo", platform.PlatformAlias); err != nil {
		t.Error(err)
	}
	if err := testValues("platform display name", "foo", platform.PlatformDispName); err != nil {
		t.Error(err)
	}
	if op := reflect.DeepEqual(category, platform.Category); !op {
		t.Error(unexpectedValueString("platform category", category, platform.Category))
	}
	if op := reflect.DeepEqual(radar, platform.RadarOpts); !op {
		t.Error(unexpectedValueString("platform radar options", radar, platform.RadarOpts))
	}
	if err := testValues("platform plublic provider id", 12345, platform.PublicProviderID); err != nil {
		t.Error(err)
	}
	if err := testValues("description", "foo description", platform.Description); err != nil {
		t.Error(err)
	}
}

func TestPlatformDelete(t *testing.T) {
	teardown := setup()
	defer teardown()
	mux.HandleFunc("/v2/config/platforms.json/123", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	err := client.Platform.Delete(123)
	if err != nil {
		t.Error(err)
	}
}

func TestPlatformList(t *testing.T) {
	teardown := setup()
	defer teardown()
	category := map[string]interface{}{"id": float64(1)}
	radar := map[string]interface{}{"usePublicData": true}
	var plaforms []Platform
	platform1 := Platform{
		Id:               123,
		PlatformAlias:    "foo",
		PlatformDispName: "foo",
		Category:         category,
		RadarOpts:        radar,
		PublicProviderID: 12345,
		Description:      "foo description",
	}
	platform2 := Platform{
		Id:               456,
		PlatformAlias:    "bar",
		PlatformDispName: "bar",
		Category:         category,
		RadarOpts:        radar,
		PublicProviderID: 67890,
		Description:      "bar description",
	}

	plaforms = append(plaforms, platform1, platform2)
	mux.HandleFunc("/v2/config/platforms.json/", func(w http.ResponseWriter, r *http.Request) {
		responseBodyObj := plaforms
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		responseBody, _ := json.Marshal(responseBodyObj)
		fmt.Fprint(w, string(responseBody))
	})

	platformlist, err := client.Platform.List()
	if err != nil {
		t.Error(err)
	}
	for index, platform := range platformlist {
		if !reflect.DeepEqual(platform, plaforms[index]) {
			t.Error(unexpectedValueString("plaforms parameter", plaforms, plaforms[index]))
		}
	}
}
