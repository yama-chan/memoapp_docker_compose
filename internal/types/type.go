package types

import (
	"memoapp/model"
	"net/url"
)

type (
	Parameters url.Values
	// Results メモリスト
	Memos []*model.Memo

	Contents interface{}
	Results  []interface{}
)
