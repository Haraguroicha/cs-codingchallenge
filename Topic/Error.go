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
