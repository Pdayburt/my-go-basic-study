package domain

// User 领域对象
type User struct {
	Id       int64 `json:"id"`
	Email    string
	Password string
}
