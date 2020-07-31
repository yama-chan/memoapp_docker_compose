package handler

import (
	"fmt"
	"log"
	"memoapp/internal/database"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ***********************************************************************
// https://medium.com/veltra-engineering/echo-middleware-in-golang-90e1d301eb27

//	Middlewareの実行順序

// 	middleware-Pre  : before
// 	middleware-Use-1: before
// 	middleware-Use-2: before
// 	middleware-Group: before
// 	middleware-Route: before
// 	logic: main
// 	logic: defer
// 	middleware-Route: after
// 	middleware-Route: defer
// 	middleware-Group: after
// 	middleware-Group: defer
// 	middleware-Use-2: after
// 	middleware-Use-2: defer
// 	middleware-Use-1: after
// 	middleware-Use-1: defer
// 	middleware-Pre  : after
// 	middleware-Pre  : defer

//	★ 'Pre'→'Use'→'Group'→'Route'の順
//	★ 'Use'で設定された2つについては、先に設定したものから実行されている
//	★ 'defer'が実行されるタイミングは当該Middlewareの事後処理('after')直後

// ***********************************************************************

func (h *MemoHandler) WithContextGen() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		//defer内部で発生したerrorを処理するのには名前付き返り値を利用する。
		return func(c echo.Context) error {
			fmt.Println("c.Path(): " + c.Path())
			fmt.Println("c.RealIP(): " + c.RealIP())
			fmt.Printf("c.ParamValues(): %v\n", c.ParamValues())
			fmt.Printf("c.QueryParams(): %v\n", c.QueryParams())
			fmt.Printf("c.QueryString(): %v\n", c.QueryString())

			// DB接続
			database, err := database.Connect()
			if err != nil {
				log.Printf("error: failed to connection DB: %v\n", err)
				return err
			}

			// DBのClose処理
			defer func() error {
				err = database.Close()
				if err != nil {
					log.Printf("error: failed to Close DB: %v\n", err)
					return err
				}
				log.Println("info: database connection is Closed")
				return nil
			}()
			h.repo = database

			// ↑ BEFORE
			// この場合、HandlerFuncが実行されてReturnとなる
			return next(c) // HandlerFunc : func(Context) error
			// この場合、AFTERの処理は実行され、エラーを返す
			// ↓ AFTER
		}
	}
}

