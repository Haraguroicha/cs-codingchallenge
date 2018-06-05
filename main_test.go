package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
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

// for test empty topics (as initial default)
func TestGetEmptyTopics(t *testing.T) {
	testHTTPResponse(HTTPGet("/api/getTopics"), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)
		assert.Equal(t, len(responsed.Data), 0)
	})
}

// trying to add topic 1
func TestNewTopic1(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "topic 1",
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, responsed.Success, true)
		assert.Equal(t, len(responsed.Data), 1)
	})
}

// trying to add topic 2
func TestNewTopic2(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "topic 2",
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, responsed.Success, true)
		assert.Equal(t, len(responsed.Data), 2)
	})
}

// trying to add a long topic and just before the length of exceed length
func TestNewTopicLong(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, responsed.Success, true)
		assert.Equal(t, len(responsed.Data), 3)
	})
}

// trying to add a long topic and it was exceed of the maximum length
func TestNewTopicLongAndWontSuccess(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456",
	}), func(w *httptest.ResponseRecorder) {
		assert.NotEqual(t, w.Code, http.StatusOK)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, responsed.Success, false)
		assert.Equal(t, len(responsed.Data), 0)
	})
}

// trying to add a normal length topic and it should success
func TestNewTopicAsNormalAgain(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: "test again",
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, w.Code, http.StatusOK)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, responsed.Success, true)
		assert.Equal(t, len(responsed.Data), 4)
	})
}

// trying to add manys topics and can not be responsed count greater than the paging maximum
func TestInsertManysTopics(t *testing.T) {
	totalTopicsToAdd := Configs.Config.TopicsPerPage + 5
	for i := 1; i <= totalTopicsToAdd; i++ {
		testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
			TopicTitle: fmt.Sprintf("manys-topic-%d", i),
		}), func(w *httptest.ResponseRecorder) {
			assert.Equal(t, w.Code, http.StatusOK)

			var responsed *QueryResponse
			json.Unmarshal(w.Body.Bytes(), &responsed)

			assert.Equal(t, responsed.Success, true)
			assert.Equal(t, float64(len(responsed.Data)), math.Min(float64(Configs.Config.TopicsPerPage), float64(4+i)))
		})
	}
}