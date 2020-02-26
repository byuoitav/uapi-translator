package couch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// make query object

// add query object and method as parameters
func RequestRoom() {

}

func makeRequest(method, url, contentType string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
	if len(contentType) > 0 {
		req.Header.Add("Content-Type", contentType)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("bad response code - %v: %s", resp.StatusCode, bodyBytes)
	}

	return bodyBytes, nil
}
