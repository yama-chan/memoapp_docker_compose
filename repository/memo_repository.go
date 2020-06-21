package repository

import (
	"database/sql"
	"math"

	"memoapp/model"
)

func MemoCreate(memo *model.Memo) (sql.Result, error) {

	query := `INSERT INTO memos (memo)
	VALUES (:memo);`

	tx := db.MustBegin()

	res, err := tx.NamedExec(query, memo)
	if err != nil {

		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return res, nil
}

func MemoListByCursor(cursor int) ([]*model.Memo, error) {

	if cursor <= 0 {
		cursor = math.MaxInt32
	}

	query := `SELECT *
	FROM memos
	WHERE id < ?
	ORDER BY id desc
	LIMIT 10`

	memos := make([]*model.Memo, 0, 10)

	if err := db.Select(&memos, query, cursor); err != nil {
		return nil, err
	}

	return memos, nil
}

func MemoDelete(id int) error {
	query := "DELETE FROM memos WHERE id = ?"

	tx := db.MustBegin()

	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()

		return err
	}

	return tx.Commit()
}
