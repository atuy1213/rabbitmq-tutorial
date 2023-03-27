package repository

import (
	"fmt"

	"github.com/atuy1213/rabbitmq-tutorial/pkg/model"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetTasks() []*model.Task {
	ret := []*model.Task{}
	for i := 0; i < 5; i++ {
		ret = append(ret, &model.Task{
			Name: fmt.Sprintf("task%d", i),
		})
	}
	return ret
}
