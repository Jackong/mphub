package redis

import (
	"os"

	log "github.com/Sirupsen/logrus"

	r "gopkg.in/redis.v3"
)

//Client for redis
var Client *r.Client

func init() {
	addr := os.Getenv("REDIS_ADDR")
	Client = r.NewClient(&r.Options{
		Addr: addr,
		DB:   0,
	})
	log.WithField("addr", addr).Info("Starting redis")
	pong, err := Client.Ping().Result()
	if err != nil || pong != "PONG" {
		panic(err)
	}
}
