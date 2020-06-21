package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"memoapp/model"
	"memoapp/repository"

	"github.com/labstack/echo/v4"
)

func MemoIndex(c echo.Context) error {

	memos, err := repository.MemoListByCursor(0)

	if err != nil {

		c.Logger().Error(err.Error())

		return c.NoContent(http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Memos": memos,
	}

	return render(c, "src/views/index.html", data)
}

type MemoCreateOutput struct {
	Memo             *model.Memo
	Message          string
	ValidationErrors []string
}

func MemoCreate(c echo.Context) error {
	var memo model.Memo

	var out MemoCreateOutput

	if err := c.Bind(&memo); err != nil {

		c.Logger().Error(err.Error())

		return c.JSON(http.StatusBadRequest, out)

	}

	res, err := repository.MemoCreate(&memo)
	if err != nil {

		c.Logger().Error(err.Error())

		return c.JSON(http.StatusInternalServerError, out)
	}

	id, _ := res.LastInsertId()

	memo.ID = int(id)

	out.Memo = &memo

	return c.JSON(http.StatusOK, out)
}

//削除機能
func MemoDelete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := repository.MemoDelete(id); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusInternalServerError, "")
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Memo %d is deleted", id))
}
