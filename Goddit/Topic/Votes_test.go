package Topic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpVote(t *testing.T) {
	v := &Votes{0, 0, 0}
	v.SetUpVote()
	assert.Equal(t, uint64(1), v.UpVotes)
	assert.Equal(t, uint64(0), v.DownVotes)
	assert.Equal(t, int64(1), v.SumVotes)
}
func TestDownVote(t *testing.T) {
	v := &Votes{0, 0, 0}
	v.SetDownVote()
	assert.Equal(t, uint64(0), v.UpVotes)
	assert.Equal(t, uint64(1), v.DownVotes)
	assert.Equal(t, int64(-1), v.SumVotes)

}
func TestSumVote(t *testing.T) {
	v := &Votes{3, 2, 0}
	assert.Equal(t, int64(1), v.GetSum())
	assert.Equal(t, int64(0), v.SumVotes)
	v.SetSumVote()
	assert.Equal(t, int64(1), v.SumVotes)
}
