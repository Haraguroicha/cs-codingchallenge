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
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)
		assert.Equal(t, 0, len(responsed.Data))
	})
}

// trying to add topic 1
func TestNewTopic1(t *testing.T) {
	topicTitle := "topic 1"
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: topicTitle,
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		assert.Equal(t, 1, len(responsed.Data))
		assert.Equal(t, topicTitle, responsed.Data[0].TopicTitle)
	})
}

// trying to add topic 2
func TestNewTopic2(t *testing.T) {
	topicTitle := "topic 2"
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: topicTitle,
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		assert.Equal(t, 2, len(responsed.Data))
		assert.Equal(t, topicTitle, responsed.Data[1].TopicTitle)
	})
}

// trying to add a long topic and just before the length of exceed length
func TestNewTopicLong(t *testing.T) {
	topicTitle := "123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345"
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: topicTitle,
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		assert.Equal(t, 3, len(responsed.Data))
		assert.Equal(t, topicTitle, responsed.Data[2].TopicTitle)
	})
}

// trying to add a long topic and it was exceed of the maximum length
func TestNewTopicLongAndWontSuccess(t *testing.T) {
	topicTitle := "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456"
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: topicTitle,
	}), func(w *httptest.ResponseRecorder) {
		assert.NotEqual(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, false, responsed.Success)
		assert.Equal(t, 0, len(responsed.Data))
	})
}

// trying to add a normal length topic and it should success
func TestNewTopicAsNormalAgain(t *testing.T) {
	topicTitle := "test again"
	testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
		TopicTitle: topicTitle,
	}), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		assert.Equal(t, 4, len(responsed.Data))
		assert.Equal(t, topicTitle, responsed.Data[3].TopicTitle)
	})
}

func addManysTopic(t *testing.T, startAt int, additionTest func(*testing.T, string, *QueryResponse)) {
	totalTopicsToAdd := Configs.Config.TopicsPerPage + 5
	for i := 1; i <= totalTopicsToAdd; i++ {
		topicTitle := fmt.Sprintf("manys-topic-%d", i+startAt)
		testHTTPResponse(HTTPPost("/api/newTopic", &Topic.RequestOfTopic{
			TopicTitle: topicTitle,
		}), func(w *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, w.Code)

			var responsed *QueryResponse
			json.Unmarshal(w.Body.Bytes(), &responsed)

			assert.Equal(t, true, responsed.Success)
			assert.Equal(t, math.Min(float64(Configs.Config.TopicsPerPage), float64(4+i+startAt)), float64(len(responsed.Data)))
			additionTest(t, topicTitle, responsed)
			t.Log("TopicIDs", Topic.GetTopicIDs(topics))
		})
	}
}

// trying to add manys topics and can not be responsed count greater than the paging maximum
func TestInsertManysTopics(t *testing.T) {
	addManysTopic(t, 0, func(t *testing.T, topicTitle string, responsed *QueryResponse) {
		assert.Equal(t, topicTitle, topics[len(topics)-1].TopicTitle)
	})
}

// trying to UpVote some topics
func TestUpVotes(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/upVote/3", nil), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		_topicIDs := Topic.GetTopicIDs(responsed.Data)
		assert.Equal(t, []int{3, 0}, _topicIDs[0:2])
		assert.Equal(t, &Topic.Votes{UpVotes: 1, DownVotes: 0, SumVotes: 1}, responsed.Data[0].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 0, DownVotes: 0, SumVotes: 0}, responsed.Data[1].Votes)
		t.Log("TopicIDs", Topic.GetTopicIDs(topics))
	})

	testHTTPResponse(HTTPPost("/api/upVote/4", nil), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		_topicIDs := Topic.GetTopicIDs(responsed.Data)
		assert.Equal(t, []int{3, 4, 0}, _topicIDs[0:3])
		assert.Equal(t, &Topic.Votes{UpVotes: 1, DownVotes: 0, SumVotes: 1}, responsed.Data[0].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 1, DownVotes: 0, SumVotes: 1}, responsed.Data[1].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 0, DownVotes: 0, SumVotes: 0}, responsed.Data[2].Votes)
		t.Log("TopicIDs", Topic.GetTopicIDs(topics))
	})

	testHTTPResponse(HTTPPost("/api/upVote/4", nil), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		_topicIDs := Topic.GetTopicIDs(responsed.Data)
		assert.Equal(t, []int{4, 3, 0}, _topicIDs[0:3])
		assert.Equal(t, &Topic.Votes{UpVotes: 2, DownVotes: 0, SumVotes: 2}, responsed.Data[0].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 1, DownVotes: 0, SumVotes: 1}, responsed.Data[1].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 0, DownVotes: 0, SumVotes: 0}, responsed.Data[2].Votes)
		t.Log("TopicIDs", Topic.GetTopicIDs(topics))
	})
}

