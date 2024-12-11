package ciaos_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ashish9738/go-ciaos"
)

func testConfig() *ciaos.Config {
	return &ciaos.Config{
		APIURL:        "http://test-api.com",
		UserId:        "testuser",
		UserAccessKey: "testaccesskey",
	}
}

func TestPutSuccess(t *testing.T) {
	testData := []byte("test data")
	tmpFile := "test.txt"

	tmpFilePath := "./" + tmpFile
	err := ioutil.WriteFile(tmpFilePath, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	defer os.Remove(tmpFilePath)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := filepath.Base(r.URL.Path)

		if fileName != tmpFile {
			t.Errorf("Unexpected file name in URL path: %s", fileName)
		}

		expectedHeaders := "testuser"
		if r.Header.Get("User") != expectedHeaders {
			t.Errorf("Expected header User: %s, got: %s", expectedHeaders, r.Header.Get("User"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data uploaded successfully: key = test.txt"))
	}))
	defer mockServer.Close()

	cfg := testConfig()
	cfg.APIURL = mockServer.URL
	ciaos := cfg
	fileName := filepath.Base(tmpFilePath)
	response, err := ciaos.Put(fileName, tmpFilePath)
	if err != nil {
		t.Fatalf("Failed to perform PUT request: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	defer response.Body.Close()

	expectedResponse := "Data uploaded successfully: key = test.txt"
	if string(body) != expectedResponse {
		t.Errorf("Expected response: %s, got: %s", expectedResponse, string(body))
	}
}
