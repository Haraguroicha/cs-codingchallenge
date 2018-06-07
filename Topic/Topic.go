package Topic

import (
	"math"
	"sort"
	"unicode/utf8"

	"github.com/Haraguroicha/cs-codingchallenge/Configs"
	"github.com/Haraguroicha/cs-codingchallenge/Error"
	"github.com/Haraguroicha/cs-codingchallenge/Utilities"
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

// NewTopic is to create a new topic
func NewTopic(_topic string) (*ResponseOfTopic, error) {
	// only raise error when character count out of the config indicates value
	if uint64(utf8.RuneCountInString(_topic)) > Configs.Config.MaximumTopicLength {
		err := Error.RaiseExceededTopicLengthError(Configs.Config.MaximumTopicLength)
		return nil, err
	}

	// initial votes as starts at 0 inside the topic
	topic := &ResponseOfTopic{
		TopicTitle: _topic,
		Votes:      &Votes{0, 0, 0},
	}

	return topic, nil
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

// GetTopic from Topics
func GetTopic(_topics []*ResponseOfTopic, _topicID uint64) (*ResponseOfTopic, error) {
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

// GetMaxPage is for return the max page number for topics
func GetMaxPage(_topics []*ResponseOfTopic) uint64 {
	count := len(_topics)
	if count == 0 {
		return 0
	}
	return uint64(math.Ceil(float64(count) / float64(Configs.Config.TopicsPerPage)))
}

// sort by multi key is ref from https://golang.org/pkg/sort/

type lessFunc func(p1, p2 *ResponseOfTopic) bool

// MultiSorter implements the Sort interface, sorting the changes within.
type MultiSorter struct {
	changes []*ResponseOfTopic
	less    []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *MultiSorter) Sort(changes []*ResponseOfTopic) {
	ms.changes = changes
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *MultiSorter {
	return &MultiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *MultiSorter) Len() int {
	return len(ms.changes)
}

// Swap is part of sort.Interface.
func (ms *MultiSorter) Swap(i, j int) {
	ms.changes[i], ms.changes[j] = ms.changes[j], ms.changes[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *MultiSorter) Less(i, j int) bool {
	p, q := ms.changes[i], ms.changes[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

func sortBySumVotes(c1, c2 *ResponseOfTopic) bool {
	return c1.Votes.SumVotes > c2.Votes.SumVotes
}

func sortByUpVotes(c1, c2 *ResponseOfTopic) bool {
	return c1.Votes.UpVotes > c2.Votes.UpVotes
}
func sortByID(c1, c2 *ResponseOfTopic) bool {
	return c1.TopicID < c2.TopicID
}

// SortTopics for sort the topics by votes
func SortTopics(_topics []*ResponseOfTopic) {
	OrderedBy(sortBySumVotes, sortByUpVotes, sortByID).Sort(_topics)
}
