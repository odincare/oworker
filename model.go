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
type FailPayload struct {
	Args     []interface{} `json:"args"`
	Error    string        `json:"error"`
	Times    int           `json:"times"`
	MaxTry   int           `json:"max_try"`
	FailedAt time.Time     `json:"failed_at"`
}

type work struct {
	Queue   string    `json:"queue"`
	RunAt   time.Time `json:"run_at"`
	Payload Payload   `json:"payload"`
}

type failure struct {
	FailedAt  time.Time   `json:"failed_at"`
	Payload   Payload     `json:"payload"`
	Exception string      `json:"exception"`
	Error     string      `json:"error"`
	Times     int         `json:"times"`
	MaxTry    int         `json:"max_try"`
	Backtrace []string    `json:"backtrace"`
	Worker    *worker     `json:"worker"`
	Queue     string      `json:"queue"`
	ExecTime  []time.Time `json:"exec_time"`
}

type failureData struct {
	FailedAt  time.Time   `json:"failed_at"`
	Payload   Payload     `json:"payload"`
	Exception string      `json:"exception"`
	Error     string      `json:"error"`
	Times     int         `json:"times"`
	MaxTry    int         `json:"max_try"`
	Backtrace []string    `json:"backtrace"`
	Worker    string      `json:"worker"`
	Queue     string      `json:"queue"`
	ExecTime  []time.Time `json:"exec_time"`
}
