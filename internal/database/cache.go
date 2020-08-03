package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"memoapp/internal/types"
	"memoapp/model"
	"os"

	"github.com/gomodule/redigo/redis"
)

type CacheClient struct {
	Conn redis.Conn
}

var _ Client = &CacheClient{}

// ConnectRedis Redisへ接続する
func ConnectRedis() (Client, error) {

	// 環境変数
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return &CacheClient{}, errors.New("REDIS_HOST enviroment value is blank")
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		return &CacheClient{}, errors.New("REDIS_PORT enviroment value is blank")
	}
	// DB接続
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Printf("error: failed to connect Redis : %v\n", err)
		panic(fmt.Errorf("failed to connect Redis : %w\n ", err))
	}
	log.Println("info: Redisデータベースに接続しました")
	return &CacheClient{Conn: conn}, nil
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

	bytes, err := json.Marshal(types.Memos{memo})
	if err != nil {
		return nil, err
	}

	if _, err := m.Conn.Do("SET", "memos", bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

// SetByte バイト配列をキャッシュする
func (m CacheClient) SetByte(bytes []byte) error {
	if m.Conn == nil {
		return errors.New("not initialized redis conn")
	}
	if _, err := m.Conn.Do("SET", "memos", bytes); err != nil {
		return err
	}
	return nil
}

// Exists キャッシュの存在確認
func (m CacheClient) Exists() (bool, error) {
	if m.Conn == nil {
		return false, errors.New("not initialized redis conn")
	}
	exists, err := redis.Bool(m.Conn.Do("EXISTS", "memos"))
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Get キャッシュデータを取得
func (m CacheClient) Get() ([]byte, error) {
	if m.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}

	// exists, err := m.Exists()
	// if err != nil {
	// 	return nil, err
	// }
	// if !exists {
	// 	return nil, nil
	// }

	bytes, err := redis.Bytes(m.Conn.Do("GET", "memos"))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// DEL キャッシュを削除
func (m CacheClient) DEL(id int) ([]byte, error) {
	if m.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}
	if _, err := m.Conn.Do("DEL", "memos"); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s CacheClient) set(m *model.Memo) error {
	if s.Conn == nil {
		return errors.New("not initialized redis conn")
	}

	bytes, err := json.Marshal(types.Memos{m})
	if err != nil {
		return err
	}

	if _, err := s.Conn.Do("SET", "memos", bytes); err != nil {
		return err
	}

	return nil
}

func (s CacheClient) get() ([]byte, error) {
	if s.Conn == nil {
		return nil, errors.New("not initialized redis conn")
	}

	// exists, err := s.Exists()
	// if err != nil {
	// 	return nil, err
	// }
	// if !exists {
	// 	return nil, nil
	// }

	bytes, err := redis.Bytes(s.Conn.Do("GET", "memos"))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
