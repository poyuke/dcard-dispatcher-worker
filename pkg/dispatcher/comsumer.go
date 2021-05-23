package dispatcher

import (
	"fmt"
	"crypto/sha1"

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
		keys, err := redisConn.Keys("*").Result()
		if err != nil {
			logger.Error("Redis Error", zap.String("Job Worker", "get all key"), zap.String("error", err.Error()))
		}
		for _, key := range keys {
			status, err := redisConn.HMGet(key, "status").Result()
			if err != nil {
				logger.Error("Redis Error", zap.String("Job Worker", "hmget key status"), zap.String("error", err.Error()))
			}
			if status[0] == "waiting" {
				data := make(map[string]interface{})
				data["status"] = "processing"
				err := redisConn.HMSet(key, data).Err()
				if err != nil {
					logger.Error("Redis Error", zap.String("Job Worker", "hmset key status"), zap.String("error", err.Error()))
				}

				fileContent, err := redisConn.HMGet(key, "fileContent").Result()
				if err != nil {
					logger.Error("Redis Error", zap.String("Job Worker", "hmget key fileContent"), zap.String("error", err.Error()))
				}
				h := sha1.New()
				h.Write([]byte(fileContent[0].(string)))
    			bs := h.Sum(nil)
				shaStr := fmt.Sprintf("%x", bs)

				sha := make(map[string]interface{})
				sha["sha"] = shaStr
				err = redisConn.HMSet(key, sha).Err()
				if err != nil {
					logger.Error("Redis Error", zap.String("Job Worker", "hmset key sha"), zap.String("error", err.Error()))
				}

				data["status"] = "success"
				err = redisConn.HMSet(key, data).Err()
				if err != nil {
					logger.Error("Redis Error", zap.String("Job Worker", "hmset key status"), zap.String("error", err.Error()))
				}
			}
		}
	}()
	fmt.Printf(" [*] Waiting for upload file. To exit press CTRL+C")
	fmt.Println()
	<-forever
}
