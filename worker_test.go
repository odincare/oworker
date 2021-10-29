package oworker

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

var workerMarshalJSONTests = []struct {
	w        worker
	expected []byte
}{
	{
		worker{},
		[]byte(`":0-:"`),
	},
	{
		worker{
			process: process{
				Hostname: "hostname",
				Pid:      12345,
				ID:       "123",
				Queues:   []string{"high", "low"},
			},
		},
		[]byte(`"hostname:12345-123:high,low"`),
	},
}

func TestWorkerMarshalJSON(t *testing.T) {
	for _, tt := range workerMarshalJSONTests {
		actual, err := tt.w.MarshalJSON()
		if err != nil {
			t.Errorf("Worker(%#v): error %s", tt.w, err)
		} else {
			if string(actual) != string(tt.expected) {
				t.Errorf("Worker(%#v): expected %s, actual %s", tt.w, tt.expected, actual)
			}
		}
	}
}

func TestEnqueue(t *testing.T) {
	initConfig()
	expectedArgs := []interface{}{"a1", "lot", "of", "params"}
	jobName := "SomethingCool"
	queueName := "testQueue"
	expectedJob := &Job{
		Queue: queueName,
		Payload: Payload{
			Class: jobName,
			Args:  expectedArgs,
		},
	}

	workerSettings.Queues = []string{queueName}
	workerSettings.UseNumber = true
	workerSettings.ExitOnComplete = true

	err := Enqueue(expectedJob)
	if err != nil {
		t.Errorf("Error while enqueue %s", err)
	}
}

func TestWorker(t *testing.T) {
	initConfig()
	jobName := "SomethingCool"
	actualArgs := []interface{}{}
	actualQueueName := ""
	Register(jobName, func(queue string, args ...interface{}) error {
		actualArgs = args
		actualQueueName = queue
		fmt.Println(actualArgs)
		fmt.Println(actualQueueName)
		return errors.New("生成的错误")
	})

	RegisterFail(jobName, func(s string, data FailPayload) error {
		//fmt.Println("@@@@@@@@@处理错误.........@@@@@@")
		fmt.Println("错误信息", data)
		return errors.New("是大福利科技")
	})
	go func() {
		if err := Work(); err != nil {
			t.Errorf("(Enqueue) Failed on work %s", err)
		}
	}()

	time.Sleep(500 * time.Second)
}
func initConfig() {
	fmt.Println("测试输出---配置初始化....")
	redisHost := os.Getenv("REDIS_ADDR")
	pwd := os.Getenv("REDIS_PWD")
	settings := WorkerSettings{
		URI:            "redis://" + pwd + "*@" + redisHost + ":6379/10",
		Connections:    100,
		Queues:         []string{"testQueue", "delimited", "queues"},
		UseNumber:      true,
		ExitOnComplete: true,
		Concurrency:    2,
		Namespace:      "resque:",
		Interval:       5.0,
	}
	SetSettings(settings)
}
