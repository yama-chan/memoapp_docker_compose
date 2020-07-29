package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"memoapp/internal/types"
	"memoapp/model"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

type Memorepo struct {
	DB *sqlx.DB
}

var _ Database = Memorepo{}

func (m Memorepo) Close() error {
	return m.DB.Close()
}

func ConnectMySql() (Database, error) {

	// 環境変数
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return Memorepo{}, errors.New("DSN enviroment value is blank")
	}

	//db, err := sqlx.Connect("mysql", dsn) //sqlx.Connectでsqlx.Openとdb.Ping()をやっているので修正してもいいかも
	// DB接続
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Printf("error: failed to open database connection: %v\n", err)
		return Memorepo{}, err
	}

	// 参考：Go勉強会
	// コネクションの有効期限を設定しておかないと、
	// 死んだコネクションをいつまでも持ち続けるので設定するほうが良い。
	db.SetConnMaxLifetime(time.Minute)

	// 疎通確認
	if err := db.Ping(); err != nil {
		log.Printf("error: failed to Ping verifies a connection : %v\n", err)
		return Memorepo{}, err
	}

	log.Println("info: MySQLデータベースに接続しました")
	return Memorepo{DB: db}, nil
}

// Exists 存在確認
func (m Memorepo) Exists() (bool, error) {
	return false, nil
}

// func (m Memorepo) Set(memo *model.Memo) (sql.Result, error) {
func (m Memorepo) Set(memo *model.Memo) ([]byte, error) {

	query := `INSERT INTO memos (memo)
		VALUES (:memo);`

	tx := m.DB.MustBegin()
	res, err := tx.NamedExec(query, memo)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("failed to get LastInsertId : %v\n", err)
		return nil, err
	}
	memo.SetID(int(id)) // idをセット

	bytes, err := json.Marshal(types.Memos{memo})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// func (m Memorepo) Set(memo *model.Memo) (sql.Result, error) {
func (m Memorepo) SetByte(data []byte) error {

	// query := `INSERT INTO memos (memo)
	// 	VALUES (:memo);`

	// type Memo struct {
	// 	ID   int    `json:"ID"`
	// 	Memo string `json:"Memo"`
	// }
	// var memo *Memo
	// if err := json.Unmarshal(data, memo); err != nil {
	// 	return err
	// }

	// tx := m.DB.MustBegin()
	// _, err := tx.NamedExec(query, memo)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// tx.Commit()

	return nil
}

// func (m Memorepo) GetAll() ([]*model.Memo, error) {
func (m Memorepo) Get() ([]byte, error) {

	// TODO: カーソル?
	// if cursor <= 0 {
	// 	cursor = math.MaxInt32
	// }

	// TODO: プライマリーキーとは別にidがある理由は？ORDER BYはIDでなく作成日付でやるように修正
	query := `SELECT * FROM memos ORDER BY id asc;`

	// make 関数の第 1 引数([]int)が型、第 2 引数(length)が 長さ 、第 3 引数(capacity)が 容量 を意味しています。
	// 長さ が 容量 を超えた時に、その時の 容量 の倍の 容量 が新たに確保される
	// append 関数だけを使って要素を追加していくときには、長さは 0 に指定しておく
	// make 関数で長さを 0 以外の値にしたとき、初期の長さ分の要素を考慮した作りする必要があります。
	// 参考：https://qiita.com/hitode7456/items/562527069e13347b89c8

	//予め容量を1０としている理由はLIMITがあるから？しかし長さは０なので空のスライスができる
	// つまり空のスライスが出来るがappendしていって長さが10を超えた場合は容量が倍になる設定
	memos := make([]*model.Memo, 0)
	if err := m.DB.Select(&memos, query); err != nil { //Select関数内でappendしているので長さは０で可変にする
		return nil, err
	}

	bytes, err := json.Marshal(memos)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// func (m Memorepo) DEL(id int) error {
func (m Memorepo) DEL(id int) ([]byte, error) {
	query := "DELETE FROM memos WHERE id = ?"

	tx := m.DB.MustBegin()
	if _, err := tx.Exec(query, id); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("fail to Delete Exec Query: %w", err)
	}
	tx.Commit()
	bytes, err := json.Marshal(types.Memos{
		&model.Memo{
			ID: id,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fail to Marshal json: %w", err)
	}
	return bytes, nil
}
