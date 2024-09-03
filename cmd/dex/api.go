package dex

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func RunGetMethod(data interface{}) (*ExecutionResult, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))
	buf := bytes.NewBuffer(body)
	resp, err := http.Post("https://toncenter.com/api/v3/runGetMethod", "application/json", buf)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	result := &ExecutionResult{}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resBody, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
