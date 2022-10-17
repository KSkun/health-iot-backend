package param

type ReqCreateUserV1 struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,printascii,gte=6"`
}
