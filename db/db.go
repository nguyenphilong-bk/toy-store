package db

import (
	"fmt"
	"log"
	"os"

	_redis "github.com/go-redis/redis/v7"
	_ "github.com/lib/pq" //import postgres
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// Init ...
func Init() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=db port=5432 sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	var err error
	db, err = ConnectDB(dbinfo)
	if err != nil {
		log.Fatal(err)
	}

}

// ConnectDB ...
func ConnectDB(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dataSourceName, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db = db.Set("gorm:auto_preload", true)
	// db.Logger()
	sqldb, _ := db.DB()
	if err := sqldb.Ping(); err != nil {
		// mini.Log().Error().Err(err).Msg("ping failed to database")
		// mini.Log().Debug().Msg("reconnecting to database...")
		sqldb.Close()
		db = nil
		return ConnectDB(dataSourceName)
	}

	return db, nil
}

// GetDB ...
func GetDB() *gorm.DB {
	return db
}

// RedisClient ...
var RedisClient *_redis.Client

// InitRedis ...
func InitRedis(selectDB ...int) {

	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	RedisClient = _redis.NewClient(&_redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       selectDB[0],
		// DialTimeout:        10 * time.Second,
		// ReadTimeout:        30 * time.Second,
		// WriteTimeout:       30 * time.Second,
		// PoolSize:           10,
		// PoolTimeout:        30 * time.Second,
		// IdleTimeout:        500 * time.Millisecond,
		// IdleCheckFrequency: 500 * time.Millisecond,
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },
	})

}

// GetRedis ...
func GetRedis() *_redis.Client {
	return RedisClient
}
