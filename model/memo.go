package model

type Memo struct {
	ID   int    `db:"id" form:"id"`
	Memo string `db:"memo" form:"memo"`
}

func (m *Memo) SetId(id int) {
	m.ID = id
}
