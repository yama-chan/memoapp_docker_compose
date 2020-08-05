package database

import (
	"memoapp/model"
	"net/url"
)

// Client データベースクライアントのインターフェース
type Client interface {
	Set(*model.Memo) ([]byte, error)
	Get(url.Values) ([]byte, error)
	DEL(url.Values) ([]byte, error)
	Exists(url.Values) (bool, error)
	SetByte(url.Values, []byte) error
	Close() error
}

var (
	pkgName = "database"
)

// CheckCache キャッシュの有無の確認を行う
func CheckCache() (Client, error) {
	return nil, nil
}
