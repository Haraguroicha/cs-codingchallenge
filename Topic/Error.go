package Topic

import (
	"fmt"
)

// ExceededTopicLengthError is the topic length exceeded exception
type ExceededTopicLengthError struct {
	Length int
}

func (e *ExceededTopicLengthError) Error() string {
	return fmt.Sprintf("The Topic length exceeded to %d characters", e.Length)
}

// NoTopicError is the topic not found exception
type NoTopicError struct {
	TopicID int
}

func (e *NoTopicError) Error() string {
	return fmt.Sprintf("The Topic Not Found for TopicID %d", e.TopicID)
}
