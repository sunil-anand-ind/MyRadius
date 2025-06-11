package datastore

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

var (
	// Redis Server Address
	RedisServerAddr = "localhost:6379"
	// Redis Record Persistent Duration (in hours)
	RecordPersistDuration int = 24
)

/* Accounting Details */
type AccountingRecord struct {
	UserName string `json: "username"`
	NASIPAddress string `json: "nas_ip-address"`
	NASPort string `json: "nas_port"`
	AcctStatusType string `json: "acct_status_type"`
	AcctSessionId string `json: "acct_session_id"`
	FramedIpAddress string `json: "framed_ip_address"`
	CallingStationId string `json: "calling_station_id"`
	CalledStationId string `json: "called_station_id"`
	Timestamp string `json: "timestamp"`
	ClientIp string `json: "client_ip"`
	PacketType string `json: "packet_type"`
}

/* Initialize the Redis Client */
func InitRedisClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr: RedisServerAddr,
		Password : "",
		DB: 0,
	})


	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unable to initialize connection towards Redis Server", err)
	}
	log.Printf("Successfully connected to Redis Server: %s",RedisServerAddr)
}

/* Store Accounting Records into redis by formulating the JSON */
func StoreAccountingRecord(key string, record AccountingRecord) {
	jsonRecord, err := json.Marshal(record)
	if err != nil {
		log.Fatal("Unable to prepare record", err)
	}

	err = rdb.Set(ctx, key, jsonRecord, 24*time.Hour).Err()
	if err != nil {
		log.Fatal("Unable to store record " ,err)
	}
}
