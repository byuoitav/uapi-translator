package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/byuoitav/uapi-translator/log"
	"go.uber.org/zap"
)

func GetState(url, method string, responseBody interface{}) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Log.Error("failed to create new http request", zap.String("url", url), zap.Error(err))
		return err
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
