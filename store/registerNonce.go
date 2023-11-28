package store

import "context"

const SSVRegisterNonce = "ssv:register:nonce"

func (s *Store) CheckSSVRegisterNonce(ctx context.Context) bool {
	return s.redis.Exists(ctx, SSVRegisterNonce).Val() == 1
}

func (s *Store) IncrSSVRegisterNonce(ctx context.Context) {
	s.redis.Incr(ctx, SSVRegisterNonce)
}

func (s *Store) GetSSVRegisterNonce(ctx context.Context) string {
	return s.redis.Get(ctx, SSVRegisterNonce).Val()
}

func (s *Store) SetSSVRegisterNonce(ctx context.Context, nonce int) {
	s.redis.Set(ctx, SSVRegisterNonce, nonce, -1)
}
