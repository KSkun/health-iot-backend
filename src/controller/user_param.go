package controller

type ReqUserSimpleV1 struct {
	Name     string `json:"name" query:"name" validate:"required"`
	Password string `json:"password" query:"password" validate:"required,printascii,gte=6"`
}
