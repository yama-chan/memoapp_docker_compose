package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"memoapp/internal/types"
	"memoapp/model"
	"net/url"
	"os"

	"github.com/gomodule/redigo/redis"
)

type CacheClient struct {
	Conn redis.Conn
}

var (
	_ Client = CacheClient{}
)

const (
	keyprefix string = "memoapp;getmemo/"
)

// ConnectRedis Redisへ接続する
func ConnectRedis() (CacheClient, error) {

	// 環境変数
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return CacheClient{}, errors.New("REDIS_HOST enviroment value is blank")
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		return CacheClient{}, errors.New("REDIS_PORT enviroment value is blank")
	}
	// DB接続
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Printf("error: failed to connect Redis : %v\n", err)
		return CacheClient{}, fmt.Errorf("failed to connect Redis : %w\n ", err)
	}
	log.Println("info: Redisデータベースに接続しました")
	return CacheClient{Conn: conn}, nil
}

// Close 接続を閉じる
func (m CacheClient) Close() error {
	return m.Conn.Close()
}

// Set メモをキャッシュする
func (m CacheClient) Set(memo *model.Memo) ([]byte, error) {

	if m.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}

	bytes, err := json.Marshal(types.Results{memo})
	if err != nil {
		return nil, err
	}

	if _, err := m.Conn.Do("SET", "memos", bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

// SetByte バイト配列をキャッシュする
func (m CacheClient) SetByte(params url.Values, bytes []byte) error {
	if m.Conn == nil {
		return errors.New("not initialized redis conn")
	}

	cacheKey := createCacheKey(params)
	log.Printf("info: SetByte CacheKey -> %s", cacheKey)

	if _, err := m.Conn.Do("SET", cacheKey, bytes); err != nil {
		return err
	}
	return nil
}

// Exists キャッシュの存在確認
func (m CacheClient) Exists(params url.Values) (bool, error) {
	if m.Conn == nil {
		return false, errors.New("not initialized redis conn")
	}

	cacheKey := createCacheKey(params)
	log.Printf("info: Exists CacheKey -> %s", cacheKey)

	exists, err := redis.Bool(m.Conn.Do("EXISTS", cacheKey))
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Get キャッシュデータを取得
func (m CacheClient) Get(params url.Values) ([]byte, error) {
	if m.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}

	cacheKey := createCacheKey(params)
	log.Printf("info: Get CacheKey -> %s", cacheKey)

	bytes, err := redis.Bytes(m.Conn.Do("GET", cacheKey))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// DEL キーを指定してキャッシュを削除
func (m CacheClient) DEL(params url.Values) ([]byte, error) {
	if m.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}

	cacheKey := createCacheKey(params)
	log.Printf("info: DEL CacheKey -> %s", cacheKey)

	if _, err := m.Conn.Do("DEL", cacheKey); err != nil {
		return nil, err
	}

	// 接続先のRedisを一括削除
	if _, err := m.Conn.Do("flushdb"); err != nil {
		return nil, err
	}
	return nil, nil
}

// Flush キャッシュを一括削除
func (m CacheClient) Flush() error {
	if m.Conn == nil {
		return errors.New("not initialized redis conn")
	}

	// 接続先のRedisを一括削除
	if _, err := m.Conn.Do("flushdb"); err != nil {
		return err
	}
	return nil
}

func createCacheKey(params url.Values) string {
	return keyprefix + params.Encode()
}
