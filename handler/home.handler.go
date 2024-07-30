package handler

import (
	"github.com/a-h/templ"
	"github.com/alinsimion/expense_tracker/model"
	home "github.com/alinsimion/expense_tracker/view/home"
	"github.com/labstack/echo/v4"
)

func (es *ExpenseHandler) ShowHomescreen(ctx echo.Context) error {

	var returnView templ.Component
	user := ctx.Request().Context().Value(userContextKey)
	if user != nil {
		if user.(*model.User).Email != "" {
			return hxRedirect(ctx, "/expense")
		}
	} else {
		returnView = home.ShowHomeWithLayout()
	}

	return View(ctx, returnView)
}
