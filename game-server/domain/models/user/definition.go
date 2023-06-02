package user

type GenderType string

const (
	Undefined GenderType = "undefined"
	Male      GenderType = "male"
	Female    GenderType = "female"
)

type User struct {
	Login  string     `json:"login"`
	Emain  string     `json:"email"`
	Gender GenderType `json:"gender"`
	// TODO IMAGE
}
