package init

import (
	"context"
	"fmt"
	"time"

	"demo/hd/constant"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	RedisCli *redis.Client
	GormDB   *gorm.DB
)

func init() {
	initRedisClient()
	initDBClient()
}

func initRedisClient() {
	RedisCli = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", constant.RedisHost, constant.RedisPort),
		DialTimeout:  time.Second * 2,
		ReadTimeout:  time.Second * 2,
		WriteTimeout: time.Second * 2,
		PoolSize:     10,
		MinIdleConns: 10,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := RedisCli.Ping(ctx).Err(); err != nil {
		panic(err)
	}
}

func initDBClient() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", constant.MysqlUser, constant.MysqlPasswd, constant.MysqlHost, constant.MysqlPort, constant.MysqlDB)
	GormDB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 使用单数表名
			},
		},
	)

	if err != nil {
		panic(err)
	}

	sdb, err := GormDB.DB()
	if err != nil {
		panic(err)
	}
	sdb.SetMaxIdleConns(10)                  // 最大空闲连接数
	sdb.SetMaxOpenConns(100)                 // 最大连接数
	sdb.SetConnMaxLifetime(time.Minute * 10) // 设置连接空闲超时
}
