package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/byuoitav/uapi-translator/log"
	"go.uber.org/zap"
)

// ErrNotFound is the error returned by the package when a document is not found
// to fulfill a request
var ErrNotFound = errors.New("The requested document was not found")

// Service represents a database service and the config necessary to run the service
type Service struct {
	Address  string
	Username string
	Password string
}

// search represents a search parameter in a couch _find call
type search struct {
	GT    string `json:"$gt,omitempty"`
	LT    string `json:"$lt,omitempty"`
	Regex string `json:"$regex,omitempty"`
}

// query represents a query body to be sent to couch
type query struct {
	Selector map[string]interface{} `json:"selector"`
	Limit    int                    `json:"limit"`
}

func DBSearch(url, method string, query, resp interface{}) error {
	var body []byte
	var err error
	if query != nil {
		body, err = json.Marshal(query)
		if err != nil {
			log.Log.Error("failed to marshal search query into json", zap.Error(err))
			return err
		}
	}

	log.Log.Info("searching database", zap.String("method", method), zap.String("query", string(body)))
	err = makeRequest(method, url, "application/json", body, &resp)
	if err != nil {
		log.Log.Error("failed to make db search request")
		return err
	}

	return nil
}

// makeRequest makes the given request to couch and then parses the response into the
// responseBody pointer passed in
func (s *Service) makeRequest(method, path string, body []byte, responseBody interface{}) error {
	url := fmt.Sprintf("%s/%s", s.Address, path)
	log.Log.Debugf("Making couch request: %s %s", method, url)

	// Create the request
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		err = fmt.Errorf("db/makeRequest create couch request: %w", err)
		return err
	}

	// Add basic auth
	req.SetBasicAuth(s.Username, s.Password)

	// Execute the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("db/makeRequest make request: %w", err)
		return err
	}
	defer res.Body.Close()

	// Check for 404
	if res.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	// Check for non 200
	if res.StatusCode != 200 {
		return fmt.Errorf("db/makeRequest Error response from couch. Code: %d", res.StatusCode)
	}

	// Read the body
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("db/makeRequest reading response body: %w", err)
	}

	// Unmarshal
	err = json.Unmarshal(b, responseBody)
	if err != nil {
		return fmt.Errorf("db/makeRequest json unmarshal: %w", err)
	}

	return err
}

// DEPRECATED: Please use the makeRequest() function that operates on Service
// This function is left here to continue to support existing code
func makeRequest(method, url, contentType string, body []byte, responseBody interface{}) error {
	log.Log.Info("making http request", zap.String("dest-url", url))
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		log.Log.Error("failed to create new http request", zap.String("url", url), zap.Error(err))
		return err
	}

	req.SetBasicAuth(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
	if len(contentType) > 0 {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Log.Error("failed to make http request", zap.String("url", url), zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Log.Error("failed to read http response body", zap.Error(err))
		return err
	}

	if resp.StatusCode/100 != 2 {
		log.Log.Error("bad response code", zap.Int("resp code", resp.StatusCode), zap.String("body", string(b)))
		return fmt.Errorf("bad response code - %v: %s", resp.StatusCode, b)
	}

	if responseBody != nil {
		err = json.Unmarshal(b, responseBody)
		if err != nil {
			log.Log.Error("failure to unmarshal resp body", zap.String("body", string(b)), zap.Error(err))
			return err
		}
	}

	return nil
}
