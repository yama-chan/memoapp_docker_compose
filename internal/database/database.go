package database

import (
	"memoapp/model"
)

// Database データベースのインターフェース
type Database interface {
	Close() error
	Exists() (bool, error)
	Get() ([]byte, error)
	Set(*model.Memo) ([]byte, error)
	SetByte([]byte) error
	DEL(int) ([]byte, error)
	// Connect() (Database, error)
	// Set(*model.Memo) (sql.Result, error)
	// GetAll() ([]*model.Memo, error)
}

// Connect DB接続を行う
func Connect() (Database, error) {
	return ConnectMySql()
}
