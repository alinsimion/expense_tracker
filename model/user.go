package model

const (
	UserContextKey = "UserKey"
)

type User struct {
	Id        int64
	Name      string
	Email     string
	LoggedIn  bool
	AvatarUrl string
}
