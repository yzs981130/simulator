package job

type Job struct {
	ID				int
	// resource request
	CPU				int
	Mem 			int
	GPU 			int

	// relative submit time
	SubmitTime		int
	// relative running time
	CompletionTime		int

	//pending->running time
	StartTime		int

	// user group info
	User 		int

	// runtime info
	WaitingTime		int
	RunningTime		int

	// last term status
	// update when session ends
	Status			Status

	// occupied gpu info
	GPUIndex		[]int

	// custom field

}

func (job *Job) NewJob(_id int, _CPU int, _Mem int, _GPU int, _submitTime int, _completionTime int, _user int) Job {
	j := Job{
		ID: _id,
		CPU: _CPU,
		Mem: _Mem,
		GPU: _GPU,
		SubmitTime: _submitTime,
		CompletionTime: _completionTime,
		User: _user,
		WaitingTime: 0,
		RunningTime: 0,
		Status: Pending,
	}
	j.GPUIndex = make([]int, _GPU)
	return j
}


