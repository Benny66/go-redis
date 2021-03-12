package redis

import (
	"errors"
	"math"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RRedis struct {
	redisCli       *redis.Pool
	maxIdle        int
	maxActive      int
	maxIdleTimeout time.Duration
	maxTimeout     time.Duration
	lazyLimit      bool
	maxSize        int
}

// 从池里获取连接 ———— 单独的方法
func (r *RRedis) getRedisConn() redis.Conn {
	rc := r.redisCli.Get()
	// // 用完后将连接放回连接池
	// defer rc.Close()
	return rc
}

// GetAllKeys 获取所有keys
func (r *RRedis) GetAllKeys() []string {

	rc := r.getRedisConn()
	// 用完后将连接放回连接池
	defer rc.Close()

	keys, err := redis.Strings(rc.Do("KEYS", "*"))
	if err != nil {
		return make([]string, 0)
	}
	return keys
}

func (r *RRedis) Pop(key string, timeout int) (string, error) {

	start := time.Now()

	for {
		res, err := r.PopNoWait(key)
		if err != nil {
			return "", err
		} else if res == "" {
			if timeout != -1 {
				lasted := time.Now().Sub(start)
				if r.maxTimeout > lasted {
					t1 := r.maxTimeout
					t2 := time.Duration(timeout)*time.Second - lasted
					time.Sleep(time.Duration(math.Min(float64(t1), float64(t2))))
				} else {
					return "", errors.New("GET timeout")
				}
			} else {
				time.Sleep(r.maxTimeout)
			}
		} else {
			return res, nil
		}
	}
}

func (r *RRedis) PopNoWait(key string) (string, error) {

	rc := r.getRedisConn()
	defer rc.Close()

	res, err := redis.String(rc.Do("LPOP", key))

	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *RRedis) Put(key string, value string, timeout int) (int, error) {

	start := time.Now()

	for {
		res, err := r.PutNoWait(key, value)

		if err != nil {
			return 0, err
		} else if res == -1 {
			if timeout != -1 {
				lasted := time.Now().Sub(start)
				if r.maxTimeout > lasted {
					t1 := r.maxTimeout
					t2 := time.Duration(timeout)*time.Second - lasted
					time.Sleep(time.Duration(math.Min(float64(t1), float64(t2))))
				} else {
					return 0, errors.New("PUT timeout")
				}
			} else {
				time.Sleep(r.maxTimeout)
			}

		} else {
			return res, nil
		}

	}
}

func (r *RRedis) PutNoWait(key string, value string) (int, error) {

	rc := r.getRedisConn()
	defer rc.Close()

	if r.Full(key) {
		return -1, nil
	}

	res, err := redis.Int(rc.Do("RPUSH", key, value))
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *RRedis) QSize(key string) int {

	rc := r.getRedisConn()
	defer rc.Close()

	res, err := redis.Int(rc.Do("LLEN", key))
	if err != nil {
		return -1
	}

	return res
}

func (r *RRedis) Empty(key string) bool {

	rc := r.getRedisConn()
	defer rc.Close()

	res, err := redis.Int(rc.Do("LLEN", key))
	if err != nil {
		return false
	}
	if res == 0 {
		return true
	}
	return false
}

func (r *RRedis) Full(key string) bool {

	if r.maxSize != 0 && r.QSize(key) >= r.maxSize {
		return true
	}

	return false
}

func (r *RRedis) Get(key string) (string, error) {
	rc := r.getRedisConn()
	defer rc.Close()
	res, err := redis.String(rc.Do("GET", key))
	if err != nil {
		return "", err
	}
	return res, nil
}

// 保存字符串key：value， 第三个参数为保存时间（s）24 * 3600则保存一天,-1为无过期时间
func (r *RRedis) Set(key string, value string, expire int) (bool, error) {

	rc := r.getRedisConn()
	defer rc.Close()

	_, err := redis.String(rc.Do("SET", key, value))
	if err != nil {
		return false, err
	}
	if expire != -1 {
		_, err := redis.Int64(rc.Do("EXPIRE", key, expire))
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

//获取key的过期时间
func (r *RRedis) TTL(key string) (int64, error) {
	rc := r.getRedisConn()
	defer rc.Close()

	res, err := redis.Int64(rc.Do("TTL", key))
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *RRedis) HGet(hashKey, key string) (string, error) {
	rc := r.getRedisConn()
	defer rc.Close()
	res, err := redis.String(rc.Do("HGET", hashKey, key))
	if err != nil {
		return "", err
	}
	return res, nil
}

func (r *RRedis) HGetAll(hashKey string) (map[string]string, error) {
	rc := r.getRedisConn()
	defer rc.Close()
	res, err := redis.StringMap(rc.Do("HGETALL", hashKey))
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 保存哈希值key：value
func (r *RRedis) HSet(hashKey, key, value string) (bool, error) {

	rc := r.getRedisConn()
	defer rc.Close()

	_, err := redis.Int(rc.Do("HSET", hashKey, key, value))
	if err != nil {
		return false, err
	}
	return true, nil
}

//保存集合
func (r *RRedis) SAdd(setKey, value string) (int, error) {

	rc := r.getRedisConn()
	defer rc.Close()

	if r.Full(setKey) {
		return -1, nil
	}
	res, err := redis.Int(rc.Do("SADD", setKey, value))
	if err != nil {
		return 0, err
	}
	return res, nil
}

//获取集合列表
func (r *RRedis) SMembers(hashKey string) ([]string, error) {
	rc := r.getRedisConn()
	defer rc.Close()
	res, err := redis.Strings(rc.Do("SMEMBERS", hashKey))
	if err != nil {
		return []string{}, err
	}
	return res, nil
}

//移除集合中一个值
func (r *RRedis) SRem(setKey, value string) (int, error) {

	rc := r.getRedisConn()
	defer rc.Close()
	res, err := redis.Int(rc.Do("SREM", setKey, value))
	if err != nil {
		return 0, err
	}
	return res, nil
}
