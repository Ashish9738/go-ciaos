package ciaos

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func getTestConfig() Config {
	return Config{
		APIURL:        "http://test-api.com",
		UserId:        "testuser",
		UserAccessKey: "testaccesskey",
	}
}

func (m *MockClient) Put(url string, headers map[string]string, body []byte) (*http.Response, error) {
	args := m.Called(url, headers, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestCiaosClient(t *testing.T) {
	config := getTestConfig()

	ciaosClient, err := Ciaos(config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	assert.NotNil(t, ciaosClient)
}

func TestPutSuccess(t *testing.T) {
	tmpFile := "test.txt"
	err := ioutil.WriteFile(tmpFile, []byte("test data"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.Remove(tmpFile)
		if err != nil {
			t.Fatal(err)
		}
	}()

	mockClient := new(MockClient)
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("Data uploaded successfully: key = test.txt"))),
	}
	mockClient.On("Put", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)

	headers := map[string]string{"User": "testuser"}
	body := []byte("test data")
	url := "/put/test.txt"
	response, err := mockClient.Put(url, headers, body)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, 200, response.StatusCode)
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Data uploaded successfully: key = test.txt", string(respBody))
	mockClient.AssertExpectations(t)

}