// trying to add manys topics that will be append to last
func TestInsertManysTopicsAgain(t *testing.T) {
	addManysTopic(t, 25, func(t *testing.T, topicTitle string, responsed *QueryResponse) {
		assert.Equal(t, topicTitle, topics[len(topics)-1].TopicTitle)
	})
}

// trying to DownVote some topics
func TestDownVotes(t *testing.T) {
	testHTTPResponse(HTTPPost("/api/downVote/4", nil), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		_topicIDs := Topic.GetTopicIDs(responsed.Data)
		assert.Equal(t, []int{4, 3, 0}, _topicIDs[0:3])
		assert.Equal(t, &Topic.Votes{UpVotes: 2, DownVotes: 1, SumVotes: 1}, responsed.Data[0].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 1, DownVotes: 0, SumVotes: 1}, responsed.Data[1].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 0, DownVotes: 0, SumVotes: 0}, responsed.Data[2].Votes)
		t.Log("TopicIDs", Topic.GetTopicIDs(topics))
	})

	testHTTPResponse(HTTPPost("/api/downVote/4", nil), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		_topicIDs := Topic.GetTopicIDs(responsed.Data)
		assert.Equal(t, []int{3, 4, 0}, _topicIDs[0:3])
		assert.Equal(t, &Topic.Votes{UpVotes: 1, DownVotes: 0, SumVotes: 1}, responsed.Data[0].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 2, DownVotes: 2, SumVotes: 0}, responsed.Data[1].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 0, DownVotes: 0, SumVotes: 0}, responsed.Data[2].Votes)
		t.Log("TopicIDs", Topic.GetTopicIDs(topics))
	})

	testHTTPResponse(HTTPPost("/api/downVote/4", nil), func(w *httptest.ResponseRecorder) {
		assert.Equal(t, http.StatusOK, w.Code)

		var responsed *QueryResponse
		json.Unmarshal(w.Body.Bytes(), &responsed)

		assert.Equal(t, true, responsed.Success)
		_topicIDs := Topic.GetTopicIDs(responsed.Data)
		assert.Equal(t, []int{3, 0, 1}, _topicIDs[0:3])
		assert.Equal(t, 4, topics[len(topics)-1].TopicID)
		assert.Equal(t, &Topic.Votes{UpVotes: 1, DownVotes: 0, SumVotes: 1}, responsed.Data[0].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 0, DownVotes: 0, SumVotes: 0}, responsed.Data[1].Votes)
		assert.Equal(t, &Topic.Votes{UpVotes: 2, DownVotes: 3, SumVotes: -1}, topics[len(topics)-1].Votes)
		t.Log("TopicIDs", Topic.GetTopicIDs(topics))
	})
}

// trying to add manys topics that will be append to before the last one, because last one is subtotal less than 0
func TestInsertManysTopicsAgainAndAgain(t *testing.T) {
	addManysTopic(t, 50, func(t *testing.T, topicTitle string, responsed *QueryResponse) {
		assert.Equal(t, topicTitle, topics[len(topics)-2].TopicTitle) // it should be at last two to find last we added topic
	})
}
