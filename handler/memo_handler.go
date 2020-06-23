package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"memoapp/model"
	"memoapp/repository"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type (
	//htmlData　htmlテンプレートに渡すデータ型
	htmlData map[string]interface{}

	//MemoAppOutput レスポンス用のデータ型
	MemoAppOutput struct {
		Memo    *model.Memo
		Message string
	}

	//Handler メモ用ハンドラー
	handler struct {
		repo repository.Memorepo
	}
)

func ProvideHandler(db *sqlx.DB) handler {
	return handler{
		repo: repository.Memorepo{
			DB: db,
		},
	}
}

func (h handler) MemoIndex(c echo.Context) error {

	memos, err := h.repo.GetMemoList()
	if err != nil {
		c.Logger().Errorf("failed to select db request : %v\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return render(c, "src/views/index.html",
		htmlData{
			"Memos": memos,
		})
}

func (h handler) MemoCreate(c echo.Context) error {

	var (
		memo = &model.Memo{}
		// out  MemoCreateOutput
	)

	err := c.Bind(memo)
	if err != nil {
		c.Logger().Errorf("failed to Bind request params : %v\n", err)
		return c.JSON(http.StatusBadRequest,
			MemoAppOutput{Message: "入力データに誤りがあります。"},
		)
	}

	// バリデート必要？モデルからValidata関数呼び出す？

	res, err := h.repo.MemoCreate(memo)
	if err != nil {
		c.Logger().Errorf("failed to insert memo data [%v] : %v\n", memo, err)
		return c.JSON(http.StatusInternalServerError, MemoAppOutput{})
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.Logger().Errorf("failed to get LastInsertId : %v\n", err)
		return c.JSON(http.StatusInternalServerError, MemoAppOutput{})
	}
	//①なぜint型でキャストしているのか？ / ②modelに関することはmodelで関数化しよう（setIdとか）
	memo.SetId(int(id)) // idをセット

	return c.JSON(http.StatusOK, MemoAppOutput{Memo: memo})
}

//削除機能
func (h handler) MemoDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Errorf("failed to converted to type int : %v\n", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	err = h.repo.MemoDelete(id)
	if err != nil {
		c.Logger().Errorf("failed to delete memo data [id :%v]: %v\n", id, err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Memo %d is deleted : ", id))
}
