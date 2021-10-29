package oworker

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

var (
	failHandlers map[string]failedFunc
)

//处理函数定义
type failedFunc func(string, FailPayload) error

func init() {
	failHandlers = make(map[string]failedFunc)
}

func RegisterFail(class string, handler failedFunc) {
	failHandlers[class] = handler
}
func FailedInfo2String() string {
	return ""
}

type failedHandler struct{}

func newFailedHandle() *failedHandler {
	return &failedHandler{}
}
func (p *failedHandler) handle(jobs <-chan *failureData) {
	conn, err := GetConn()
	if err != nil {
		logger.Criticalf("Error on getting connection in failed handler : %v", err)
		return
	} else {
		PutConn(conn)
	}

	go func() {
		defer func() {
			conn, err := GetConn()
			if err != nil {
				return
			} else {
				PutConn(conn)
			}
		}()
		for job := range jobs {
			if hand, ok := failHandlers[job.Payload.Class]; ok {
				if job.Times != 0 && job.MaxTry != 0 {
					if job.Times >= job.MaxTry {
						continue
					}
				}
				data := FailPayload{
					Args:     job.Payload.Args,
					Error:    job.Error,
					Times:    job.Times,
					MaxTry:   job.MaxTry,
					FailedAt: job.FailedAt,
				}
				if err := hand(job.Queue, data); err != nil {
					job.Error = job.Error + "/" + err.Error()
					job.Times = job.Times + 1
					job.ExecTime = append(job.ExecTime, time.Now())
					writeData2Redis(conn, job)
				}
			} else {
				writeData2Redis(conn, job)
			}
		}
	}()
}
func writeData2Redis(conn *RedisConn, data *failureData) {
	buf, _ := json.Marshal(failedData2Model(data))
	conn.Send("RPUSH", fmt.Sprintf("%sfailed:%s", workerSettings.Namespace, data.Queue), buf)
	conn.Flush()
}

func failedData2Model(data *failureData) *failure {
	workerData := strings.Split(data.Worker, ":")
	woeker, _ := newWorker(strings.Join(workerData[:1], ":"), []string{workerData[2]})
	fai := &failure{
		FailedAt:  time.Now(),
		Payload:   data.Payload,
		Exception: "Error",
		Error:     data.Error,
		Worker:    woeker,
		MaxTry:    workerSettings.MaxRetry,
		Queue:     data.Queue,
		Times:     data.Times,
		ExecTime:  data.ExecTime,
	}
	return fai
}
