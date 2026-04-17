package user

//go:generate easyjson -all user.go
type User struct {
	Email string
}
