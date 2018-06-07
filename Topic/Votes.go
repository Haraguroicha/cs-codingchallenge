package Topic

// The Votes structure
type Votes struct {
	UpVotes   uint64 `json:"upVotes"`
	DownVotes uint64 `json:"downVotes"`
	SumVotes  int64  `json:"sumVotes"`
}

// GetSum is for calculate the Sum to SumVotes value
func (v *Votes) GetSum() int64 {
	return int64(v.UpVotes - v.DownVotes)
}

// SetSumVote is for set the sum value from GetSum
func (v *Votes) SetSumVote() {
	v.SumVotes = v.GetSum()
}

// SetUpVote is for update the UpVotes value and also update the SumVotes
func (v *Votes) SetUpVote() {
	v.UpVotes++
	v.SetSumVote()
}

// SetDownVote is for update the DownVotes value and also update the SumVotes
func (v *Votes) SetDownVote() {
	v.DownVotes++
	v.SetSumVote()
}
