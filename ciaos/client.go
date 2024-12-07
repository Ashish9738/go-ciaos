package ciaos

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Ashish9738/go-ciaos/config"
	flatbufferHandler "github.com/Ashish9738/go-ciaos/utils/handlers"
)

type CiaosInstance struct {
	config  config.Config
	headers map[string]string
}

func Ciaos(config config.Config) (*CiaosInstance, error) {
	if config.UserId == "" {
		return nil, fmt.Errorf("user id must not be empty")
	}

	if config.APIURL == "" {
		return nil, fmt.Errorf("api url must not be empty")
	}

	headers := map[string]string{
		"User": config.UserId,
	}

	return &CiaosInstance{
		config:  config,
		headers: headers,
	}, nil
}

func (c *CiaosInstance) Put(filePath string, key string) (*http.Response, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file_path cannot be empty or None")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	if key == "" {
		key = filePath
	}

	flatBufferData, err := flatbufferHandler.CreateFlatBuffer([][]byte{data})
	if err != nil {
		return nil, fmt.Errorf("failed to create a FlatBuffers Data: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/put/%s", c.config.APIURL, key), bytes.NewReader(flatBufferData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	for headerKey, headerValue := range c.headers {
		req.Header.Set(headerKey, headerValue)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP PUT req failed: %v", err)
	}

	return resp, nil

}

func (c *CiaosInstance) PutBinary(key string, dataList [][]byte) (*http.Response, error) {
	flatBufferData, err := flatbufferHandler.CreateFlatBuffer(dataList)
	if err != nil {
		return nil, fmt.Errorf("failed to create FlatBuffer Data: %v", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/put/%s", c.config.APIURL, key), bytes.NewReader(flatBufferData))

	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	for headerKey, headerValue := range c.headers {
		req.Header.Set(headerKey, headerValue)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP PUT req failed: %v", err)
	}

	return resp, nil
}

func (c *CiaosInstance) UpdateKey(oldKey string, newKey string) (string, error) {

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/update_key/%s/%s", c.config.APIURL, oldKey, newKey), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create post request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http error during key update: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read the request body: %v", err)
	}

	return string(body), nil
}

func (c *CiaosInstance) Update(key string, dataList [][]byte) (*http.Response, error) {
	flatBufferData, err := flatbufferHandler.CreateFlatBuffer(dataList)
	if err != nil {
		return nil, fmt.Errorf("failed to create FlatBuffer Data: %v", err)
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/update/%s", c.config.APIURL, key),
		bytes.NewReader(flatBufferData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http error during update: %v", err)
	}

	return resp, nil
}

func (c *CiaosInstance) Append(key string, dataList [][]byte) (*http.Response, error) {
	flatBufferData, err := flatbufferHandler.CreateFlatBuffer(dataList)
	if err != nil {
		return nil, fmt.Errorf("failed to create FlatBuffer Data: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/append/%s", c.config.APIURL, key),
		bytes.NewReader(flatBufferData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("http error during append: %v", err)
	}

	return resp, nil
}

func (c *CiaosInstance) Delete(key string) (*http.Response, error) {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/delete/%s", c.config.APIURL, key),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %v", err)
	}

	for headerKey, headerValue := range c.headers {
		req.Header.Set(headerKey, headerValue)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error during deletion: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("failed to delete key '%s': status code %d", key, resp.StatusCode)
	}

	return resp, nil
}

func (c *CiaosInstance) Get(key string) ([][]byte, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/get/%s", c.config.APIURL, key),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error during retrieval: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	fileDataList, err := flatbufferHandler.ParseFlatBuffer(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing FlatBuffer: %v", err)
	}

	return fileDataList, nil

}
