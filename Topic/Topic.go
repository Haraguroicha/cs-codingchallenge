package Topic

import (
	"github.com/Haraguroicha/cs-codingchallenge/Configs"
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
		err := &ExceededTopicLengthError{Configs.Config.MaximumTopicLength}
		return nil, err
	}

	votes := &Votes{0, 0, 0}
	topic := &ResponseOfTopic{
		TopicTitle: _topic,
		Votes:      votes,
	}

	return topic, nil
}
