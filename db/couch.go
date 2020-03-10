package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/byuoitav/uapi-translator/log"
	"go.uber.org/zap"
)

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
