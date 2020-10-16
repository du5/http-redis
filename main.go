package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
)

type TomlConfig struct {
	RedisDB RedisDB `toml:"redis"`
}

type RedisDB struct {
	Connection *redis.Client
	DBHost     string `toml:"host"`
	DBPort     int    `toml:"port"`
	DBPass     string `toml:"pass"`
	DBName     int    `toml:"name"`
}

type RedisStats struct {
	Stats   bool
	Message string
}

type RedisVal struct {
	Stats bool
	Val   string
}

func (rdb *RedisDB) open() error {
	connection := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdb.DBHost, rdb.DBPort),
		Password: rdb.DBPass,
		DB:       rdb.DBName,
	})
	rdb.Connection = connection

	return rdb.Connection.Ping().Err()
}

func (rdb *RedisDB) close() error {
	return rdb.Connection.Close()
}

func checkStatsHelper(err error, ifok string) (bool, string) {
	if err != nil {
		return false, err.Error()
	} else {
		return true, ifok
	}
}

func (rdb *RedisDB) ping() RedisStats {
	stats := RedisStats{}
	err := rdb.Connection.Ping().Err()
	stats.Stats, stats.Message = checkStatsHelper(err, "redis db is running.")
	return stats
}

func (rdb *RedisDB) get(c echo.Context) error {
	val := RedisVal{}
	redisVal, err := rdb.Connection.Get(c.Param("key")).Result()
	val.Stats, val.Val = checkStatsHelper(err, redisVal)
	jsonStr, _ := json.Marshal(val)
	return c.JSONBlob(http.StatusOK, jsonStr)
}

func (rdb *RedisDB) rdbping(c echo.Context) error {
	jsonStr, _ := json.Marshal(rdb.ping())
	return c.JSONBlob(http.StatusOK, jsonStr)
}

func main() {
	var conf TomlConfig
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}

	rdb := conf.RedisDB
	if err := rdb.open(); err != nil {
		panic(err)
	}
	defer rdb.close()
	e := echo.New()
	e.GET("/", rdb.rdbping)
	e.GET("/get/:key", rdb.get)
	e.Logger.Fatal(e.Start(":80"))
}
