package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	mysqlConfig "github.com/zngue/zng_app/db/mysql"
	redisConfig "github.com/zngue/zng_app/db/redis"
	"github.com/zngue/zng_app/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func NewRedis(fns ...redisConfig.Fn) (*redis.Client, func(), error) {
	var config = &redisConfig.Option{
		Password: "",
		Port:     6379,
		Database: 0,
	}
	for _, fn := range fns {
		fn(config)
	}
	if config.Host == "" {
		return nil, nil, fmt.Errorf("redis host is empty")
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     30,
		PoolTimeout:  30 * time.Second,
		MinIdleConns: 10,
		DB:           config.Database,
	})
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		defer func(redisClient *redis.Client) {
			redisErr := redisClient.Close()
			if redisErr != nil {
				fmt.Println(redisErr)
				return
			}
		}(redisClient)
	}
	return redisClient, cleanup, nil
}
func NewDB(fns ...mysqlConfig.Fn) (db *gorm.DB, err error) {
	var config = &mysqlConfig.Option{
		Port: 3306,
	}
	for _, fn := range fns {
		fn(config)
	}
	if config.Username == "" {
		err = fmt.Errorf("mysql username is empty")
		return
	}
	if config.Password == "" {
		err = fmt.Errorf("mysql password is empty")
		return
	}
	if config.Host == "" {
		err = fmt.Errorf("mysql host is empty")
		return
	}
	if config.Database == "" {
		err = fmt.Errorf("mysql database name is empty")
		return
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		"Asia%2FShanghai",
	)
	var (
		sqlDB *sql.DB
	)
	var dbConfig = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   config.TablePrefix,
		},
		SkipDefaultTransaction: true,
	}
	if config.Logger {
		dbConfig.Logger = log.NewLog(config.LoggerConfig)
	}
	db, err = gorm.Open(mysql.Open(dsn), dbConfig)
	if err != nil {
		return
	}
	sqlDB, err = db.DB()
	if err != nil {
		return
	}
	//设置连接池
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(20)
	//可以服用的最大时间
	sqlDB.SetConnMaxLifetime(1700 * time.Second)
	return
}
