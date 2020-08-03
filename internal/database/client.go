package database

import (
	"memoapp/model"
)

// Client データベースクライアントのインターフェース
type Client interface {
	Close() error
	Exists() (bool, error)
	Get() ([]byte, error)
	Set(*model.Memo) ([]byte, error)
	SetByte([]byte) error
	DEL(int) ([]byte, error)
}

var (
	pkgName = "database"
)

// Connect DB接続を行う
func Connect() (Client, error) {
	return ConnectMySql()
}
