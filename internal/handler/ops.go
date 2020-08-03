package handler

import (
	"fmt"
	"log"
	"memoapp/internal/database"
	"memoapp/model"

	"github.com/labstack/echo/v4"
)

type (
	// ChainHandler
	// ChainHandler struct {
	// 	Mdws     []Middleware
	// 	Endpoint Middleware
	// 	chain    Middleware
	// 	info     CacheConfig
	// }

	// step func(step) step

	CacheHandler struct {
		HasCache bool
		cache    database.Client
		client   database.Client
		// Mdws     []Middleware
		info CacheConfig
	}

	// CacheConfig データベース実行用のクライアント
	CacheConfig struct {
		useCache    bool
		setCache    bool
		removeCache bool
	}

	Middlewares []Middleware

	getOps    func() ([]byte, error)
	setOps    func(*model.Memo) ([]byte, error)
	deleteOps func(int) ([]byte, error)
)

func GetOps(ctx echo.Context, info CacheConfig) ([]byte, error) {
	hdr, err := New(info)
	if err != nil {
		return nil, fmt.Errorf("failed to New *CacheHandler: [%s]%w\n ", pkgName2, err)
	}

	data, err := hdr.client.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to New *CacheHandler: [%s]%w\n ", pkgName2, err)
	}

	if info.setCache {
		err = hdr.setCache(ctx, data)
		if err != nil {
			return data, fmt.Errorf("failed to Set Cache: [%s]%w\n ", pkgName2, err)
		}
	}

	if info.removeCache {
		err = hdr.clearCache(ctx)
		if err != nil {
			return data, fmt.Errorf("failed to Clear Cache: [%s]%w\n ", pkgName2, err)
		}
	}

	return data, err
}
func SetOps(ctx echo.Context, info CacheConfig, memo *model.Memo) ([]byte, error) {
	hdr, err := New(info)
	if err != nil {
		return nil, fmt.Errorf("failed to New *CacheHandler: [%s]%w\n ", pkgName2, err)
	}

	data, err := hdr.client.Set(memo)
	if err != nil {
		return nil, fmt.Errorf("failed to New *CacheHandler: [%s]%w\n ", pkgName2, err)
	}
	//一旦Postではキャッシュしない
	// if info.setCache {
	// 	err = hdr.setCache(ctx, data)
	// 	if err != nil {
	// 		return data, fmt.Errorf("failed to Set Cache: [%s]%w\n ", pkgName2, err)
	// 	}
	// }

	if info.removeCache {
		err = hdr.clearCache(ctx)
		if err != nil {
			return data, fmt.Errorf("failed to Clear Cache: [%s]%w\n ", pkgName2, err)
		}
	}

	return data, err
}

func New(info CacheConfig) (*CacheHandler, error) {
	var (
		redis database.Client
	)
	if info.useCache {
		redis, err := database.ConnectRedis()
		if err != nil {
			log.Printf("error: failed to Connect Redis : %v\n", err)
			return nil, fmt.Errorf("failed to Connect Redis: [%s]%w\n ", pkgName2, err)
		}
		cached, err := redis.Exists()
		if err != nil {
			log.Printf("error: failed to Get cached data : %v\n", err)
			return nil, fmt.Errorf("failed to Get cached data: [%s]%w\n ", pkgName2, err)
		}
		if cached {
			log.Printf("info: Found form Redis Memo cached data")
			return &CacheHandler{
				info:   info,
				cache:  redis,
				client: redis,
			}, nil
		}
		log.Printf("info: Not Found form Redis Memo cached data")
	}
	mysql, err := database.ConnectMySql()
	if err != nil {
		log.Printf("error: failed to Connect MySql : %v\n", err)
		return nil, fmt.Errorf("failed to Connect MySql: [%s]%w\n ", pkgName2, err)
	}
	return &CacheHandler{
		info:   info,
		cache:  redis,
		client: mysql,
	}, nil
}

func (c *CacheHandler) setCache(ctx echo.Context, data []byte) error {
	if c.cache == nil {
		redis, err := database.ConnectRedis()
		if err != nil {
			log.Printf("error: failed to Connect Redis : %v\n", err)
			return fmt.Errorf("failed to Connect Redis: [%s]%w\n ", pkgName2, err)
		}
		c.cache = redis
	}
	ctx.Response().After(func() {
		// after Response
		err := c.cache.SetByte(data)
		if err != nil {
			// とりあえずログのみ出力
			log.Printf("error: fail to SetByte Error: %v\n", err)
			// return c.NoContent(http.StatusInternalServerError)
		}
		log.Println("info: memo data is cached")
	})
	return nil
}

func (c *CacheHandler) clearCache(ctx echo.Context) error {
	// 型チェック
	switch c.client.(type) {
	case database.CacheClient:
		// 既にキャッシュされたデータを返している場合はキャッシュを削除する必要が無いのでreturnする
		log.Println("info: data is already cached.")
		return nil
	default:
		// Redisに接続
		r, err := database.ConnectRedis()
		if err != nil {
			log.Printf("error: Redisへの接続に失敗しました。: %v\n", err)
			return fmt.Errorf("failed to connection Redis: %w", err)
		}
		c.cache = r
	}

	ctx.Response().After(func() {
		// after Response
		_, err := c.cache.DEL(0) // ここにキーが引数として入る
		if err != nil {
			// とりあえずログのみ出力
			log.Printf("error: fail to SetByte Error: %v\n", err)
			// return c.NoContent(http.StatusInternalServerError)
		}
		log.Println("info: memo data is cached")
	})
	return nil
}

// // Chain returns a Middlewares type from a slice of middleware handlers.
// func Chain(middlewares ...func(Middleware) Middleware) []Middleware {
// 	return Middlewares(middlewares)
// }

// func (mx *ChainHandler) Middlewares() Middlewares {
// 	return mx.middlewares
// }

// // Handler builds and returns a http.Handler from the chain of middlewares,
// // with `h http.Handler` as the final handler.
// func (mws Middlewares) Handler(h http.Handler) http.Handler {
// 	return &ChainHandler{mws, h, chain(mws, h)}
// }

// // HandlerFunc builds and returns a http.Handler from the chain of middlewares,
// // with `h http.Handler` as the final handler.
// func (mws Middlewares) HandlerFunc(h http.HandlerFunc) http.Handler {
// 	return &ChainHandler{mws, h, chain(mws, h)}
// }

// // chain builds a http.Handler composed of an inline middleware stack and endpoint
// // handler in the order they are passed.
// func chain(middlewares []func(http.Handler) http.Handler, endpoint http.Handler) http.Handler {
// 	// Return ahead of time if there aren't any middlewares for the chain
// 	if len(middlewares) == 0 {
// 		return endpoint
// 	}

// 	// Wrap the end handler with the middleware chain
// 	h := middlewares[len(middlewares)-1](endpoint)
// 	for i := len(middlewares) - 2; i >= 0; i-- {
// 		h = middlewares[i](h)
// 	}

// 	return h
// }
