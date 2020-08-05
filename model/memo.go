package model

import "errors"

type Memo struct {
	ID   int    `db:"id" form:"id" json:"ID"`
	Memo string `db:"memo" form:"memo" json:"Memo"`
}

// SetID MemoのIDの設定を行う
func (m *Memo) SetID(id int) {
	m.ID = id
}

// Validate Memoのバリデーション関数
func (m *Memo) Validate() error {
	if m.ID < 0 {
		return errors.New("Validate Error[model.Memo]: memo.ID should not be less then 0. ")
	}
	if m.Memo == "" {
		return errors.New("Validate Error[model.Memo]: memo.Memo should not be empty. ")
	}
	return nil
}
