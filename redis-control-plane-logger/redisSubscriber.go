package main

import (
	"context"
	"log"
	"os"
	"strings"
	"github.com/go-redis/redis/v8"
)

var (
    // Redis Server Address
    RedisServerAddr = "redis:6379"
	subctx = context.Background()
	subrdb *redis.Client
	keyPattern = "radius:acct:"
	radiusLog =  "/var/log/radius_updates.log"
)

/* Main entry point of the Redis Control plane logger program */
func main() {

	go InitSubscriber()

	select{}
}

/* Initialize the Redis Subscriber */
func InitSubscriber() {
	subrdb = redis.NewClient(&redis.Options{
		Addr: RedisServerAddr,
		Password : "",
		DB: 0,
	})


	_, err := subrdb.Ping(subctx).Result()
	if err != nil {
		log.Fatal("Unable to initialize connection towards Redis Server", err)
	}

	rdbSubscriber := subrdb.Subscribe(subctx, "__keyevent@0__:set")
	defer rdbSubscriber.Close()

	log.Printf("Successfully subscribed to Redis Server: %s",RedisServerAddr)

	// Creates a file and prints in a specific timestamp format.
	logFile, err := os.OpenFile(radiusLog, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("ERROR: Unable to create logfile %s", radiusLog)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)


	ch := rdbSubscriber.Channel()

	/* Filter the keys based on the pattern provided */
	for msg := range ch {
		if strings.HasPrefix(msg.Payload, keyPattern) {
			log.Printf("Received Update for key : %s", msg.Payload)
		}

	}
}
