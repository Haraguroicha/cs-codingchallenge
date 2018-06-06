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

var topics []*Topic.ResponseOfTopic

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	Configs.Config = Configs.NewConfig("conf/config.yaml")

	topics := []*Topic.ResponseOfTopic{}

	log.Printf("Topic Count: %d", len(topics))

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
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/api/getTopics/*page", GetTopics)

	router.POST("/api/newTopic", NewTopic)

	router.POST("/api/upVote/:topic/*page", UpTopic)

	router.POST("/api/downVote/:topic/*page", DownTopic)

	return router
}

// Pages for indicates the response position
type Pages struct {
	CurrentPage int `json:"currentPage"`
	LastPage    int `json:"lastPage"`
}

// QueryResponse is the Response of XHR Request structure, it always have Success field to indicate the request is success or not
type QueryResponse struct {
	Data    []*Topic.ResponseOfTopic `json:"data"`
	Pages   *Pages                   `json:"pages"`
	Success bool                     `json:"success"`
}

// GetTopics Handler
func GetTopics(c *gin.Context) {
	var _page = c.Param("page")
	if len(_page) <= 1 {
		_page = "1"
	} else {
		_page = _page[1:len(_page)]
	}
	page, err := strconv.Atoi(_page)
	if page <= 0 {
		err := Error.RaisePageParameterInvalidError(_page)
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	if err != nil {
		err := Error.RaisePageParameterInvalidError(_page)
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	// just trying to get first we want, there can not sort during users get the top list, that will be an impact to the system
	starts := Configs.Config.TopicsPerPage * (page - 1)
	maxTopicsCount := int(math.Min(float64(starts+Configs.Config.TopicsPerPage), float64(len(topics))))
	if starts > len(topics) {
		err := Error.RaisePageInvalidError(_page)
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
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
	topic.TopicID = len(topics) + 1
	topics = append(topics, topic)
	Topic.SortTopics(topics)
	GetTopics(c)
}

// UpTopic Handler
func UpTopic(c *gin.Context) {
	topicID, err := strconv.Atoi(c.Param("topic"))
	if err != nil {
		err := Error.RaiseTopicParameterInvalidError(c.Param("topic"))
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	topic, err := Topic.GetTopic(topics, topicID)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	topic.Votes.SetUpVote()
	Topic.SortTopics(topics)
	GetTopics(c)
}

// DownTopic Handler
func DownTopic(c *gin.Context) {
	topicID, err := strconv.Atoi(c.Param("topic"))
	if err != nil {
		err := Error.RaiseTopicParameterInvalidError(c.Param("topic"))
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	topic, err := Topic.GetTopic(topics, topicID)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	topic.Votes.SetDownVote()
	Topic.SortTopics(topics)
	GetTopics(c)
}
