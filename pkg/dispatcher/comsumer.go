package dispatcher

import (
	"crypto/sha1"
	"fmt"
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
					if status[0] == "waiting" {
						data := make(map[string]interface{})
						data["status"] = "processing"
						data["process_time"] = time.Now().Unix()
						err = redisConn.HMSet(dcardKey, data).Err()
						if err != nil {
							logger.Error("Redis Error", zap.String("Job Worker", "hmset dcard key status"), zap.String("error", err.Error()))
						}

						fileContent, err := redisConn.HMGet(dcardKey, "fileContent").Result()
						if err != nil {
							logger.Error("Redis Error", zap.String("Job Worker", "hmget key fileContent"), zap.String("error", err.Error()))
						}

						// hash file content
						h := sha1.New()
						h.Write([]byte(fileContent[0].(string)))
						bs := h.Sum(nil)
						shaStr := fmt.Sprintf("%x", bs)

						sha := make(map[string]interface{})
						sha["sha"] = shaStr
						sha["status"] = "success"
						sha["process_time"] = ""
						err = redisConn.HMSet(dcardKey, sha).Err()
						if err != nil {
							logger.Error("Redis Error", zap.String("Job Worker", "hmset key sha and status"), zap.String("error", err.Error()))
						} else {
							// delete scan key
							_, err = redisConn.Del(key).Result()
							if err != nil {
								logger.Error("Redis Error", zap.String("Job Worker", "delete scan key"), zap.String("error", err.Error()))
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
