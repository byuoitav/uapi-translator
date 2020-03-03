package couch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func DBSearch(url, method string, query, resp interface{}) error {
	var body []byte
	var err error
	if query != nil {
		body, err = json.Marshal(query)
		if err != nil {
			return err
		}
	}

	err = makeRequest(method, url, "application/json", body, &resp)
	if err != nil {
		return err
	}

	return nil
}

func makeRequest(method, url, contentType string, body []byte, responseBody interface{}) error {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.SetBasicAuth(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
	if len(contentType) > 0 {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("bad response code - %v: %s", resp.StatusCode, b)
	}

	if responseBody != nil {
		err = json.Unmarshal(b, responseBody)
		if err != nil {
			return err
		}
	}

	return nil
}
