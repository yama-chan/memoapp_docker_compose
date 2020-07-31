package main

import (
	"context"
	"log"
	"memoapp/internal/handler"
	"os"
	"os/signal"
	"syscall"

	// 参考：Go勉強会
	// DBのドライバパッケージを読み込む。
	// ドライバパッケージの読み込みは、mainパッケージで実施したほうが良い。
	"github.com/comail/colog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)

	e := echo.New()

	//サーバ起動
	err := (start_application(e, ":8080"))
	if err != nil {
		// サーバエラー時はFatal（= os.Exit(1)）
		e.Logger.Fatal(err)
	}
}

// start_application アプリケーションを起動する
func start_application(e *echo.Echo, port string) error {

	//  ハンドラー生成
	handler.ProvideHandler(e)
	// handler.ProvideHandler(e)

	// echoログのフォーマット
	logger := middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: logFormat(),
		Output: os.Stdout,
	})

	// ミドルウェア
	e.Use(
		logger,
		middleware.Recover(),
		middleware.Gzip(), //HTTPレスポンスをGzip圧縮して返す
		// hdr.TestGetBodyDump(),
		// hdr.WithContextGen(),
		// hdr.WithProviderFinalizer(),
	)
	// 静的ファイル
	e.Static("/styles", "src/styles")

	// ファイル
	e.File("/favicon.ico", "images/favicon.ico")

	// インデックス画面を表示
	e.GET("/", index)

	// サーバー起動
	// ゴルーチン/チャネル
	errCh := make(chan error, 1)
	go func() {
		errCh <- e.Start(port)
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// チャネルの受信待ち
	select {
	case err := <-errCh:
		return err
	case <-quit:
		err := e.Shutdown(ctx)
		if err != nil {
			return err
		}
		log.Println("info: Shutdown gracefully...")
		return nil
	}
}

func index(c echo.Context) error {
	return render(c, "src/views/index.html", nil)
}

func logFormat() string {
	var format string
	format += "\n[  echo ]"
	format += "time:${time_rfc3339}\n"
	format += "- method:${method}\t"
	format += "status:${status}\n"
	format += "- error:${error}\n"
	format += "- path:${path}\t"
	format += "uri:${uri}\t"
	format += "host:${host}\t"
	format += "remote_ip:${remote_ip}\n"
	format += "- bytes_in:${bytes_in}\t"
	format += "bytes_out:${bytes_out}\n"
	format += "- latency:${latency}\t"
	format += "latency_human:${latency_human}\n\n"
	// format += "forwardedfor:${header:x-forwarded-for}\n"
	// format += "referer:${referer}\n"
	// format += "user_agent:${user_agent}\n"
	// format += "request_id:${id}\n"

	return format
}
