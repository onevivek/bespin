package database

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/onevivek/bespin/go_lang/gmail/account"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(address string) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := rdb.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}

	return &Redis{client: rdb}, nil
}

func (r *Redis) Get(ctx context.Context, id string) (account.Account, error) {
	resp, err := r.client.HGetAll(ctx, r.key(id)).Result()
	if err != nil {
		return account.Account{}, fmt.Errorf("error getting account %s from redis: %w", id, err)
	}

	return account.ParseMap(resp)
}

func (r *Redis) Save(ctx context.Context, id string, a account.Account) error {
	value := account.ParseAccount(a)
	if err := r.client.HSet(ctx, r.key(id), value).Err(); err != nil {
		return fmt.Errorf("error saving account %s to redis: %w", id, err)
	}

	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) key(accountId string) string {
	return fmt.Sprintf("account:%s", accountId)
}
