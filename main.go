package main

import (
	"log"
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

// GetTopics Handler
func GetTopics(c *gin.Context) {
	c.JSON(http.StatusOK, topics)
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
	topics = append(topics, topic)
	c.JSON(http.StatusOK, topics)
}
