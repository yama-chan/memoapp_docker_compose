package main

import (
	"errors"
	"log"
	"os"
	"time"

	"memoapp/handler"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// const tmplPath = "src/view"

// var db *sqlx.DB

func main() {
	e := echo.New()

	//サーバ起動
	e.Logger.Fatal(start_application(e))
}

func connectDB(e *echo.Echo) (*sqlx.DB, error) {

	dsn := os.Getenv("DSN")
	if dsn == "" {
		return nil, errors.New("DSN enviroment value is blank")
	}

	// db, err := sqlx.Connect("mysql", dsn) //sqlx.Connectでsqlx.Openとdb.Ping()をやっているので修正してもいいかも
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

	if err := db.Ping(); err != nil {
		e.Logger.Errorf("failed to Ping verifies a connection : %v\n", err)
		return nil, err
	}

	log.Println("データベースに接続しました")
	return db, nil
}

func start_application(e *echo.Echo) error {

	// DB接続
	db, err := connectDB(e)
	if err != nil {
		e.Logger.Errorf("failed to  connection DB: %v\n", err)
		return err
	}
	defer db.Close()
	// repository.SetDB(db)

	//静的ファイル
	e.Static("/styles", "src/styles")
	//ミドルウェア
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.Gzip(),
	)
	//ルーティング
	hdr := handler.ProvideHandler(db)
	e.POST("/", hdr.MemoCreate)
	e.GET("/", hdr.MemoIndex)
	e.DELETE("/:id", hdr.MemoDelete)

	//ゴルーチン/チャネル
	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Start(":8080")
	}()

	select {
	case err := <-errCh:
		return err
	}
}
