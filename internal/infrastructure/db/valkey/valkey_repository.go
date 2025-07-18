package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

type ValkeyRepository struct {
	client valkey.Client
}

func NewValkeyRepository() (ValkeyRepository, error) {
	client, err := NewConnection()
	if err != nil {
		return ValkeyRepository{}, err
	}

	return ValkeyRepository{
		client: client,
	}, nil
}

func (r *ValkeyRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Do(ctx, r.client.B().Get().Key(key).Build()).ToString()
}

func (r *ValkeyRepository) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	val := toString(value)
	cmd := r.client.B().Set().Key(key).Value(val)
	if ttl > 0 {
		cmd.ExSeconds(int64(ttl))
	}
	return r.client.Do(ctx, cmd.Build()).Error()
}

func (r *ValkeyRepository) Delete(ctx context.Context, keys ...string) error {
	return r.client.Do(ctx, r.client.B().Del().Key(keys...).Build()).Error()
}

func (r *ValkeyRepository) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Do(ctx, r.client.B().Exists().Key(key).Build()).AsInt64()
	return count > 0, err
}

func (r *ValkeyRepository) Expire(ctx context.Context, key string, ttl int) error {
	return r.client.Do(ctx, r.client.B().Expire().Key(key).Seconds(int64(ttl)).Build()).Error()
}

func (r *ValkeyRepository) Increment(ctx context.Context, key string) (int64, error) {
	return r.client.Do(ctx, r.client.B().Incr().Key(key).Build()).AsInt64()
}

func (r *ValkeyRepository) HSet(ctx context.Context, key string, values map[string]interface{}) error {
	cmd := r.client.B().Hset().Key(key).FieldValue()
	for key, value := range values {
		cmd.FieldValue(key, toString(value))
	}
	return r.client.Do(ctx, cmd.Build()).Error()
}

func (r *ValkeyRepository) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.Do(ctx, r.client.B().Hget().Key(key).Field(field).Build()).ToString()
}

func (r *ValkeyRepository) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	cmd := r.client.B().Hgetall().Key(key).Build()
	res, err := r.client.Do(ctx, cmd).AsStrMap()
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
}

func (r *ValkeyRepository) LPush(ctx context.Context, key string, values ...interface{}) error {
	strValues := make([]string, 0, len(values))
	for _, v := range values {
		strValues = append(strValues, toString(v))
	}

	cmd := r.client.B().Lpush().Key(key).Element(strValues...).Build()
	return r.client.Do(ctx, cmd).Error()
}

func (r *ValkeyRepository) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	cmd := r.client.B().Lrange().Key(key).Start(start).Stop(stop).Build()
	return r.client.Do(ctx, cmd).AsStrSlice()
}

func (r *ValkeyRepository) Close() error { return nil }

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}
