package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"memoapp/internal/database"
	"memoapp/model"

	"log"

	"github.com/labstack/echo/v4"
)

type (
	// MemoHandler2 メモ用ハンドラー
	MemoHandler2 struct {
		echo        *echo.Echo
		middlewares []Middleware
		client      database.Client
	}

	endPointHandler2 func(c echo.Context) ([]byte, error)
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
	}{
		{
			"GET",
			"/list",
			hdr.MemoIndex,
		},
		{
			"POST",
			"/",
			hdr.MemoCreate,
		},
		{
			"DELETE",
			"/:id",
			hdr.MemoDelete,
		},
	}
	for _, r := range routes {
		e.Add(r.method, r.path, r.handlerFunc)
	}
	// e.GET("/list", hdr.cacheEndpointHandler(hdr.MemoIndex))
	// e.POST("/", hdr.endpointHandler(hdr.MemoCreate))
	// e.DELETE("/:id", hdr.endpointHandler(hdr.MemoDelete))
	return hdr
}

func (h *MemoHandler2) MemoIndex(c echo.Context) error {
	cacheInfo := CacheConfig{
		useCache:    true,
		setCache:    true,
		removeCache: false,
	}
	// memos, err := h.client.Get()
	memos, err := GetOps(c, cacheInfo)
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
		memo = &model.Memo{}
	)

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

	memoData, err := h.client.Set(memo)
	if err != nil {
		log.Printf("error: pkg=%s データ挿入エラー : %v\n", pkgName2, err)
		return fmt.Errorf("failed to insert memo data :[%s] %w\n ", pkgName2, err)
	}

	log.Printf("info: pkg=%s データ作成OK\n", pkgName2)
	return c.JSONBlob(http.StatusOK, memoData)
}

// MemoDelete メモ削除
func (h *MemoHandler2) MemoDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("error: データ型の変換エラー（int） : pkg=%s %v\n", pkgName2, err)
		return fmt.Errorf("failed to converted to type int :[%s] %w\n ", pkgName2, err)
	}

	memoID, err := h.client.DEL(id)
	if err != nil {
		log.Printf("error: データ削除エラー :[%s] %v\n", pkgName2, err)
		return fmt.Errorf("failed to delete memo data: [%s] %w\n ", pkgName2, err)
	}

	log.Printf("info: pkg=%s データ削除OK", pkgName2)
	return c.JSONBlob(http.StatusOK, memoID)
}
