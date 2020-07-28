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

		fmt.Println("c.Path(): " + c.Path())
		fmt.Println("c.RealIP(): " + c.RealIP())
		fmt.Printf("c.ParamValues(): %v\n", c.ParamValues())
		fmt.Printf("c.QueryParams(): %v\n", c.QueryParams())
		fmt.Printf("c.QueryString(): %v\n", c.QueryString())

		// DB接続
		database, err := h.Connect()
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

		data, err := handler(c)

		if err != nil {
			log.Printf("error: Internal Server Error: %v\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if data == nil {
			return c.NoContent(http.StatusOK)
		}
		if !h.HasCache {
			err := database.SetByte(data)
			if err != nil {
				log.Printf("error: fail to SetByte Error: %v\n", err)
				return c.NoContent(http.StatusInternalServerError)
			}
			log.Println("info: memo data is cached")
		}
		return c.JSONBlob(http.StatusOK, data)
	}
}
func (h *MemoHandler) endpointHandler(handler endPointHandler) echo.HandlerFunc {
	return func(c echo.Context) error {

		fmt.Println("c.Path(): " + c.Path())
		fmt.Println("c.RealIP(): " + c.RealIP())
		fmt.Printf("c.ParamValues(): %v\n", c.ParamValues())
		fmt.Printf("c.QueryParams(): %v\n", c.QueryParams())
		fmt.Printf("c.QueryString(): %v\n", c.QueryString())

		// DB接続
		database, err := database.ConnectMySql()
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

		data, err := handler(c)

		if err != nil {
			log.Printf("error: Internal Server Error: %v\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if data == nil {
			return c.NoContent(http.StatusOK)
		}
		return c.JSONBlob(http.StatusOK, data)
	}
}

func (h *MemoHandler) WithProviderFinalizer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// ↑ BEFORE
			err := next(c) // HandlerFunc : func(Context) error
			// この場合、AFTERの処理は実行され、エラーを返す
			// ↓ AFTER
			if err != nil {
				return err
			}
			return nil
			// log.Println("info: database connection is Closing...")
			// return h.repo.Close()
		}
	}
}

// func (controller ControllerBase) withProviderClient() echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			// TODO: ミドルウェアでプロバイダーにClientをもたせるようにする。
// 			// controller.Provider.ProvideStorageOperator().

// 			// ↑ BEFORE
// 			err := next(c) // HandlerFunc : func(Context) error
// 			// この場合、AFTERの処理は実行され、エラーを返す
// 			// ↓ AFTER
// 			if err != nil {
// 				return err
// 			}
// 			finalizeError := controller.Provider.Finalize(c.Request().Context())
// 			if finalizeError != nil {
// 				return finalizeError
// 			}
// 			return nil
// 		}
// 	}
// }
