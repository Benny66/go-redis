package main

/*
 * @Descripttion:
 * @version: v1.0.0
 * @Author: shahao
 * @Date: 2021-03-11 11:43:57
 * @LastEditors: shahao
 * @LastEditTime: 2021-03-12 15:50:43
 */

import (
	"encoding/json"
	"fmt"
	"time"

	red "goRedis/redis"
	//"github.com/garyburd/redigo/redis"
)

func main() {
	// 测试redis的操作
	testRedisOperate()
	//testRedis()
}

type Doctor struct {
	ID      int64
	Name    string
	Age     int
	Sex     int
	AddTime time.Time
}

// 测试redis的操作
func testRedisOperate() {
	redisCLi, err := red.ProduceRedis("127.0.0.1", "6379", "", 0, 100, true)
	if err != nil {
		fmt.Println("redis连接错误！err>>>", err.Error())
		return
	}

	//保存key，字符串
	res, err := redisCLi.Set("key", "value", 10)
	if err != nil {
		fmt.Println("dataStr", err)
		return
	}
	fmt.Println("set string:", res)
	// data, err := redisCLi.GetNoWait("key")
	// fmt.Println("value>>> ", data)

	// 查询key
	keyStr, err2 := redisCLi.Get("key")
	if err != nil {
		fmt.Println("get dataStr:", err2)
		return
	}
	fmt.Println("get string", keyStr)
	//保存json字符串
	doctor := Doctor{1, "钟南山", 83, 1, time.Now()}
	doctorJson, _ := json.Marshal(doctor)

	res, err = redisCLi.Set("doctor2", string(doctorJson), 1000)
	if err != nil {
		fmt.Println("dataStr", err)
		return
	}
	fmt.Println("set json string:", res)

	expire, err := redisCLi.TTL("doctor2")
	if err != nil {
		fmt.Println("expire", err)
		return
	}
	fmt.Println("过期时间：", expire)

	result3, err := redisCLi.Get("doctor2")
	if err != nil {
		fmt.Println("doctor2 err:", err)
		return
	}
	fmt.Println("doctor2:", result3)
	bytes := []byte(result3)
	testT := Doctor{}
	err = json.Unmarshal(bytes, &testT)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println(testT)

	res, err = redisCLi.HSet("users", "name", "benny")
	if err != nil {
		fmt.Println("1err:", err)
		return
	}
	fmt.Println(res)
	res, err = redisCLi.HSet("users", "password", "123456")
	if err != nil {
		fmt.Println("2err:", err)
		return
	}
	fmt.Println(res)
	name, err := redisCLi.HGet("users", "name")
	if err != nil {
		fmt.Println("3err:", err)
		return
	}
	fmt.Println(name)
	hash, err := redisCLi.HGetAll("users")
	if err != nil {
		fmt.Println("4err:", err)
		return
	}
	fmt.Println(hash)

	sk, err := redisCLi.SAdd("setKey", "123")
	if err != nil {
		fmt.Println("5err:", err)
		return
	}
	sk, err = redisCLi.SAdd("setKey", "456")
	if err != nil {
		fmt.Println("6err:", err)
		return
	}
	fmt.Println(sk)

	skk, err := redisCLi.SMembers("setKey")
	if err != nil {
		fmt.Println("7err:", err)
		return
	}
	fmt.Println(skk)

	sk2, err := redisCLi.SRem("setKey", "123")
	if err != nil {
		fmt.Println("8err:", err)
		return
	}
	fmt.Println(sk2)
	skk, err = redisCLi.SMembers("setKey")
	if err != nil {
		fmt.Println("8err:", err)
		return
	}
	fmt.Println(skk)

	allKeysLst := redisCLi.GetAllKeys()

	fmt.Println("key>>> ", allKeysLst)

}

func testRedis() {

}
