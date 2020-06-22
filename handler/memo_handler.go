package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"memoapp/model"
	"memoapp/repository"

	"github.com/labstack/echo/v4"
)

type (
	htmlData      map[string]interface{}
	MemoAppOutput struct {
		Memo             *model.Memo
		Message          string
		ValidationErrors []string //なぜスライス？
	}
)

// type MemoAppOutput struct {
// 	Memo             *model.Memo
// 	Message          string
// 	ValidationErrors []string //なぜスライス？
// }

func MemoIndex(c echo.Context) error {

	memos, err := repository.GetMemoList()
	if err != nil {
		c.Logger().Errorf("failed to select db request : %v\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return render(c, "src/views/index.html",
		htmlData{
			"Memos": memos,
		})
}

func MemoCreate(c echo.Context) error {

	var (
		memo = &model.Memo{}
		// out  MemoCreateOutput
	)

	err := c.Bind(memo)
	if err != nil {
		c.Logger().Errorf("failed to Bind request params : %v\n", err)
		return c.JSON(http.StatusBadRequest,
			MemoAppOutput{ValidationErrors: []string{
				err.Error()},
			})
	}

	// バリデート必要？モデルからValidata関数呼び出す？

	res, err := repository.MemoCreate(memo)
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
func MemoDelete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Errorf("failed to converted to type int : %v\n", err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	err = repository.MemoDelete(id)
	if err != nil {
		c.Logger().Errorf("failed to delete memo data [id :%v]: %v\n", id, err)
		return c.JSON(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Memo %d is deleted : ", id))
}
