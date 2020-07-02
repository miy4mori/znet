package loop

import "bytes"

type process func(buffer *bytes.Buffer)

type Task interface {
	Buffer() *bytes.Buffer
	Process()
	LoadProcess(process process)
}

type task struct {
	buffer  *bytes.Buffer
	process process
}

func NewTask(buffer *bytes.Buffer) Task {
	return &task{
		buffer: buffer,
		process: func(buffer *bytes.Buffer) {
			// nothing to do
		},
	}
}

func NewTaskWithProcess(buffer *bytes.Buffer, process process) Task {
	return &task{
		buffer:  buffer,
		process: process,
	}
}

func (t *task) Buffer() *bytes.Buffer {
	return t.buffer
}

func (t *task) Process() {
	t.process(t.buffer)
}

func (t *task) LoadProcess(process process) {
	t.process = process
}
