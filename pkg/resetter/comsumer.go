package resetter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"dispatcher-worker/pkg/log"

	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Component that include initial object
type Component struct {
	Logger *zap.Logger
}

// Execute is a rabbitmq comsumer
func Execute(cmd *cobra.Command, args []string) {
	var component Component

	// Init LOG
	logger, err := log.InitLog()
	if err != nil {
		logger.Fatal("Init log system fail", zap.String("message", err.Error()))
	}
	component.Logger = logger
	defer logger.Sync() // flushes buffer, if any

	redisConn := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: "",
		DB:       0,
	})
	err = redisConn.Ping().Err()
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.String("error", err.Error()))
	}

	forever := make(chan bool)
	go func() {
		for {
			var cursor uint64
			for {
				var keys []string
				var err error
				keys, cursor, err = redisConn.Scan(cursor, "scan-*", 20).Result()
				if err != nil {
					logger.Error("Redis Error", zap.String("Job Worker", "get scan key"), zap.String("error", err.Error()))
				}
				for _, key := range keys {
					dcardKey := strings.Replace(key, "scan", "dcard", 1)
					status, err := redisConn.HMGet(dcardKey, "status").Result()
					if err != nil {
						logger.Error("Redis Error", zap.String("Job Worker", "hmget dcard key status"), zap.String("error", err.Error()))
					}
					if status[0] == "processing" {
						processTime, err := redisConn.HMGet(dcardKey, "process_time").Result()
						if err != nil {
							logger.Error("Redis Error", zap.String("Job Worker", "hmget dcard key process_time"), zap.String("error", err.Error()))
						}

						pTime, _ := strconv.ParseInt(fmt.Sprintf("%v", processTime[0]), 10, 64)
						if time.Now().Unix()-pTime >= 60 {
							data := make(map[string]interface{})
							data["status"] = "waiting"
							data["sha"] = ""
							data["process_time"] = ""
							err = redisConn.HMSet(dcardKey, data).Err()
							if err != nil {
								logger.Error("Redis Error", zap.String("Job Worker", "hmset dcard key status"), zap.String("error", err.Error()))
							}
						}
					}
				}
				if cursor == 0 {
					break
				}
			}
		}
	}()
	fmt.Printf(" [*] To exit press CTRL+C")
	fmt.Println()
	<-forever
}
