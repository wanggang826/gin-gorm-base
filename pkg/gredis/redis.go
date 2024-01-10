package gredis

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"

	"gin-gorm-base/pkg/logging"
	"gin-gorm-base/pkg/setting"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			//redis.DialDatabase(setting.RedisSetting.Db)
			if err != nil {
				logging.Error("Redis Dial error:" + err.Error())
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					c.Close()
					logging.Error("Redis AUTH error:" + err.Error())
					return nil, err
				}
			}
			if setting.RedisSetting.Db > 0 {
				if _, err := c.Do("SELECT", setting.RedisSetting.Db); err != nil {
					c.Close()
					logging.Error("Redis SELECT:" + err.Error())
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func SetInt(key string, data int, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetString(key string, data string, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		return err
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// Get get a key
func GetInt(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

// Get get a key
func GetString(key string) (string, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return reply, nil
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func HmSet(key string, time int, data ...interface{}) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("HMSet", key, data)
	if err != nil {
		return err
	}
	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

func HGetInt(key, valueKey string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	val, err := redis.Int(conn.Do("HGet", key, valueKey))
	if err != nil {
		return 0, err
	}
	return val, nil
}

func HGetSting(key string, valueKey string) (string, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("HGet", key, valueKey))
	if err != nil {
		return "", err
	}
	return reply, nil
}

func HGetAll(key string) ([]interface{}, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	values, err := redis.Values(conn.Do("HGetAll", key))

	if err != nil {
		return nil, err
	}
	return values, nil
}

func LPush(key string, time int, data ...interface{}) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("LPush", key, data)
	if err != nil {
		return err
	}
	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

func RPush(key string, time int, data ...interface{}) error {
	conn := RedisConn.Get()
	defer conn.Close()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = conn.Do("RPush", key, value)
	if err != nil {
		return err
	}
	if time > 0 {
		_, err = conn.Do("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

func LPop(key string) (string, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.String(conn.Do("LPop", key))
	if err != nil {
		return "", err
	}
	return reply, nil
}

func LLen(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("LLen", key))
	if err != nil {
		return 0, err
	}
	return reply, nil
}

func LRange(key string) ([]string, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	r, err := redis.Strings(conn.Do("LRange", key, 0, -1))
	if err != nil {
		return nil, err
	}
	return r, nil
}
