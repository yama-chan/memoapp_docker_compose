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

type MemoCache struct {
	Conn redis.Conn
}

var _ Database = &MemoCache{}

func (m MemoCache) Close() error {
	return m.Conn.Close()
}

// ConnectRedis Redisへ接続する
func ConnectRedis() (Database, error) {

	// 環境変数
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return &MemoCache{}, errors.New("REDIS_HOST enviroment value is blank")
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		return &MemoCache{}, errors.New("REDIS_PORT enviroment value is blank")
	}
	// DB接続
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Printf("error: failed to connect Redis : %v\n", err)
		panic(err)
	}
	log.Println("info: Redisデータベースに接続しました")
	return &MemoCache{Conn: conn}, nil
}

// Set メモをキャッシュする
func (m *MemoCache) Set(memo *model.Memo) ([]byte, error) {
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
func (m *MemoCache) SetByte(bytes []byte) error {
	if m.Conn == nil {
		return errors.New("not initialized redis conn")
	}
	if _, err := m.Conn.Do("SET", "memos", bytes); err != nil {
		return err
	}
	return nil
}

// Exists 存在確認
func (m *MemoCache) Exists() (bool, error) {
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
func (m *MemoCache) Get() ([]byte, error) {
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
func (m *MemoCache) DEL(id int) error {
	if m.Conn == nil {
		return errors.New("not initialized redis conn")
	}

	if _, err := m.Conn.Do("DEL", "memos"); err != nil {
		return err
	}

	log.Printf("info: remove memo%d\n", id)
	return nil
}

func (s *MemoCache) set(m *model.Memo) error {
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

func (s *MemoCache) get() ([]byte, error) {
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
