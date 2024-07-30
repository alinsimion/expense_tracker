package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	"github.com/alinsimion/expense_tracker/model"
	"github.com/alinsimion/expense_tracker/util"
	"github.com/labstack/echo/v4"
)

const (
	toCookieName          = "to"
	accessTokenCookieName = "access_token"
	userContextKey        = "UserKey"
)

func View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	ctx := c.Request().Context()
	// fmt.Println(ctx.Value("User"))

	newCtx := context.WithValue(ctx, "User", c.Get("User"))

	return cmp.Render(newCtx, c.Response().Writer)
}

func (eh *ExpenseHandler) WithUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Request().Cookie(accessTokenCookieName)

		if err != nil {
			return next(c)
		}

		resp, err := Client.Auth.User(c.Request().Context(), cookie.Value)

		if err != nil {
			return next(c)
		}

		// fmt.Println(resp.AppMetadata)
		// user := model.User{
		// 	Email:    resp.Email,
		// 	LoggedIn: true,
		// }

		user := eh.store.GetUserByEmail(resp.Email)

		if user.Email == "" {
			slog.Error("No user in DB, strange")
		} else {
			user.LoggedIn = true
			avatarUrl, ok := resp.UserMetadata["avatar_url"]
			if ok {
				user.AvatarUrl = avatarUrl.(string)
			}
		}

		ctx := context.WithValue(c.Request().Context(), userContextKey, user)
		req := c.Request().WithContext(ctx)

		c.SetRequest(req)
		return next(c)

	}
}

func (eh *ExpenseHandler) WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		cookie, err := c.Request().Cookie(accessTokenCookieName)
		redirectPath := util.GetFullUrl("/login")

		if err != nil {
			slog.Error("Could not find cookie", "err", err.Error())

			path := c.Request().URL.Path

			if len(path) > 0 {
				redirectCookie := &http.Cookie{
					Value:    path,
					Name:     toCookieName,
					Path:     util.GetFullUrl("/"),
					HttpOnly: true,
					Secure:   true,
				}

				http.SetCookie(c.Response(), redirectCookie)
			}
			http.Redirect(c.Response(), c.Request(), redirectPath, http.StatusSeeOther)
			return err
		}

		resp, err := Client.Auth.User(c.Request().Context(), cookie.Value)

		if err != nil {

			// fmt.Println("Debugging")
			// fmt.Println(cookie)
			// fmt.Println(err.Error())
			// fmt.Println(resp)

			slog.Error("Could not auth with supabase", "err", err.Error())
			http.Redirect(c.Response(), c.Request(), redirectPath, http.StatusSeeOther)
			return nil
		}

		user := eh.store.GetUserByEmail(resp.Email)

		if user.Email == "" {
			slog.Error("No user in DB, strange")

			user := model.User{
				Email: resp.Email,
				// AvatarUrl: avatarUrl,
			}

			eh.store.CreateUser(user)

		} else {
			user.LoggedIn = true
			avatarUrl, ok := resp.UserMetadata["avatar_url"]
			if ok {
				user.AvatarUrl = avatarUrl.(string)
			}
		}

		// slog.Error("User din context inainte", "user_context", c.Request().Context().Value(userContextKey))

		ctx := context.WithValue(c.Request().Context(), userContextKey, user)

		req := c.Request().WithContext(ctx)

		c.SetRequest(req)

		// slog.Error("User din context dupa", "user_context", c.Request().Context().Value(userContextKey))

		return next(c)
	}
}

func hxRedirect(ctx echo.Context, url string) error {
	req := ctx.Request()
	res := ctx.Response()

	if len(req.Header.Get("HX-Request")) > 0 {
		res.Header().Set("HX-Redirect", url)
		res.WriteHeader(http.StatusSeeOther)
		return nil
	}
	http.Redirect(res, req, url, http.StatusSeeOther)
	return nil
}

func setAuthCookies(ctx echo.Context, accessToken string) {
	cookie := &http.Cookie{
		Value:    accessToken,
		Name:     accessTokenCookieName,
		Path:     util.GetFullUrl("/"),
		HttpOnly: true,
		Secure:   true,
		MaxAge:   72000,
	}

	http.SetCookie(ctx.Response(), cookie)

}
