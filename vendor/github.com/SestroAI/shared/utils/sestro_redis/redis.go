package sestro_redis

import (
	"fmt"
	"github.com/SestroAI/shared/logger"
	"github.com/garyburd/redigo/redis"
	"os"
)

var redisPool *redis.Pool

func init() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if redisAddr == "" || redisPassword == "" {
		logger.Errorf("No REDIS_ADDR env variable found. Please ignore if you are not using redis")
		os.Exit(-1)
	}

	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", redisAddr)
			if redisPassword == "" {
				return conn, err
			}
			if err != nil {
				return nil, err
			}
			if _, err := conn.Do("AUTH", redisPassword); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		Wait:      false,
		MaxActive: 30, //30 is the limit on free redis cache instance
	}
}

type RedisConn struct {
	Conn redis.Conn
}

func GetNewRedisConnection() RedisConn {
	redisConn := redisPool.Get()
	return RedisConn{redisConn}
}

func (rconn *RedisConn) SaveKeyValueInRedis(key string, value interface{}) error {
	_, err := rconn.Conn.Do("SET", key, value)
	if err != nil {
		logger.Errorf("Unable to set Redis value = %s for key = %s", value, key)
		return err
	}
	return nil
}

func (rconn *RedisConn) GetKeyValueFromRedis(key string) (interface{}, error) {
	value, err := rconn.Conn.Do("GET", key)
	if err != nil {
		logger.Errorf("Unable to get Redis value key = %s", key)
		return nil, err
	}
	return value, nil
}

func main() {
	fmt.Println("Hello, playground")
	sr_conn := GetNewRedisConnection()

	redis_key := "341414" + "::" + "32423"
	value, err := sr_conn.GetKeyValueFromRedis(redis_key)
	fmt.Print(value, err)
}
