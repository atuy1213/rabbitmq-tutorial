package model

import "time"

type Task struct {
	Name string
}

func (r *Task) DoSomething() error {
	time.Sleep(5 * time.Second)
	return nil
}
