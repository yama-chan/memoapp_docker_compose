package handler

import (
	"fmt"
	"memoapp/model"
	"net/http"

	"memoapp/internal/database"

	"log"

	"github.com/labstack/echo/v4"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

type (
	// MemoHandler2 メモ用ハンドラー2
	MemoHandler2 struct {
		HasCache bool
		echo     *echo.Echo
		// middlewares []Middleware
		// client      database.Client
	}
	Middlewares []echo.MiddlewareFunc
)

var (
	pkgName2 = "handler2"
)

// ProvideHandler メモハンドラーからルーティングを設定する
func ProvideHandler2(e *echo.Echo) *MemoHandler2 {
	hdr := &MemoHandler2{echo: e}
	routes := []struct { // ルート関数ごとに設定値を記載
		method      string
		path        string
		handlerFunc echo.HandlerFunc
		middlewares Middlewares
	}{
		{
			"GET",
			"/list",
			hdr.MemoIndex,
			Middlewares{
				hdr.UseCache(),
				hdr.SetCache(),
			},
		},
		{
			"POST",
			"/",
			hdr.MemoCreate,
			Middlewares{
				UseMySQL(),
				hdr.ClearCache(),
			},
		},
		{
			"DELETE",
			"/",
			hdr.MemoDelete,
			Middlewares{
				UseMySQL(),
				hdr.ClearCache(),
			},
		},
	}
	for _, r := range routes {
		e.Add(r.method, r.path, r.handlerFunc, r.middlewares...)
	}
	// e.GET("/list", hdr.cacheEndpointHandler(hdr.MemoIndex))
	// e.POST("/", hdr.endpointHandler(hdr.MemoCreate))
	// e.DELETE("/:id", hdr.endpointHandler(hdr.MemoDelete))
	return hdr
}

func (h *MemoHandler2) MemoIndex(c echo.Context) error {
	var (
		ctx        = c.Request().Context()
		parameters = c.Request().URL.Query()
	)

	// contextからデータベースのクライアントを取得
	client, ok := ctx.Value("storeKey").(database.Client)
	if !ok {
		return fmt.Errorf("failed to Get client from context: [%s]\n ", pkgName2)
	}

	// 取得
	memos, err := client.Get(parameters)
	if err != nil {
		log.Printf("error: failed to Get memo data : %v\n", err)
		return fmt.Errorf("failed to Get memo data: [%s]%w\n ", pkgName2, err)
	}

	log.Printf("info: pkg=%s データ取得OK\n", pkgName2)
	return c.JSONBlob(http.StatusOK, memos)
}

// MemoCreate メモ作成
func (h *MemoHandler2) MemoCreate(c echo.Context) error {

	var (
		ctx  = c.Request().Context()
		memo = &model.Memo{}
	)

	// フォーム値とバインド
	err := c.Bind(memo)
	if err != nil {
		log.Printf("error: 入力データに誤りがあります。:[%s] %v\n", pkgName2, err)
		return fmt.Errorf("failed to Bind request params :[%s] %v\n ", pkgName2, err)
	}

	// バリデートを実行
	err = memo.Validate()
	if err != nil {
		log.Printf("error: バリデーションでエラーが発生しました。:[%s] %v\n", pkgName2, err)
		return fmt.Errorf("validation error:[%s] %w\n ", pkgName2, err)
	}

	// contextからデータベースのクライアントを取得
	client, ok := ctx.Value("storeKey").(database.Client)
	if !ok {
		return fmt.Errorf("failed to Get client from context: [%s]\n ", pkgName2)
	}

	// 作成
	memoData, err := client.Set(memo)
	if err != nil {
		log.Printf("error: pkg=%s データ挿入エラー : %v\n", pkgName2, err)
		return fmt.Errorf("failed to insert memo data :[%s] %w\n ", pkgName2, err)
	}

	log.Printf("info: pkg=%s データ作成OK\n", pkgName2)
	return c.JSONBlob(http.StatusOK, memoData)
}

// MemoDelete メモ削除
func (h *MemoHandler2) MemoDelete(c echo.Context) error {
	var (
		ctx        = c.Request().Context()
		parameters = c.Request().URL.Query()
	)

	// contextからデータベースのクライアントを取得
	client, ok := ctx.Value("storeKey").(database.Client)
	if !ok {
		return fmt.Errorf("failed to Get client from context: [%s]\n ", pkgName2)
	}

	// 削除
	memoID, err := client.DEL(parameters)
	if err != nil {
		log.Printf("error: データ削除エラー :[%s] %v\n", pkgName2, err)
		return fmt.Errorf("failed to delete memo data: [%s] %w\n ", pkgName2, err)
	}

	log.Printf("info: pkg=%s データ削除OK", pkgName2)
	return c.JSONBlob(http.StatusOK, memoID)
}
