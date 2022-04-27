package RedisTool

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
	"unsafe"
)

/**
    @author: WangYuding
    @since: 2022/4/27
    @desc: //Redis的工具库
**/
var ctx = context.Background()

func NewRedis(db int) *redis.Client { //将数据库连接操作打包为方法使用newRdis(0)方法带入数据库名调用即可
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", //数据库默认安装在开发机，监听localhost，默认端口为6379
		Password: "123456",         // simply password set
		DB:       db,               // use default DB
	})
	return rdb //返回数据库客户端
}

// GetAllKeys 获取该数据库里所有的key
func GetAllKeys(db int) []string {
	rdb := NewRedis(db)
	defer rdb.Close()
	keys, err := rdb.Keys(ctx, "*").Result()
	CheckError(err)
	return keys
}

func GetValueByKey(db int, key string) (string, error) {
	rdb := NewRedis(db)
	defer rdb.Close()
	val, err := rdb.Get(ctx, key).Result() //使用IdTime获取Message
	CheckError(err)
	return val, err
}

func SetRecord(db int, key string, data []byte) bool {
	rdb := NewRedis(db)
	err := rdb.Set(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		log.Println(err)
		return false
	} //保存Message,保存24小时24*time.Hour
	err = rdb.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// CheckError error处理,可以使用客户端log处理的逻辑将错误信息收集保存到数据库,这里不在展开
func CheckError(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}

// StrTobyte 将string转为[]byte
func StrTobyte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// ByteTostr 将[]byte转为string
func ByteTostr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// DeletSomething 该方法可以删除Redis里任意数据库 Db 的任意 Key
func DeletSomething(key string, Db int) {
	rdb := NewRedis(Db)
	rdb.Del(ctx, key).Err()
	rdb.Close()
}