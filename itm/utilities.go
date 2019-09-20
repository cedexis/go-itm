package itm

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"strings"
)

func unexpectedValueString(label string, expected interface{}, got interface{}) string {
	return fmt.Sprintf("Unexpected value [%s]\nExpected: %v\nGot: %v", label, expected, got)
}

func newUnexpectedValueError(label string, expected interface{}, got interface{}) error {
	return fmt.Errorf(unexpectedValueString(label, expected, got))
}

func testValues(label string, expected interface{}, got interface{}) (err error) {
	if expected != got {
		err = fmt.Errorf(unexpectedValueString(label, expected, got))
	}
	return
}

type Token struct {
	Value   string `json:"value"`
	Type    string `json:"tokenType"`
	Expired bool   `json:"expired"`
}

// Get client token for accessing ITM APIs using client ID and client Secret
func GetToken(clientId string, clientSecret string) string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	reqBody := strings.NewReader("client_id=" + clientId + "&client_secret=" + clientSecret + "&grant_type=client_credentials")
	req, err := http.NewRequest("POST", "https://api.cedexis.com/api/oauth/token", reqBody)
	if err != nil {
		logger.Error.Println("Request for Client's Bearer Token failed with error", err.Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error.Println("Request for Client's Bearer Token failed with error", err.Error())
	}
	defer httpResp.Body.Close()
	if 200 != httpResp.StatusCode {
		logger.Error.Println("Request for Bearer Token failed with error ", httpResp)
		return ""
	}
	respBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		logger.Error.Println("Request for Client's Bearer Token failed with error", err.Error())
		return ""
	}
	var resp response
	resp.Body = respBody
	var token Token
	json.Unmarshal(resp.Body, &token)
	return token.Value
}
