package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Haraguroicha/cs-codingchallenge/Configs"
	"github.com/Haraguroicha/cs-codingchallenge/Topic"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine

func init() {
	Configs.Config = Configs.NewConfig("conf/config.yaml")
	r = getRouter()
}

// Helper function to process a request and test its response
func testHTTPResponse(req *http.Request, f func(w *httptest.ResponseRecorder)) {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	f(w)
}

func HTTPRequest(_method string, _url string, _data io.Reader) *http.Request {
	req, _ := http.NewRequest(_method, _url, _data)
	return req
}

func HTTPGet(_url string) *http.Request {
	return HTTPRequest("GET", _url, nil)
}

func HTTPPost(_url string, _data interface{}) *http.Request {
	jsonPayload, _ := json.Marshal(_data)
	return HTTPRequest("POST", _url, bytes.NewBuffer(jsonPayload))
}

func TestGetEmptyTopics(t *testing.T) {
	testHTTPResponse(HTTPGet("/api/getTopics"), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed []*Topic.ResponseOfTopic
		json.Unmarshal(w.Body.Bytes(), &responsed)
		assert.Equal(t, len(responsed), 0)
	})
}

func TestNewTopic1(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "topic 1",
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed []*Topic.ResponseOfTopic
		json.Unmarshal(w.Body.Bytes(), &responsed)
		assert.Equal(t, len(responsed), 1)
	})
}

func TestNewTopic2(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "topic 2",
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed []*Topic.ResponseOfTopic
		json.Unmarshal(w.Body.Bytes(), &responsed)
		assert.Equal(t, len(responsed), 2)
	})
}

func TestNewTopicLong(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed []*Topic.ResponseOfTopic
		json.Unmarshal(w.Body.Bytes(), &responsed)
		assert.Equal(t, len(responsed), 3)
	})
}

func TestNewTopicLongAndWontSuccess(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456",
	}), func(w *httptest.ResponseRecorder) {
		assert.NotEqual(t, w.Code, http.StatusOK)

		var responsed []*Topic.ResponseOfTopic
		json.Unmarshal(w.Body.Bytes(), &responsed)
		assert.Equal(t, len(responsed), 0)
	})
}

func TestNewTopicAsNormalAgain(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "test again",
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed []*Topic.ResponseOfTopic
		json.Unmarshal(w.Body.Bytes(), &responsed)
		assert.Equal(t, len(responsed), 4)
	})
}
