package visual

type Step int

const (
	Idle Step = iota
	MST
	Matching
	Tour
	Done
)
