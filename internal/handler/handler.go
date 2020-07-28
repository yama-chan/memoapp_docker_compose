package handler

// Handler データベースのインターフェース
type Handler interface {
	Init() error
	Connect() error
	Get() error
	Set() error
}
