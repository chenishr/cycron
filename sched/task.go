package sched

type Task struct {
	id			 int
	taskName     string
	taskType     int
	description  string
	cronSpec     string
	concurrent   int
	command      string
	status       int
}

func (t *Task)InitTasks() (err error) {

	return nil
}