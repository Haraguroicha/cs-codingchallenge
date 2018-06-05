package Error

import (
	"fmt"
)

// ErrorMessage is the topic length exceeded exception
type ErrorMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Message)
}

// RaiseExceededTopicLengthError for raise ExceededTopicLengthError
func RaiseExceededTopicLengthError(length int) *ErrorMessage {
	return &ErrorMessage{
		Title:   "ExceededTopicLengthError",
		Message: fmt.Sprintf("The Topic length exceeded to %d characters", length),
		Success: false,
	}
}

// RaiseNoTopicError for raise NoTopicError
func RaiseNoTopicError(topicID int) *ErrorMessage {
	return &ErrorMessage{
		Title:   "NoTopicError",
		Message: fmt.Sprintf("The Topic Not Found for TopicID %d", topicID),
		Success: false,
	}
}

// RaiseNoTopicParameterFoundError for raise NoTopicParameterFoundError
func RaiseNoTopicParameterFoundError() *ErrorMessage {
	return &ErrorMessage{
		Title:   "NoTopicParameterFoundError",
		Message: "There is not found for parameter to indicates the TopicID",
		Success: false,
	}
}

// RaiseTopicParameterInvalidError for raise TopicParameterInvalidError
func RaiseTopicParameterInvalidError(invalidTopicID string) *ErrorMessage {
	return &ErrorMessage{
		Title:   "TopicParameterInvalidError",
		Message: fmt.Sprintf("There is invalid TopicID of %s", invalidTopicID),
		Success: false,
	}
}
