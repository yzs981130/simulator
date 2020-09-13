package job

type Status int

const (
	Pending Status = iota
	Running
	Finishing
	Completed
)