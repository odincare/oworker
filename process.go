package oworker

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type process struct {
	Hostname string
	Pid      int
	ID       string
	Queues   []string
}

func newProcess(id string, queues []string) (*process, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &process{
		Hostname: hostname,
		Pid:      os.Getpid(),
		ID:       id,
		Queues:   queues,
	}, nil
}

func (p *process) String() string {
	return fmt.Sprintf("%s:%d-%s:%s", p.Hostname, p.Pid, p.ID, strings.Join(p.Queues, ","))
}

func (p *process) open(conn *RedisConn) error {
	deleteByKey(conn, fmt.Sprintf("%sworker:%s*", workerSettings.Namespace, p.Hostname))
	deleteByKey(conn, fmt.Sprintf("%sstat:processed:%s*", workerSettings.Namespace, p.Hostname))
	deleteByKey(conn, fmt.Sprintf("%sstat:failed:%s*", workerSettings.Namespace, p.Hostname))
	conn.Send("SADD", fmt.Sprintf("%sworkers", workerSettings.Namespace), p)
	conn.Send("SET", fmt.Sprintf("%sstat:processed:%v", workerSettings.Namespace, p), "0")
	conn.Send("SET", fmt.Sprintf("%sstat:failed:%v", workerSettings.Namespace, p), "0")
	conn.Flush()

	return nil
}

func (p *process) close(conn *RedisConn) error {
	logger.Infof("%v shutdown", p)
	conn.Send("SREM", fmt.Sprintf("%sworkers", workerSettings.Namespace), p)
	conn.Send("DEL", fmt.Sprintf("%sstat:processed:%s", workerSettings.Namespace, p))
	conn.Send("DEL", fmt.Sprintf("%sstat:failed:%s", workerSettings.Namespace, p))
	conn.Flush()

	return nil
}

func (p *process) start(conn *RedisConn) error {
	conn.Send("SET", fmt.Sprintf("%sworker:%s:started", workerSettings.Namespace, p), time.Now().String())
	conn.Flush()

	return nil
}

func (p *process) finish(conn *RedisConn) error {
	conn.Send("DEL", fmt.Sprintf("%sworker:%s", workerSettings.Namespace, p))
	conn.Send("DEL", fmt.Sprintf("%sworker:%s:started", workerSettings.Namespace, p))
	conn.Flush()

	return nil
}

func (p *process) fail(conn *RedisConn) error {
	conn.Send("INCR", fmt.Sprintf("%sstat:failed", workerSettings.Namespace))
	conn.Send("INCR", fmt.Sprintf("%sstat:failed:%s", workerSettings.Namespace, p))
	conn.Flush()

	return nil
}

func (p *process) queues(strict bool) []string {
	// If the queues order is strict then just return them.
	if strict {
		return p.Queues
	}

	// If not then we want to to shuffle the queues before returning them.
	queues := make([]string, len(p.Queues))
	for i, v := range rand.Perm(len(p.Queues)) {
		queues[i] = p.Queues[v]
	}
	return queues
}
