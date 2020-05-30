package models

import (
	"fmt"
	"go.uber.org/zap"
	"time"

	"github.com/spf13/viper"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *Databases

type Databases struct {
	MysqlDb *sqlx.DB
	RedisDB *redis.Client
}

// Init 初始化MySQL连接
func (db *Databases) InitMySql() (err error) {
	// "user:password@tcp(host:port)/dbname"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.db"),
	)
	db.MysqlDb, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return
	}
	//连通性测试
	err = db.MysqlDb.Ping()
	// 设置最大连接数
	db.MysqlDb.SetMaxOpenConns(viper.GetInt("mysql.maxconns"))
	// 设置最大空闲连接数
	db.MysqlDb.SetMaxIdleConns(viper.GetInt("mysql.idleconns"))
	return
}

// Close 程序退出时释放MySQL连接
// 不直接对外暴露db变量，而是对外暴露一个Close方法
func (db *Databases) CloseMysql() {
	err := db.MysqlDb.Close()
	zap.L().Error("Mysql数据库关闭失败", zap.Error(err))
}

//初始化redis
func (db *Databases) InitRedis() (err error) {
	db.RedisDB = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		DB:          viper.GetInt("redis.db"), // use default DB
		DialTimeout: viper.GetDuration("redis.timeout") * time.Second,
		//ReadTimeout: viper.GetDuration("redis.timeout")*time.Second,
		//WriteTimeout: viper.GetDuration("redis.timeout")*time.Second,
		//PoolTimeout: viper.GetDuration("redis.timeout")*time.Second,
	})
	_, err = db.RedisDB.Ping().Result()
	return
}

func (db *Databases) CloseRedis() {
	err := db.RedisDB.Close()
	zap.L().Error("Redis关闭失败", zap.Error(err))
}
