package db

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	mysqlCfg "github.com/zngue/zng_app/db/mysql"
	redisCfg "github.com/zngue/zng_app/db/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

func NewRedis(fns ...redisCfg.RedisFn) (*redis.Client, func(), error) {
	var config = &redisCfg.RedisOption{
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
	_, err := redisClient.Ping().Result()
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
		fmt.Println("redis close")
	}
	return redisClient, cleanup, nil
}
func NewDB(fns ...mysqlCfg.MysqlFn) (db *gorm.DB, err error) {
	var config = &mysqlCfg.MysqlOption{
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
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`,加S
			SingularTable: true,
		},
		//对于写操作（创建、更新、删除），为了确保数据的完整性，GORM 会将它们封装在事务内运行。但这会降低性能，你可以在初始化时禁用这种方式
		SkipDefaultTransaction: true,
	})
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
