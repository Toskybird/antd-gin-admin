package scanjob

import (
	"context"

	repointerfaces "antd-gin-admin-backend/internal/repository/interfaces"
	redisclient "antd-gin-admin-backend/pkg/redis"
)

const defaultQueueKey = "scan:jobs:queue"

// RedisQueue enqueues scan job codes into a Redis list.
type RedisQueue struct {
	client *redisclient.Client
	key    string
}

func NewRedisQueue(client *redisclient.Client) repointerfaces.ScanJobQueue {
	return &RedisQueue{client: client, key: defaultQueueKey}
}

func (q *RedisQueue) Enqueue(ctx context.Context, jobCode string) error {
	if q == nil || q.client == nil {
		return nil
	}
	return q.client.LPush(ctx, q.key, jobCode)
}

func (q *RedisQueue) ListEnqueued(ctx context.Context) ([]string, error) {
	if q == nil || q.client == nil {
		return []string{}, nil
	}
	return q.client.LRange(ctx, q.key, 0, -1)
}
