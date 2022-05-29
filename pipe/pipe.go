package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	Pipe struct {
		Ctx
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			DefaultTask(&TL).Job(),
		),
	)
}
