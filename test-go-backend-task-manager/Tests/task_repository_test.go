package tests

import (
	"context"

	"github.com/stretchr/testify/suite"
)

type testRepositorySuite struct {
	suite.Suite
}

func (taskRepo *taskRepository) FetchAllTasks(cxt context.Context) ([]domain.Task, *domain.TaskError)
