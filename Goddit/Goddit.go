package Goddit

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/Haraguroicha/cs-codingchallenge/Goddit/Configs"
	"github.com/Haraguroicha/cs-codingchallenge/Goddit/Error"
	"github.com/Haraguroicha/cs-codingchallenge/Goddit/Topic"
	"github.com/gin-gonic/gin"
)

// Goddit is an instance
type Goddit struct {
	*http.Server
	Config *Configs.Configs
	Topics []*Topic.ResponseOfTopic
	Router *gin.Engine
	Port   string
	isTest bool
}

// NewService create server instance
func NewService(isTest bool, port string, configFile string) *Goddit {
	return &Goddit{
		Config: Configs.NewConfig(configFile),
		Topics: []*Topic.ResponseOfTopic{},
		Router: gin.New(),
		Port:   port,
		isTest: isTest,
	}
}

// Start the server
func (g *Goddit) Start() {
	GetRouter(g, g.Router)

	g.Server = &http.Server{
		Addr:    ":" + g.Port,
		Handler: g.Router,
	}

	go func() {
		// service connections
		if err := g.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// Stop the server
func (g *Goddit) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := g.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

// GetRouter is for get the default gin routes
func GetRouter(g *Goddit, r *gin.Engine) {
	if g.isTest == false {
		r.Use(gin.Logger())
	}
	r.LoadHTMLGlob("templates/*.tmpl.html")
	r.Static("/static", "static")

	r.GET("/", func(c *gin.Context) {
		// pass Config data to template
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"Config": g.Config})
	})

	r.GET("/api/getTopics/*page", g.GetTopics)

	r.POST("/api/newTopic", g.NewTopic)

	r.POST("/api/upVote/:topic/*page", g.UpTopic)

	r.POST("/api/downVote/:topic/*page", g.DownTopic)
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

// GetMaxPage is for return the max page number for topics
func (g *Goddit) GetMaxPage(_topics []*Topic.ResponseOfTopic) uint64 {
	count := len(_topics)
	if count == 0 {
		return 0
	}
	return uint64(math.Ceil(float64(count) / float64(g.Config.TopicsPerPage)))
}

// GetTopics Handler
func (g *Goddit) GetTopics(c *gin.Context) {
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
	starts := g.Config.TopicsPerPage * (page - 1)
	maxTopicsCount := uint64(math.Min(float64(starts+g.Config.TopicsPerPage), float64(len(g.Topics))))
	// if we request a page is out of bound, raise error
	if starts > uint64(len(g.Topics)) {
		err := Error.RaisePageInvalidError(_page)
		c.JSON(http.StatusExpectationFailed, err)
		return
	}
	// strip the range only we want
	_topics := g.Topics[starts:maxTopicsCount]

	c.JSON(http.StatusOK, &QueryResponse{
		Data: _topics,
		Pages: &Pages{
			CurrentPage: page,
			LastPage:    g.GetMaxPage(g.Topics),
		},
		Success: true,
	})
}

// NewTopic Handler
func (g *Goddit) NewTopic(c *gin.Context) {
	var req Topic.RequestOfTopic
	c.BindJSON(&req)

	// only raise error when character count out of the config indicates value
	if uint64(utf8.RuneCountInString(req.TopicTitle)) > g.Config.MaximumTopicLength {
		err := Error.RaiseExceededTopicLengthError(g.Config.MaximumTopicLength)
		if err != nil {
			c.JSON(http.StatusExpectationFailed, err)
			return
		}
	}
	// initial votes as starts at 0 inside the topic
	topic := &Topic.ResponseOfTopic{
		TopicTitle: req.TopicTitle,
		Votes:      &Topic.Votes{0, 0, 0},
	}

	topic.TopicID = uint64(len(g.Topics)) + 1
	g.Topics = append(g.Topics, topic)
	Topic.SortTopics(g.Topics)
	g.GetTopics(c)
}

func (g *Goddit) getTopicByRequest(c *gin.Context) *Topic.ResponseOfTopic {
	topicID, err := strconv.ParseUint(c.Param("topic"), 10, 64)
	if err != nil {
		err := Error.RaiseTopicParameterInvalidError(c.Param("topic"))
		c.JSON(http.StatusExpectationFailed, err)
		return nil
	}
	topic, err := Topic.GetTopicByID(g.Topics, topicID)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err)
		return nil
	}
	return topic
}

// UpTopic Handler
func (g *Goddit) UpTopic(c *gin.Context) {
	topic := g.getTopicByRequest(c)
	if topic != nil {
		topic.Votes.SetUpVote()
		Topic.SortTopics(g.Topics)
		g.GetTopics(c)
	}
}

// DownTopic Handler
func (g *Goddit) DownTopic(c *gin.Context) {
	topic := g.getTopicByRequest(c)
	if topic != nil {
		topic.Votes.SetDownVote()
		Topic.SortTopics(g.Topics)
		g.GetTopics(c)
	}
}
