package db

type RedisInterface interface {
	Get(key string) (string, error)
	Del(key string) (int64, error)
	HGetAll(key string) ([]string, error)
	HGet(key, field string) (string, error)
	HSet(key, field string, value interface{}) (string, error)
	SetWithExpire(key string, value interface{}, expire int) (string, error)
	SetWithoutExpire(key string, value interface{}) (string, error)
	ExtendExpiration(key string, expire int) (int, error)
	AddInSet(key, value string) (int, error)
	GetSetMembers(key string) ([]string, error)
	GetSetLength(key string) (int, error)
	GetNElementOfSet(key string, n int) ([]string, error)
	PushNElementToSet(values []interface{}) (int, error)
	HMSet(key string, kv map[string]interface{}) (string, error)
	HMGet(key string, fields ...string) ([]string, error)
	HDel(key string, fields ...string) (int, error)
	Append(key, value interface{}) (int, error)
	Exist(key string) (int, error)
	TimeToLive(key string) (int, error)
	Ping() error
}
