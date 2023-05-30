package user

type GenderType int

const (
	Undefined GenderType = 0
	Male      GenderType = 1
	Female    GenderType = 2
)

type User struct {
	Login  string
	Emain  string
	Gender GenderType
	// TODO IMAGE
}
