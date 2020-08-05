package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"memoapp/internal/types"
	"memoapp/model"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type (

	// MemoAppOutput レスポンス用のデータ型
	// TODO: 最終的にoutputはこれにする
	MemoAppOutput struct {
		Results types.Contents `json:"Results"`
		Message string         `json:"Message"`
	}

	// MySQLClient MySQL接続用の構造体
	MySQLClient struct {
		DB *sqlx.DB
	}
)

var _ Client = MySQLClient{}

func ConnectMySql() (MySQLClient, error) {

	// 環境変数
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return MySQLClient{}, errors.New("DSN enviroment value is blank")
	}

	// DB接続 + 疎通確認
	//db, err := sqlx.Connect("mysql", dsn) //sqlx.Connectでsqlx.Openとdb.Ping()をやっているので修正してもいいかも
	// if err != nil {
	// 	log.Printf("error: failed to connect database: %w\n", err)
	// 	return MySQLClient{}, err
	// }

	// DB接続
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Printf("error: failed to open database connection: %v\n", err)
		return MySQLClient{}, err
	}

	// 疎通確認
	if err := db.Ping(); err != nil {
		log.Printf("error: failed to Ping verifies a connection : %v\n", err)
		return MySQLClient{}, err
	}

	// 参考：Go勉強会
	// コネクションの有効期限を設定しておかないと、
	// 死んだコネクションをいつまでも持ち続けるので設定するほうが良い。
	db.SetConnMaxLifetime(time.Minute)

	log.Println("info: MySQLデータベースに接続しました")
	return MySQLClient{DB: db}, nil
}

// Close 接続を閉じる
func (m MySQLClient) Close() error {
	return m.DB.Close()
}

// Exists 存在確認
func (m MySQLClient) Exists(params url.Values) (bool, error) {
	return false, nil
}

// func (m Memorepo) Set(memo *model.Memo) (sql.Result, error) {
func (m MySQLClient) Set(memo *model.Memo) ([]byte, error) {

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
		log.Printf("error: failed to get LastInsertId : %v\n", err)
		return nil, err
	}
	memo.SetID(int(id)) // idをセット

	bytes, err := json.Marshal(MemoAppOutput{
		Results: types.Results{memo},
	})
	if err != nil {
		log.Printf("error: failed to Marshal memo: %v\n", types.Results{memo})
		return nil, err
	}

	return bytes, nil
}

// func (m Memorepo) Set(memo *model.Memo) (sql.Result, error) {
func (m MySQLClient) SetByte(params url.Values, data []byte) error {

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

// func (m Memorepo) GetAll() (types.Results, error) {
func (m MySQLClient) Get(params url.Values) ([]byte, error) {
	var (
		memo = params.Get("memo")
		err  error
	)

	query := "SELECT * FROM memos"
	if memo != "" {
		query += " WHERE memo LIKE " + "'%" + memo + "%'"
	}
	query += " ORDER BY id asc;"
	log.Printf("info: Select query -> %v\n ", query)

	// make 関数の第 1 引数([]int)が型、第 2 引数(length)が 長さ 、第 3 引数(capacity)が 容量 を意味しています。
	// 長さ が 容量 を超えた時に、その時の 容量 の倍の 容量 が新たに確保される
	// append 関数だけを使って要素を追加していくときには、長さは 0 に指定しておく
	// make 関数で長さを 0 以外の値にしたとき、初期の長さ分の要素を考慮した作りする必要があります。
	// 参考：https://qiita.com/hitode7456/items/562527069e13347b89c8

	//予め容量を1０としている理由はLIMITがあるから？しかし長さは０なので空のスライスができる
	// つまり空のスライスが出来るがappendしていって長さが10を超えた場合は容量が倍になる設定
	memos := make([]*model.Memo, 0)
	err = m.DB.Select(&memos, query) //Select関数内でappendしているので長さは０で可変にする
	if err != nil {                  //Select関数内でappendしているので長さは０で可変にする
		log.Printf("error: failed to NamedQuery: [%s] %v\n ", pkgName, err)
		return nil, fmt.Errorf("failed to NamedQuery: [%s] %w\n ", pkgName, err)
	}

	bytes, err := json.Marshal(
		MemoAppOutput{
			// Results: results,
			Results: memos,
		})
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (m MySQLClient) DEL(params url.Values) ([]byte, error) {
	var (
		id = params.Get("id")
	)
	memoID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("error: failed to converted to type int: pkg=%s %v\n", pkgName, err)
		return nil, fmt.Errorf("failed to converted to type int: [%s] %w\n ", pkgName, err)
	}

	query := "DELETE FROM memos WHERE id = ?"

	tx := m.DB.MustBegin()
	if _, err := tx.Exec(query, memoID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("fail to Delete Exec Query: %w", err)
	}

	tx.Commit()
	bytes, err := json.Marshal(
		MemoAppOutput{
			Results: types.Results{
				&model.Memo{
					ID: memoID,
				},
			},
		})
	if err != nil {
		return nil, fmt.Errorf("fail to Marshal json: %w", err)
	}
	return bytes, nil
}
