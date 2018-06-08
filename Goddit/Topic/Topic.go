package Topic

import (
	"github.com/Haraguroicha/cs-codingchallenge/Goddit/Error"
	"github.com/Haraguroicha/cs-codingchallenge/Goddit/Utilities"
)

// RequestOfTopic struct
type RequestOfTopic struct {
	TopicTitle string `json:"topicTitle"`
}

// ResponseOfTopic struct
type ResponseOfTopic struct {
	TopicID    uint64 `json:"topicId"`
	TopicTitle string `json:"topicTitle"`
	Votes      *Votes `json:"votes"`
}

// GetTopicIDs for only map the TopicID as array
func GetTopicIDs(_topics []*ResponseOfTopic) []uint64 {
	// read id as interface array
	_ids := Utilities.Map(_topics, func(val interface{}) interface{} {
		return val.(*ResponseOfTopic).TopicID
	}).([]interface{})
	// convert back to int array
	_topicIDs := make([]uint64, len(_ids))
	for i, v := range _ids {
		_topicIDs[i] = v.(uint64)
	}
	// return the int array
	return _topicIDs
}

// GetTopicByID from Topics
func GetTopicByID(_topics []*ResponseOfTopic, _topicID uint64) (*ResponseOfTopic, error) {
	// find the topic by id
	for _, t := range _topics {
		if t.TopicID == _topicID {
			return t, nil
		}
	}
	// only raise error when we not found the requested topic id
	err := Error.RaiseNoTopicError(_topicID)
	return nil, err
}
