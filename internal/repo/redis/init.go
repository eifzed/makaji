package redis

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/eifzed/joona/internal/config"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/prometheus/common/log"
)

var (
	ErrKeyNotFound = errors.New("Key Not Found")

	dial = redigo.Dial
)

// NetworkTCP for tcp
const NetworkTCP = "tcp"

// Options for redis
type Options struct {
	MaxIdle       int
	MaxActive     int
	TimeoutSecond int
	Wait          bool
	AuthKey       string
	Cfg           *config.Config
	Address       string
}

var authClientFunc = authorizeClient

// New redis connection
func New(opt Options) *Store {
	return &Store{
		Pool: &redigo.Pool{
			MaxIdle:     opt.MaxIdle,
			MaxActive:   opt.MaxActive,
			IdleTimeout: time.Duration(opt.TimeoutSecond) * time.Second,
			Dial: func() (redigo.Conn, error) {
				c, err := dial(NetworkTCP, opt.Address)
				if err != nil {
					log.Errorln("[Redis Pool]:", err.Error())
				}

				if len(opt.AuthKey) > 0 {
					if _, err = authClientFunc(c, opt.AuthKey); err != nil {
						log.Errorln("[Redis Auth]:", err.Error())
					}
				}

				return c, err
			},
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				// _, err := c.Do("PING")
				return nil
			},
		},
	}
}

func authorizeClient(c redigo.Conn, key string) (reply interface{}, err error) {
	return c.Do("AUTH", key)
}

// Store object
type Store struct {
	Pool  Handler
	mutex sync.Mutex
}

// Handler handler for cache
type Handler interface {
	Get() redigo.Conn
	GetContext(context.Context) (redigo.Conn, error)
}

func (r *Store) Ping() error {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}

// Get string value
func (r *Store) Get(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	resp, err := redigo.String(conn.Do("GET", key))
	if err == redigo.ErrNil {
		return "", ErrKeyNotFound
	}
	return resp, err
}

// Del key value
func (r *Store) Del(key string) (int64, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	resp, err := redigo.Int64(conn.Do("DEL", key))
	if err == redigo.ErrNil {
		return 0, ErrKeyNotFound
	}
	return resp, err
}

// HGetAll key and value
func (r *Store) HGetAll(key string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Strings(conn.Do("HGETALL", key))
}

// HGet key and value
func (r *Store) HGet(key, field string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.String(conn.Do("HGET", key, field))
}

// HSet set has map
func (r *Store) HSet(key, field string, value interface{}) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	resp, err := redigo.String(conn.Do("HSET", key, field, value))
	return resp, err
}

// SetWithExpire will be used to set the value with expire
func (r *Store) SetWithExpire(key string, value interface{}, expire int) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	resp, err := redigo.String(conn.Do("SET", key, value, "EX", expire))
	return resp, err
}

// SetWithoutExpire will be used to set the value without expire
func (r *Store) SetWithoutExpire(key string, value interface{}) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	resp, err := redigo.String(conn.Do("SET", key, value))
	return resp, err
}

// ExtendExpiration expiration redis
func (r *Store) ExtendExpiration(key string, expire int) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	resp, err := redigo.Int(conn.Do("EXPIRE", key, expire))
	return resp, err
}

// AddInSet will be used to add value in set
func (r *Store) AddInSet(key, value string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("SADD", key, value))
}

// GetSetMembers will be used to get the set memebers
func (r *Store) GetSetMembers(key string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Strings(conn.Do("SMEMBERS", key))
}

// GetSetLength will be used to get the set length
func (r *Store) GetSetLength(key string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("SCARD", key))
}

// GetNElementOfSet to get the first N elements of set
func (r *Store) GetNElementOfSet(key string, n int) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Strings(conn.Do("SPOP", key, n))
}

// PushNElementToSet will be used to push n elements to set
func (r *Store) PushNElementToSet(values []interface{}) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("SADD", values...))
}

// HMSet function
// please use basic types only (no struct, array, or map) for kv value
func (r *Store) HMSet(key string, kv map[string]interface{}) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	var (
		args = make([]interface{}, 1+(len(kv)*2))
		idx  = 1
	)
	args[0] = key
	for k, v := range kv {
		args[idx] = k
		args[idx+1] = v
		idx += 2
	}
	return redigo.String(conn.Do("HMSET", args...))
}

// HMGet keys and value
func (r *Store) HMGet(key string, fields ...string) ([]string, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	return redigo.Strings(conn.Do("HMGET", args...))
}

// HDel fields of a key
func (r *Store) HDel(key string, fields ...string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	args := make([]interface{}, len(fields)+1)
	args[0] = key
	for i, field := range fields {
		args[i+1] = field
	}
	return redigo.Int(conn.Do("HDEL", args...))
}

// Append string to existing value in the key
func (r *Store) Append(key, value interface{}) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("APPEND", key, value))
}

// Exist check if key exists in redis
func (r *Store) Exist(key string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("EXISTS", key))
}

// TimeToLive is used to check a key expiration time
func (r *Store) TimeToLive(key string) (int, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	return redigo.Int(conn.Do("TTL", key))
}
