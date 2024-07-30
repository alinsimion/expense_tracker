package handler

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/alinsimion/expense_tracker/model"
	"github.com/alinsimion/expense_tracker/util"
	view "github.com/alinsimion/expense_tracker/view/auth"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
)

var (
	Client *supabase.Client
)

func (eh *ExpenseHandler) HandleLoginCreate(ctx echo.Context) error {

	req := ctx.Request()

	credentials := supabase.UserCredentials{
		Email:    req.FormValue("email"),
		Password: req.FormValue("password"),
	}

	resp, err := Client.Auth.SignIn(req.Context(), credentials)

	if err != nil {
		slog.Error("Login errr", "Error", err.Error())
		loginErrors := view.LoginErrors{}
		loginErrors.InvalidCredentials = "The credentials you have entered are invalid"
		loginView := view.LoginForm(credentials, loginErrors)
		return View(ctx, loginView)
	}

	setAuthCookies(ctx, resp.AccessToken)

	user := eh.store.GetUserByEmail(resp.User.Email)

	if user.Email == "" {
		name, ok := resp.User.UserMetadata["name"].(string)
		if !ok {
			name = ""
		}

		avatarUrl, ok := resp.User.UserMetadata["avatar_url"].(string)
		if !ok {
			avatarUrl = ""
		}

		tempUser := model.User{
			Email:     resp.User.Email,
			Name:      name,
			AvatarUrl: avatarUrl,
		}
		eh.store.CreateUser(tempUser)
	}
	path, err := req.Cookie(toCookieName)

	toDeleteCookie := &http.Cookie{
		Name:     toCookieName,
		Value:    "",
		Path:     util.GetFullUrl("/"),
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(ctx.Response(), toDeleteCookie)

	redirectPath := util.GetFullUrl("/expense")

	if err != nil {
		slog.Error("Did not find redirect cookie", "err", err.Error())

	} else {
		if len(path.Value) > 0 {
			redirectPath = path.Value
		}
	}

	return hxRedirect(ctx, redirectPath)
}
func (eh *ExpenseHandler) HandleLoginCreateProvider(ctx echo.Context) error {

	slog.Debug("Redirect url", "url", os.Getenv("REDIRECT_URL"))
	resp, err := Client.Auth.SignInWithProvider(
		supabase.ProviderSignInOptions{
			Provider:   "google",
			RedirectTo: os.Getenv("REDIRECT_URL"),
		},
	)

	if err != nil {
		slog.Error("Error with google", "err", err.Error())
		return err
	}

	slog.Debug("Redirect resp url", "url", resp.URL)

	http.Redirect(ctx.Response(), ctx.Request(), resp.URL, http.StatusSeeOther)

	return nil
}
func (eh *ExpenseHandler) HandleLoginIndex(ctx echo.Context) error {

	loginView := view.ShowLoginWithLayout()
	return View(ctx, loginView)
}
func (eh *ExpenseHandler) HandleSignupIndex(ctx echo.Context) error {

	signupView := view.ShowSignupWithLayout()
	return View(ctx, signupView)
}
func (eh *ExpenseHandler) HandleSignupCreate(ctx echo.Context) error {
	req := ctx.Request()

	credentials := view.SignupCreds{
		Email:           req.FormValue("email"),
		Password:        req.FormValue("password"),
		ConfirmPassword: req.FormValue("confirm_password"),
	}

	// TODO: Compare Passwod to Confirm Password and return FormErrors if missmatch

	user, err := Client.Auth.SignUp(ctx.Request().Context(), supabase.UserCredentials{
		Email:    credentials.Email,
		Password: credentials.Password,
	})

	if err != nil {
		slog.Error("Error at signup", "err", err.Error())
		return err
	}

	// setAuthCookies(ctx, user.acc)

	return View(ctx, view.SignupSucces(user.Email))
}
func (eh *ExpenseHandler) HandleAuthCallback(ctx echo.Context) error {
	accessToken := ctx.Request().URL.Query().Get("access_token")

	if len(accessToken) == 0 {
		return View(ctx, view.CallbackScript())
	}

	setAuthCookies(ctx, accessToken)

	path, err := ctx.Request().Cookie(toCookieName)

	toDeleteCookie := &http.Cookie{
		Name:     toCookieName,
		Value:    "",
		Path:     util.GetFullUrl("/"),
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(ctx.Response(), toDeleteCookie)

	redirectPath := util.GetFullUrl("/expense")

	if err != nil {
		slog.Error("Did not find redirect cookie", "err", err.Error())

	} else {
		if len(path.Value) > 0 {
			redirectPath = util.GetFullUrl(path.Value)
		}
	}

	http.Redirect(ctx.Response(), ctx.Request(), redirectPath, http.StatusSeeOther)

	return nil
}
func (eh *ExpenseHandler) HandleLogoutCreate(ctx echo.Context) error {

	cookie := &http.Cookie{
		Value:    "",
		Name:     accessTokenCookieName,
		MaxAge:   -1,
		HttpOnly: true,
		Path:     util.GetFullUrl("/"),
		Secure:   true,
	}
	http.SetCookie(ctx.Response(), cookie)
	http.Redirect(ctx.Response(), ctx.Request(), util.GetFullUrl("/login"), http.StatusSeeOther)

	return nil
}