func (h *MemoHandler) cacheEndpointHandler(handler endPointHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.HasCache = false

		fmt.Println("c.Path(): " + c.Path())
		fmt.Println("c.RealIP(): " + c.RealIP())
		fmt.Printf("c.ParamValues(): %v\n", c.ParamValues())
		fmt.Printf("c.QueryParams(): %v\n", c.QueryParams())
		fmt.Printf("c.QueryString(): %v\n", c.QueryString())

		// Redis接続
		redis, err := database.ConnectRedis()
		if err != nil {
			log.Printf("error: failed to connection DB: %v\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		// RedisのClose処理
		defer redis.Close()

		// キャッシュの確認
		cached, err := redis.Exists()
		if err != nil {
			log.Printf("error: failed to Get cached data : %v\n", err)
			return fmt.Errorf("failed to Get cached data: [%s]%w\n ", pkgName, err)
		}
		// キャッシュが存在する場合
		if cached {
			log.Printf("info: Found form Redis Memo cached data")
			h.repo = redis
			data, err := handler(c)
			if err != nil {
				log.Printf("error: Internal Server Error: %v\n", err)
				return c.NoContent(http.StatusInternalServerError)
			}
			h.HasCache = true
			if data == nil {
				return c.NoContent(http.StatusNoContent)
			}
			return c.JSONBlob(http.StatusOK, data)
		}
		// キャッシュが存在しない場合
		log.Printf("info: Not Found form Redis Memo cached data")
		data, err := h.execMySQLHandler(handler, c)
		// レスポンスが書き込まれた後にキャッシュに格納
		if err != nil {
			log.Printf("error: Internal Server Error: %v\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if data == nil {
			return c.NoContent(http.StatusNoContent)
		}
		c.Response().After(func() {
			// after Response
			err = redis.SetByte(data)
			if err != nil {
				// とりあえずログのみ出力
				log.Printf("error: fail to SetByte Error: %v\n", err)
				// return c.NoContent(http.StatusInternalServerError)
			}
			log.Println("info: memo data is cached")
		})
		return c.JSONBlob(http.StatusOK, data)
	}
}

func (h *MemoHandler) endpointHandler(handler endPointHandler, cacheClear bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.HasCache = false

		fmt.Println("c.Path(): " + c.Path())
		fmt.Println("c.RealIP(): " + c.RealIP())
		fmt.Printf("c.ParamValues(): %v\n", c.ParamValues())
		fmt.Printf("c.QueryParams(): %v\n", c.QueryParams())
		fmt.Printf("c.QueryString(): %v\n", c.QueryString())

		data, err := h.execMySQLHandler(handler, c)
		if err != nil {
			log.Printf("error: Internal Server Error: %v\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		// ハンドラーにエラーが無ければキャッシュをクリア
		// if cacheClear {
		// 	err = h.clearRedisCache()
		// 	if err != nil {
		// 		// とりあえずログのみ出力
		// 		log.Printf("error: fail to clear cache: %v\n", err)
		// 		// return c.NoContent(http.StatusInternalServerError)
		// 	}
		// 	log.Println("info: cache is cleared")
		// }

		// c.response.Write(byte)が呼ばれた場合に以下のAfter(func())が実行される
		// ※ c.NoContentだとWriteがよばれないので実行されない。
		c.Response().After(func() {
			if cacheClear {
				err = h.clearRedisCache()
				if err != nil {
					// とりあえずログのみ出力
					log.Printf("error: fail to clear cache: %v\n", err)
					// return c.NoContent(http.StatusInternalServerError)
				}
				log.Println("info: cache is cleared")
			}
		})

		// if data == nil {
		// 	return c.NoContent(http.StatusOK)
		// }
		return c.JSONBlob(http.StatusOK, data)
	}
}

func (h *MemoHandler) execMySQLHandler(handler endPointHandler, c echo.Context) ([]byte, error) {
	// MySQLに接続
	database, err := database.ConnectMySql()
	if err != nil {
		log.Printf("error: failed to connection DB: %v\n", err)
		return nil, fmt.Errorf("failed to connection DB: %w", err)
	}
	// RedisのClose処理
	defer database.Close()

	h.repo = database
	return handler(c)
}

func (h *MemoHandler) clearRedisCache() error {
	var (
		redis database.Database
	)

	// 型チェック
	switch h.repo.(type) {
	case database.MemoCache:
		redis = h.repo
	default:
		// Redisに接続
		r, err := database.ConnectRedis()
		if err != nil {
			log.Printf("error: Redisへの接続に失敗しました。: %v\n", err)
			return fmt.Errorf("failed to connection Redis: %w", err)
		}
		redis = r
	}
	// RedisのClose処理
	defer redis.Close()

	// キャッシュの削除
	_, err := redis.DEL(0)
	return err
}

// func (h *MemoHandler) TestGetBodyDump() echo.MiddlewareFunc {
// 	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
// 		// var res map[string]interface{}
// 		log.Printf("info: resBody: %v", resBody)
// 		var out database.MemoAppOutput
// 		err := json.Unmarshal(resBody, &out)
// 		if err != nil {
// 			log.Printf("error: fail to json.Unmarshal: %v\n", err)
// 		}
// 		if h.HasCache {
// 			out.Message = "cache data"
// 			// b, _ := json.Marshal(out)
// 			// c.Response().Flush()
// 			// c.Response().Writer.Write(b)
// 			log.Printf("info: response: %v", out)
// 		}
// 	})
// }
