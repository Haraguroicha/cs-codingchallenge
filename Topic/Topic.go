package Topic

import (
	"github.com/Haraguroicha/cs-codingchallenge/Configs"
	"github.com/Haraguroicha/cs-codingchallenge/Utilities"
)

// RequestOfTopic struct
type RequestOfTopic struct {
	TopicTitle string `json:"topicTitle"`
}

// ResponseOfTopic struct
type ResponseOfTopic struct {
	TopicID    int    `json:"topicId"`
	TopicTitle string `json:"topicTitle"`
	Votes      *Votes `json:"votes"`
}

// NewTopic is to create a new topic
func NewTopic(_topic string) (*ResponseOfTopic, error) {
	if len(_topic) > Configs.Config.MaximumTopicLength {
		err := &ExceededTopicLengthError{Length: Configs.Config.MaximumTopicLength}
		return nil, err
	}

	votes := &Votes{0, 0, 0}
	topic := &ResponseOfTopic{
		TopicTitle: _topic,
		Votes:      votes,
	}

	return topic, nil
}

// GetTopicIDs for only map the TopicID as array
func GetTopicIDs(_topics []*ResponseOfTopic) []int {
	_ids := Utilities.Map(_topics, func(val interface{}) interface{} {
		return val.(*ResponseOfTopic).TopicID
	}).([]interface{})
	_topicIDs := make([]int, len(_ids))
	for i, v := range _ids {
		_topicIDs[i] = v.(int)
	}
	return _topicIDs
}

// SortTopics for sort the topics by votes
func SortTopics(_topics []*ResponseOfTopic) []*ResponseOfTopic {
	return _topics
}
