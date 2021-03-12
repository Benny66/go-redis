package redis

/*
 * @Descripttion:
 * @version: v1.0.0
 * @Author: shahao
 * @Date: 2021-03-11 11:43:23
 * @LastEditors: shahao
 * @LastEditTime: 2021-03-12 15:44:40
 */

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

// 父级interface
type RedisFatherInterface interface {
	// 获取所有keys
	GetAllKeys() []string
}

// ListInterface 操作list接口
type ListInterface interface {
	// "继承"父类的所有方法
	RedisFatherInterface

	//队列
	PopNoWait(key string) (string, error) // ~di: interface{}类型、数据压缩
	Pop(key string, timeout int) (string, error)
	Put(key string, value string, timeout int) (int, error)
	PutNoWait(key string, value string) (int, error)
	QSize(key string) int
	Empty(key string) bool
	Full(key string) bool

	//保存字符串key:value
	Set(key string, value string, expire int) (bool, error)
	//获取字符串
	Get(key string) (string, error)
	//获取key的生存时间
	TTL(key string) (int64, error)
	//获取单个哈希值
	HGet(hashKey, key string) (string, error)
	//获取所有哈希值
	HGetAll(hashKey string) (map[string]string, error)
	//设置哈希值
	HSet(hashKey, key, value string) (bool, error)

	//添加集合中的一个值
	SAdd(setKey, value string) (int, error)
	//移除集合中的一个值
	SRem(setKey, value string) (int, error)
	//获取集合的所有成员
	SMembers(setKey string) ([]string, error)
}

// redis所有操作的接口
type RedisInterface interface {
	ListInterface
}

// 工厂函数，要求对应的结构体必须实现 RedisInterface 中的所有方法
// 如果只想实现某一些方法，就返回"有这些方法的结构体"就好了
func ProduceRedis(host, port, password string, db, maxSize int, lazyLimit bool) (RedisInterface, error) {

	// 要求RRedis结构体实现返回的接口中所有的方法！
	redisObj := &RRedis{
		maxIdle:        100,
		maxActive:      130,
		maxIdleTimeout: time.Duration(60) * time.Second,
		maxTimeout:     time.Duration(30) * time.Second,
		lazyLimit:      lazyLimit,
		maxSize:        maxSize,
	}
	// 建立连接池
	redisObj.redisCli = &redis.Pool{
		MaxIdle:     redisObj.maxIdle,
		MaxActive:   redisObj.maxActive,
		IdleTimeout: redisObj.maxIdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial(
				"tcp",
				host+":"+port, // address
				redis.DialPassword(password),
				redis.DialDatabase(int(db)),
				redis.DialConnectTimeout(redisObj.maxTimeout),
				redis.DialReadTimeout(redisObj.maxTimeout),
				redis.DialWriteTimeout(redisObj.maxTimeout),
			)
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}

	return redisObj, nil
}
