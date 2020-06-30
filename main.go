package main

import (
	"errors"
	"log"
	"os"
	"time"

	"memoapp/handler"

	// 参考：Go勉強会
	// DBのドライバパッケージを読み込む。
	// ドライバパッケージの読み込みは、mainパッケージで実施したほうが良い。
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	//サーバ起動
	err := (start_application(e))
	if err != nil {
		// サーバエラー時はFatal（= os.Exit(1)）
		log.Fatal(err)
	}
}

func connectDB(e *echo.Echo) (*sqlx.DB, error) {

	//環境変数
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return nil, errors.New("DSN enviroment value is blank")
	}

	//db, err := sqlx.Connect("mysql", dsn) //sqlx.Connectでsqlx.Openとdb.Ping()をやっているので修正してもいいかも
	//DB接続
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		e.Logger.Errorf("failed to open database connection: %v\n", err)
		// e.Logger.Fatal(err)
		return nil, err
	}

	// 参考：Go勉強会
	// コネクションの有効期限を設定しておかないと、
	// 死んだコネクションをいつまでも持ち続けるので設定するほうが良い。
	db.SetConnMaxLifetime(time.Minute)

	//疎通確認
	if err := db.Ping(); err != nil {
		e.Logger.Errorf("failed to Ping verifies a connection : %v\n", err)
		return nil, err
	}

	log.Println("データベースに接続しました")
	return db, nil
}

func start_application(e *echo.Echo) error {

	//DB接続
	db, err := connectDB(e)
	if err != nil {
		e.Logger.Errorf("failed to connection DB: %v\n", err)
		return err
	}

	//deferでClose
	defer db.Close()

	//静的ファイル
	e.Static("/styles", "src/styles")

	//ミドルウェア
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.Gzip(), //HTTPレスポンスをGzip圧縮して返す
	)

	//ルーティング
	handler.ProvideHandler(e, db)

	//ゴルーチン/チャネル
	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Start(":8080")
	}()

	//チャネルの受信待ち
	select {
	case err := <-errCh:
		return err
	}
	//TODO: gracefulにサーバ停止する処理も追加する。現状ではシグナルを考慮しない
}
