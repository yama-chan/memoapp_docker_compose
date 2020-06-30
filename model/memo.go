package model

import "errors"

type Memo struct {
	ID   int    `db:"id" form:"id"`
	Memo string `db:"memo" form:"memo"`
}

// SetID MemoのIDの設定を行う
func (m *Memo) SetID(id int) {
	m.ID = id
}

// Validate Memoのバリデーション関数
func (m *Memo) Validate() error {
	if m.ID >= 0 {
		return errors.New("model Memo: Validate Error")
	}
	return nil
}
