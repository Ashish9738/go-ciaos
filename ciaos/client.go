package ciaos

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	flatbufferHandler "github.com/Ashish9738/go-ciaos/utils/handlers"
)

func Ciaos(config Config) (*Config, error) {
	if config.UserId == "" {
		return nil, fmt.Errorf("user id must not be empty")
	}

	if config.APIURL == "" {
		return nil, fmt.Errorf("api url must not be empty")
	}

	return &config, nil
}

func (config *Config) Put(filePath string, key string) (*http.Response, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file_path cannot be empty")
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
	fmt.Println("Received key:", key)

	if key == "" {
		key = filepath.Base(filePath)
		fmt.Println("without key passing", key)
	}

	flatBufferData, err := flatbufferHandler.CreateFlatBuffer([][]byte{data})
	if err != nil {
		return nil, fmt.Errorf("failed to create FlatBuffer Data: %v", err)
	}

	fmt.Println("with key passing", key)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/put/%s", config.APIURL, key), bytes.NewReader(flatBufferData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}

	req.Header.Set("User", config.UserId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP PUT request failed: %v", err)
	}

	return resp, nil
}

func (config *Config) PutBinary(key string, dataList [][]byte) (*http.Response, error) {
	flatBufferData, err := flatbufferHandler.CreateFlatBuffer(dataList)
	if err != nil {
		return nil, fmt.Errorf("failed to create FlatBuffer Data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/put/%s", config.APIURL, key), bytes.NewReader(flatBufferData))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	req.Header.Set("User", config.UserId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP POST request failed: %v", err)
	}

	return resp, nil
}

func (config *Config) UpdateKey(oldKey string, newKey string) (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/update_key/%s/%s", config.APIURL, oldKey, newKey), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create POST request: %v", err)
	}

	req.Header.Set("User", config.UserId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP error during key update: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}

func (config *Config) Update(key string, dataList [][]byte) (*http.Response, error) {
	flatBufferData, err := flatbufferHandler.CreateFlatBuffer(dataList)
	if err != nil {
		return nil, fmt.Errorf("failed to create FlatBuffer Data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/update/%s", config.APIURL, key), bytes.NewReader(flatBufferData))
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %v", err)
	}

	req.Header.Set("User", config.UserId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error during update: %v", err)
	}

	return resp, nil
}

func (config *Config) Append(key string, dataList [][]byte) (*http.Response, error) {
	flatBufferData, err := flatbufferHandler.CreateFlatBuffer(dataList)
	if err != nil {
		return nil, fmt.Errorf("failed to create FlatBuffer Data: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/append/%s", config.APIURL, key), bytes.NewReader(flatBufferData))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %v", err)
	}

	req.Header.Set("User", config.UserId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error during append: %v", err)
	}

	return resp, nil
}

func (config *Config) Delete(key string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/delete/%s", config.APIURL, key), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %v", err)
	}

	req.Header.Set("User", config.UserId)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP error during deletion: %v", err)
	}

	return resp, nil
}

func (config *Config) Get(key string) ([][]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/get/%s", config.APIURL, key), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	req.Header.Set("User", config.UserId)

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
