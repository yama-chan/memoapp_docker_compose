package repository

import (
	"database/sql"

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

func GetMemoList() ([]*model.Memo, error) {

	// TODO: カーソル? ORDER BYはIDでなく作成日付でやるように修正
	// if cursor <= 0 {
	// 	cursor = math.MaxInt32
	// }

	// TODO: プライマリーキーとは別にidがある理由は？
	query := `SELECT * FROM memos ORDER BY id desc;`

	// make 関数の第 1 引数([]int)が型、第 2 引数(length)が 長さ 、第 3 引数(capacity)が 容量 を意味しています。
	// 長さ が 容量 を超えた時に、その時の 容量 の倍の 容量 が新たに確保される
	// append 関数だけを使って要素を追加していくときには、長さは 0 に指定しておく
	// make 関数で長さを 0 以外の値にしたとき、初期の長さ分の要素を考慮した作りする必要があります。
	// 参考：https://qiita.com/hitode7456/items/562527069e13347b89c8

	//予め容量を1０としている理由はLIMITがあるから？しかし長さは０なので空のスライスができる
	// つまり空のスライスが出来るがappendしていって長さが10を超えた場合は容量が倍になる設定
	memos := make([]*model.Memo, 0)
	if err := db.Select(&memos, query); err != nil { //Select関数内でappendしているので長さは０で可変にする
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
