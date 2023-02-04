package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// TryUnmarshalJSON attempts to call json.Unmarshal on data and handles errors.
func TryUnmarshalJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("TryUnmarshalJson : error unmarshaling JSON: %v", err)
	}
	return nil
}

// TryGet attempts to make a GET request to the specified URL and either returns the body or handles error creation.
func TryGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %s", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %s", err)
	}

	return body, nil
}

// TryGet attempts to make a POST request to the specified URL and either returns the body or handles error creation.
func TryPostJSON(url string, payload interface{}) ([]byte, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error while marshaling payload: %s", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("POST request failed: %s\n", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %s", err)
	}

	return body, nil
}
