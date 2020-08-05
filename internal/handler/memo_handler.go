package handler

import (
	"fmt"

	"memoapp/internal/database"
	"memoapp/model"

	"log"

	"github.com/labstack/echo/v4"
)

type (
	// MemoHandler メモ用ハンドラー
	MemoHandler struct {
		HasCache bool
		Client   database.Client
		echo     *echo.Echo
	}

	EndPointHandler func(c echo.Context) ([]byte, error)
)

var (
	pkgName = "handler"
)

// ProvideHandler メモハンドラーからルーティングを設定する
func ProvideHandler(e *echo.Echo) *MemoHandler {
	hdr := &MemoHandler{echo: e}
	routes := []struct {
		method     string
		path       string
		callback   EndPointHandler
		cache      bool // キャッシュをするかどうか
		cacheClear bool // レスポンス返却後、キャッシュをリセットするかどうか
	}{
		{
			"GET",
			"/list",
			hdr.MemoIndex,
			true,
			false,
		},
		{
			"POST",
			"/",
			hdr.MemoCreate,
			false,
			true,
		},
		{
			"DELETE",
			"/",
			hdr.MemoDelete,
			false,
			true,
		},
	}
	for _, r := range routes {
		if r.cache {
			e.Add(r.method, r.path, hdr.cacheEndpointHandler(r.callback))
		} else {
			e.Add(r.method, r.path, hdr.endpointHandler(r.callback, r.cacheClear))
		}
	}
	// e.GET("/list", hdr.cacheEndpointHandler(hdr.MemoIndex))
	// e.POST("/", hdr.endpointHandler(hdr.MemoCreate))
	// e.DELETE("/:id", hdr.endpointHandler(hdr.MemoDelete))
	return hdr
}

func (h *MemoHandler) MemoIndex(c echo.Context) ([]byte, error) {

	memos, err := h.Client.Get(c.Request().URL.Query())
	if err != nil {
		log.Printf("error: failed to Get memo data : %v\n", err)
		return nil, fmt.Errorf("failed to Get memo data: [%s]%w\n ", pkgName, err)
	}

	log.Printf("info: pkg=%s データ取得OK\n", pkgName)
	return memos, nil

}

// MemoCreate メモ作成
func (h *MemoHandler) MemoCreate(c echo.Context) ([]byte, error) {

	var (
		memo = &model.Memo{}
	)

	err := c.Bind(memo)
	if err != nil {
		log.Printf("error: 入力データに誤りがあります。:[%s] %v\n", pkgName, err)
		return nil, fmt.Errorf("failed to Bind request params :[%s] %v\n ", pkgName, err)
	}

	// バリデートを実行
	err = memo.Validate()
	if err != nil {
		log.Printf("error: バリデーションでエラーが発生しました。:[%s] %v\n", pkgName, err)
		return nil, fmt.Errorf("validation error:[%s] %w\n ", pkgName, err)
	}

	memoData, err := h.Client.Set(memo)
	if err != nil {
		log.Printf("error: pkg=%s データ挿入エラー : %v\n", pkgName, err)
		return nil, fmt.Errorf("failed to insert memo data :[%s] %w\n ", pkgName, err)
	}

	log.Printf("info: pkg=%s データ作成OK\n", pkgName)
	return memoData, nil
}

// MemoDelete メモ削除
func (h *MemoHandler) MemoDelete(c echo.Context) ([]byte, error) {
	// id, err := strconv.Atoi(c.Param("id"))
	// if err != nil {
	// 	log.Printf("error: データ型の変換エラー（int） : pkg=%s %v\n", pkgName, err)
	// 	return nil, fmt.Errorf("failed to converted to type int :[%s] %w\n ", pkgName, err)
	// }

	memoID, err := h.Client.DEL(c.Request().URL.Query())
	if err != nil {
		log.Printf("error: データ削除エラー :[%s] %v\n", pkgName, err)
		return nil, fmt.Errorf("failed to delete memo data: [%s] %w\n ", pkgName, err)
	}

	log.Printf("info: pkg=%s データ削除OK", pkgName)
	return memoID, nil
}
