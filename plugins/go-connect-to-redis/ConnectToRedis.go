package main

import (
	"context"
	"net/http"
	"time"

	"github.com/TykTechnologies/tyk/config"
	logger "github.com/TykTechnologies/tyk/log"
	"github.com/TykTechnologies/tyk/storage"
)

const pluginDefaultKeyPrefix = "Plugin-data:"

// Global redis variables
var conf config.Config
var rc *storage.RedisController

var store = storage.RedisCluster{KeyPrefix: pluginDefaultKeyPrefix}
var log = logger.Get()

func tykStoreData(key, value string) {
	ttl := int64(1000)
	store.SetKey(key, value, ttl)
}

func tykGetData(key string) (string, error) {
	val, err := store.GetKey(key)
	return val, err
}

func establishRedisConnection() {
	// Retrieve global configs
	conf = config.Global()

	// Create a Redis Controller, which will handle the Redis connection for the storage
	rc = storage.NewRedisController(context.Background())

	// Create a storage object, which will handle Redis operations using "apikey-" key prefix
	store = storage.RedisCluster{KeyPrefix: pluginDefaultKeyPrefix, HashKeys: conf.HashKeys, RedisController: rc}

	// Perform Redis connection
	go rc.ConnectToRedis(context.Background(), nil, &conf)
	for i := 0; i < 5; i++ { // max 5 attempts - should only take 2
		time.Sleep(5 * time.Millisecond)
		if rc.Connected() {
			log.Info("Redis Controller connected")
			break
		}
		log.Warn("Redis Controller not connected, will retry")
	}

	// Error handling Redis connection
	if !rc.Connected() {
		log.Error("Could not connect to storage")
		panic("Couldn't esetablished a connection to redis")
	}
}

func init() {
	establishRedisConnection()
}

// CallRedis to poke stuff in there for funzies
func CallRedis(rw http.ResponseWriter, r *http.Request) {
	var redisKey = "I"
	var redisValue = "Woz Here"
	log.Info("Start CallRedis Plugin")

	tykStoreData(redisKey, redisValue)
	log.Info("CallRedis: ", redisKey, " saved to redis")

	val, err := tykGetData(redisKey)
	if err != nil {
		log.Error("Plugin redis error ", err)
	}
	log.Info("CallRedis: value retrieved: " + val)

	log.Info("End CallRedis Plugin")

}

func main() {}
