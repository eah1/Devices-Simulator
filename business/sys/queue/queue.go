// Package queue provides support for access the queue.
package queue

import (
	"device-simulator/app/config"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// NewQueue create a config queue.
func NewQueue(config config.Config, log *zap.SugaredLogger) *asynq.Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.QueueHost + ":" + config.QueuePort},
		asynq.Config{
			IsFailure: func(err error) bool {
				log.Errorf("asynq server exec task IsFailure ======== >>>>>>>>>>>  err : %+v \n", err)

				return true
			},
			Concurrency: config.QueueConcurrency,
			Queues: map[string]int{
				"emails":  1,
				"metrics": 9,
			},
			Logger: log,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Second * 10
			},
		},
	)

	return srv
}
