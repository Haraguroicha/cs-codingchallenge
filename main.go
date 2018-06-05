package main

import (
	"log"
	"math"
	"net/http"
	"os"

	"github.com/Haraguroicha/cs-codingchallenge/Configs"
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

	router := getRouter()

	router.Run(":" + port)
}

func getRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/api/getTopics", GetTopics)

	router.POST("/api/newTopic", NewTopic)

	return router
}

// QueryResponse is the Response of XHR Request structure, it always have Success field to indicate the request is success or not
type QueryResponse struct {
	Data    []interface{} `json:"data"`
	Success bool          `json:"success"`
}

// GetTopics Handler
func GetTopics(c *gin.Context) {
	// just trying to get first we want, there can not sort during users get the top list, that will be an impact to the system
	maxTopicsCount := int(math.Min(float64(Configs.Config.TopicsPerPage), float64(len(topics))))
	_topics := make([]interface{}, len(topics))
	for i, v := range topics {
		_topics[i] = v
	}
	_topics = _topics[0:maxTopicsCount]

	c.JSON(http.StatusOK, &QueryResponse{Data: _topics, Success: true})
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
	topic.TopicID = len(topics)
	topics = append(topics, topic)
	GetTopics(c)
}
