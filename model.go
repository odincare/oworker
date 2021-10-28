package oworker

import "time"

type Payload struct {
	Class string        `json:"class"`
	Args  []interface{} `json:"args"`
}

type Job struct {
	Queue   string
	Payload Payload
}

type work struct {
	Queue   string    `json:"queue"`
	RunAt   time.Time `json:"run_at"`
	Payload Payload   `json:"payload"`
}

type failure struct {
	FailedAt  time.Time `json:"failed_at"`
	Payload   Payload   `json:"payload"`
	Exception string    `json:"exception"`
	Error     string    `json:"error"`
	Backtrace []string  `json:"backtrace"`
	Worker    *worker   `json:"worker"`
	Queue     string    `json:"queue"`
}

//处理函数定义
type workerFunc func(string, ...interface{}) error
