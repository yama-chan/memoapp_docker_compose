package model

type Memo struct {
	ID   int    `db:"id" form:"id"`
	Memo string `db:"memo" form:"memo"`
}
