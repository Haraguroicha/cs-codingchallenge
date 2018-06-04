package Topic

// The Votes structure
type Votes struct {
	UpVotes   int `json:"upVotes"`
	DownVotes int `json:"downVotes"`
	SumVotes  int `json:"sumVotes"`
}

// GetSum is for calculate the Sum to SumVotes value
func (v *Votes) GetSum() int {
	return v.UpVotes - v.DownVotes
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
