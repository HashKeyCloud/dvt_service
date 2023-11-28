package store

import (
	"context"
)

const SSVClusterSnapshot = "ssv:cluster:snapshot"

func (s *Store) CheckSSVClusterSnapshot(ctx context.Context, key string) bool {
	return s.redis.HExists(ctx, SSVClusterSnapshot, key).Val()
}

func (s *Store) GetSSVClusterSnapshot(ctx context.Context, key string) string {
	return s.redis.HGet(ctx, SSVClusterSnapshot, key).Val()
}

func (s *Store) SetSSVClusterSnapshot(ctx context.Context, key, cluster string) {
	s.redis.HSet(ctx, SSVClusterSnapshot, key, cluster)
}

func (s *Store) DelSSVClusterSnapshot(ctx context.Context, key string) {
	s.redis.HDel(ctx, SSVClusterSnapshot, key)
}
