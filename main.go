package main

import (
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/Haraguroicha/cs-codingchallenge/Configs"
	"github.com/Haraguroicha/cs-codingchallenge/Error"
	"github.com/Haraguroicha/cs-codingchallenge/Topic"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

// topics for shared variable between main program and unit test
var topics []*Topic.ResponseOfTopic

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// set config read from conf/config.yaml to Congigs package shared variable
	Configs.Config = Configs.NewConfig("conf/config.yaml")

	// initial topics as empty array
	topics = []*Topic.ResponseOfTopic{}

	router := getRouter(false)

	router.Run(":" + port)
}

func getRouter(isTest bool) *gin.Engine {
	router := gin.New()
	if isTest == false {
		router.Use(gin.Logger())
	}
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		// pass Config data to template
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"Config": Configs.Config})
	})

	router.GET("/api/getTopics/*page", GetTopics)

	router.POST("/api/newTopic", NewTopic)

	router.POST("/api/upVote/:topic/*page", UpTopic)

	router.POST("/api/downVote/:topic/*page", DownTopic)

	return router
}

// Pages for indicates the response position
type Pages struct {
	CurrentPage uint64 `json:"currentPage"`
	LastPage    uint64 `json:"lastPage"`
}

// QueryResponse is the Response of XHR Request structure,
// it always have Success field to indicate the request is success or not,
// and there have Pages structure to indicate current page and last page count
type QueryResponse struct {
	Data    []*Topic.ResponseOfTopic `json:"data"`
	Pages   *Pages                   `json:"pages"`
	Success bool                     `json:"success"`
}

// GetTopics Handler
func GetTopics(c *gin.Context) {
	var _page = c.Param("page")
	// param read page always prefix a slash `/`, we need to strip that
	if len(_page) <= 1 {
		_page = "1"
	} else {
		_page = _page[1:len(_page)]
	}
	page, err := strconv.ParseUint(_page, 10, 64)
	// after pre-process the page parameter, convert to int
	// if there has some error, e.g. non-numeric character included, raise the error
	if err != nil {
		err := Error.RaisePageParameterInvalidError(_page)
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	// if there has less then 1, e.g. -1 or 0, also raise the error
	if page <= 0 {
		err := Error.RaisePageParameterInvalidError(_page)
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	// just trying to get first we want,
	// because we can not sort the data structure during users get the top list,
	// that will be an impact to the system
	starts := Configs.Config.TopicsPerPage * (page - 1)
	maxTopicsCount := uint64(math.Min(float64(starts+Configs.Config.TopicsPerPage), float64(len(topics))))
	// if we request a page is out of bound, raise error
	if starts > uint64(len(topics)) {
		err := Error.RaisePageInvalidError(_page)
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	// strip the range only we want
	_topics := topics[starts:maxTopicsCount]

	c.JSON(http.StatusOK, &QueryResponse{
		Data: _topics,
		Pages: &Pages{
			CurrentPage: page,
			LastPage:    Topic.GetMaxPage(topics),
		},
		Success: true,
	})
}

// NewTopic Handler
func NewTopic(c *gin.Context) {
	var req Topic.RequestOfTopic
	c.BindJSON(&req)
	topic, err := Topic.NewTopic(req.TopicTitle)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	topic.TopicID = uint64(len(topics)) + 1
	topics = append(topics, topic)
	Topic.SortTopics(topics)
	GetTopics(c)
}

func getTopicByRequest(c *gin.Context) *Topic.ResponseOfTopic {
	topicID, err := strconv.ParseUint(c.Param("topic"), 10, 64)
	if err != nil {
		err := Error.RaiseTopicParameterInvalidError(c.Param("topic"))
		c.JSON(http.StatusExpectationFailed, err)
		return nil
	}
	topic, err := Topic.GetTopic(topics, topicID)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err)
		return nil
	}
	return topic
}

// UpTopic Handler
func UpTopic(c *gin.Context) {
	topic := getTopicByRequest(c)
	if topic != nil {
		topic.Votes.SetUpVote()
		Topic.SortTopics(topics)
		GetTopics(c)
	}
}

// DownTopic Handler
func DownTopic(c *gin.Context) {
	topic := getTopicByRequest(c)
	if topic != nil {
		topic.Votes.SetDownVote()
		Topic.SortTopics(topics)
		GetTopics(c)
	}
}
